// Package ai 提供服务端的智能体对话端点：在服务端跑一个 LLM 工具调用循环，
// 复用既有 MCP 工具回答运维问题（哪些任务失败了、某任务详情、主机列表等）。
package ai

import (
	"context"
	"encoding/json"
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
Use the provided tools to look up real data before answering — never fabricate task names, ids, statuses, or log contents.
For task-log execution status: 0 = failed, 1 = running, 2 = success (finished), 3 = cancelled.
Be concise and answer in the same language the user used.`

// ChatMessage 是请求中的一条对话消息。
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Messages []ChatMessage `json:"messages"`
}

// Chat 运行一个有界的 LLM 工具调用循环并返回最终回复。
// 响应数据：{ "reply": string, "tools_used": []string }，
// tools_used 按工具实际被调用的顺序记录（含重复），不去重。
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

	messages := buildMessages(req.Messages)
	isAdmin := user.IsAdmin(c)
	tools := mcp.AgentToolDefs()

	ctx, cancel := context.WithTimeout(c.Request.Context(), chatTimeout)
	defer cancel()

	toolsUsed := make([]string, 0)
	lastContent := ""

	for i := 0; i < maxIterations; i++ {
		msg, err := client.ChatWithTools(ctx, messages, tools)
		if err != nil {
			logger.Errorf("AI对话#调用LLM失败#%s", err)
			base.RespondError(c, i18n.T(c, "ai_chat_failed"))
			return
		}

		if len(msg.ToolCalls) == 0 {
			base.RespondSuccess(c, "", gin.H{
				"reply":      msg.Content,
				"tools_used": toolsUsed,
			})
			return
		}

		lastContent = msg.Content
		messages = append(messages, msg)
		for _, tc := range msg.ToolCalls {
			toolsUsed = append(toolsUsed, tc.Function.Name)
			messages = append(messages, llm.Message{
				Role:       "tool",
				ToolCallID: tc.ID,
				Content:    callTool(tc.Function.Name, tc.Function.Arguments, isAdmin),
			})
		}
	}

	// 迭代用尽仍未给出终答：尽力返回最后一段助手内容，否则报错。
	if strings.TrimSpace(lastContent) != "" {
		base.RespondSuccess(c, "", gin.H{
			"reply":      lastContent,
			"tools_used": toolsUsed,
		})
		return
	}
	base.RespondError(c, i18n.T(c, "ai_chat_failed"))
}

// buildMessages 在用户消息前注入系统提示词。
func buildMessages(in []ChatMessage) []llm.Message {
	out := make([]llm.Message, 0, len(in)+1)
	out = append(out, llm.Message{Role: "system", Content: systemPrompt})
	for _, m := range in {
		out = append(out, llm.Message{Role: m.Role, Content: m.Content})
	}
	return out
}

// callTool 执行一次工具调用，把结果或错误信息序列化为字符串作为 tool 消息内容。
func callTool(name, args string, isAdmin bool) string {
	result, err := mcp.CallTool(name, []byte(args), isAdmin)
	if err != nil {
		return err.Error()
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		return err.Error()
	}
	return string(encoded)
}
