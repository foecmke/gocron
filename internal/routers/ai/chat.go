// Package ai 提供服务端的智能体对话端点：在服务端跑一个 LLM 工具调用循环，
// 复用既有 MCP 工具回答运维问题（哪些任务失败了、某任务详情、主机列表等），
// 并以 SSE 流式把内容增量、工具调用、工具结果推送给前端。
package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocronx-team/gocron/internal/mcp"
	"github.com/gocronx-team/gocron/internal/modules/i18n"
	"github.com/gocronx-team/gocron/internal/modules/llm"
	"github.com/gocronx-team/gocron/internal/modules/logger"
	"github.com/gocronx-team/gocron/internal/routers/base"
	"github.com/gocronx-team/gocron/internal/routers/user"
)

const (
	chatTimeout   = 120 * time.Second
	maxIterations = 6
)

// systemPrompt 约束模型角色与行为：用提供的工具回答运维问题，简洁，不编造任务数据。
const systemPrompt = `You are the AI ops assistant embedded in gocron, a distributed cron task scheduler.
Users ask operational questions about scheduled tasks, their execution logs, and the hosts that run them.

Operating principles:
- Tool use: call a tool ONLY when you need live data from this system — e.g. which tasks/hosts exist, execution logs, or why a specific run failed. For how-to / syntax / conceptual questions (e.g. how to create a task, what cron syntax to use), answer DIRECTLY from the knowledge below; do NOT call any tool.
- CRITICAL: never end your turn by only announcing an action (e.g. "let me check the tasks first"). In a single turn you must EITHER actually emit the tool call(s) you need, OR give the complete final answer. Do not stop after a preamble.
- When you do use tools, look up real data before concluding — never fabricate task names, ids, statuses, or log contents.
- For task-log execution status: 0 = failed, 1 = running, 2 = success (finished), 3 = cancelled.
- Be concise and answer in the same language the user used.

Cron syntax (IMPORTANT — gocron uses SECOND-level cron, not standard 5-field Unix cron):
- A spec has 6 space-separated fields: second minute hour day-of-month month day-of-week (seconds come FIRST).
- Day-of-week: 0 = Sunday, 1-5 = Mon-Fri, 6 = Saturday.
- When sub-second precision is not needed, the second field is 0. Examples:
  every minute "0 * * * * *"; every 5 minutes "0 */5 * * * *"; every 20 seconds "*/20 * * * * *";
  daily at 09:30 "0 30 9 * * *"; weekdays at 09:00 "0 0 9 * * 1-5"; 1st of month at 00:00 "0 0 0 1 * *".
- Shortcut descriptors are also supported: @yearly, @monthly, @weekly, @daily (midnight), @hourly,
  @every <duration> (e.g. "@every 30s", "@every 1m20s"), and @reboot (run once at startup).
- Do NOT describe gocron as using 5-field cron; that is incorrect for this system.`

// ChatMessage 是请求中的一条对话消息。
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Messages []ChatMessage `json:"messages"`
}

// sseEvent 是发送给前端的一条 SSE 事件。
type sseEvent struct {
	event string
	data  any
}

// Chat 运行一个有界的 LLM 工具调用循环，并以 SSE 流式推送结果。
// 事件契约：
//   - reasoning    {"content": "<delta>"}                       思考过程增量（思考型模型）
//   - message      {"content": "<delta>"}                       内容增量
//   - tool_call    {"id","name","arguments"}                    模型决定调用工具
//   - tool_result  {"id","name","ok": true|false}               工具执行完成（不回传结果体）
//   - error        {"message": "<msg>"}                          运行期错误
//   - done         {}                                           始终最后发送
//
// 请求校验在写入 SSE 头之前完成，校验失败/未配置时以普通 JSON 错误响应返回，
// 与应用其余接口保持一致。
func Chat(c *gin.Context) {
	var req chatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		base.RespondError(c, i18n.T(c, "param_error"))
		return
	}
	if len(req.Messages) == 0 || strings.TrimSpace(req.Messages[len(req.Messages)-1].Content) == "" {
		base.RespondError(c, i18n.T(c, "param_error"))
		return
	}

	client, err := llm.FromSettings()
	if err != nil {
		base.RespondError(c, i18n.T(c, "llm_not_configured"))
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	sendEvent := func(ev sseEvent) {
		payload, err := json.Marshal(ev.data)
		if err != nil {
			payload = []byte(`{}`)
		}
		_, _ = fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", ev.event, payload)
		c.Writer.Flush()
	}

	messages := buildMessages(req.Messages)
	isAdmin := user.IsAdmin(c)
	tools := mcp.AgentToolDefs()

	ctx, cancel := context.WithTimeout(c.Request.Context(), chatTimeout)
	defer cancel()

	defer sendEvent(sseEvent{event: "done", data: map[string]any{}})

	for i := 0; i < maxIterations; i++ {
		msg, err := client.ChatStream(ctx, messages, tools,
			func(delta string) {
				sendEvent(sseEvent{event: "message", data: map[string]string{"content": delta}})
			},
			func(delta string) {
				sendEvent(sseEvent{event: "reasoning", data: map[string]string{"content": delta}})
			})
		if err != nil {
			logger.Errorf("AI对话#调用LLM失败#轮次%d#%s", i, err)
			sendEvent(sseEvent{event: "error", data: map[string]string{"message": i18n.T(c, "ai_chat_failed")}})
			return
		}
		logger.Infof("AI对话#轮次%d#内容长度%d#工具数%d", i, len(msg.Content), len(msg.ToolCalls))

		// 没有工具调用：模型已通过 message 事件流出终答。
		if len(msg.ToolCalls) == 0 {
			return
		}

		messages = append(messages, msg)
		for _, tc := range msg.ToolCalls {
			sendEvent(sseEvent{event: "tool_call", data: map[string]string{
				"id":        tc.ID,
				"name":      tc.Function.Name,
				"arguments": tc.Function.Arguments,
			}})

			result, terr := safeCallTool(tc.Function.Name, tc.Function.Arguments, isAdmin)
			if terr != nil {
				logger.Errorf("AI对话#工具失败#%s#args=%s#%s", tc.Function.Name, tc.Function.Arguments, terr)
			} else {
				logger.Infof("AI对话#工具成功#%s", tc.Function.Name)
			}
			sendEvent(sseEvent{event: "tool_result", data: map[string]any{
				"id":   tc.ID,
				"name": tc.Function.Name,
				"ok":   terr == nil,
			}})

			messages = append(messages, llm.Message{
				Role:       "tool",
				ToolCallID: tc.ID,
				Content:    toolResultContent(result, terr),
			})
		}
	}

	// 达到最大轮次仍未给出终答（模型一直在调工具/反复打转）。
	logger.Errorf("AI对话#达到最大轮次%d仍无终答", maxIterations)
	sendEvent(sseEvent{event: "error", data: map[string]string{"message": i18n.T(c, "ai_chat_failed")}})
}

// safeCallTool 执行工具调用并兜底 panic，避免单个工具异常中断整个 SSE 流。
func safeCallTool(name, args string, isAdmin bool) (result any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("tool %s panicked: %v", name, r)
		}
	}()
	return mcp.CallTool(name, []byte(args), isAdmin)
}

// buildMessages 在用户消息前注入系统提示词（含当前服务器时间）。
func buildMessages(in []ChatMessage) []llm.Message {
	prompt := systemPrompt + fmt.Sprintf("\n\nCurrent server time: %s", time.Now().Format("2006-01-02 15:04:05 MST"))
	out := make([]llm.Message, 0, len(in)+1)
	out = append(out, llm.Message{Role: "system", Content: prompt})
	for _, m := range in {
		out = append(out, llm.Message{Role: m.Role, Content: m.Content})
	}
	return out
}

// toolResultContent 把工具结果或错误信息序列化为 tool 消息内容。
func toolResultContent(result any, err error) string {
	if err != nil {
		return err.Error()
	}
	encoded, mErr := json.Marshal(result)
	if mErr != nil {
		return mErr.Error()
	}
	return string(encoded)
}
