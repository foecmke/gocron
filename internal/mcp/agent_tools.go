package mcp

import (
	"encoding/json"
	"fmt"

	"github.com/gocronx-team/gocron/internal/modules/llm"
)

// AgentToolDefs 以 OpenAI function 形式返回 5 个工具定义，供产品内 AI 对话循环复用，
// 与 MCP server 注册的工具保持同一套底层实现（tools.go），是工具逻辑的单一来源。
func AgentToolDefs() []llm.Tool {
	return []llm.Tool{
		tool("list_tasks", "List scheduled tasks with optional name/tag/status filters and pagination.", `{
			"type": "object",
			"properties": {
				"name": {"type": "string", "description": "Fuzzy filter by task name"},
				"tag": {"type": "string", "description": "Fuzzy filter by tag"},
				"status": {"type": "integer", "description": "Filter by status: 0 disabled, 1 enabled; omit for all"},
				"page": {"type": "integer", "description": "Page number, starts at 1, default 1"},
				"page_size": {"type": "integer", "description": "Items per page, default 20, max 100"}
			}
		}`),
		tool("get_task", "Get the full configuration of a single scheduled task by id.", `{
			"type": "object",
			"properties": {
				"id": {"type": "integer", "description": "Task id"}
			},
			"required": ["id"]
		}`),
		tool("query_task_logs", "Query task execution logs, filterable by task id, execution status, keyword, and time range.", `{
			"type": "object",
			"properties": {
				"task_id": {"type": "integer", "description": "Filter by task id; omit for all tasks"},
				"status": {"type": "integer", "description": "Filter by execution status: 0 failed, 1 running, 2 success (finished), 3 cancelled"},
				"keyword": {"type": "string", "description": "Fuzzy match in task name or execution output"},
				"start_time": {"type": "string", "description": "Only logs at/after this time. Format '2006-01-02 15:04:05' or '2006-01-02' in server timezone. Use the current server time to compute ranges like 'last night'."},
				"end_time": {"type": "string", "description": "Only logs strictly before this time. Same format as start_time."},
				"page": {"type": "integer", "description": "Page number, starts at 1, default 1"},
				"page_size": {"type": "integer", "description": "Items per page, default 20, max 100"}
			}
		}`),
		tool("list_hosts", "List all execution nodes (hosts).", `{
			"type": "object",
			"properties": {}
		}`),
		tool("run_task", "Trigger an immediate manual run of a task by id (requires admin).", `{
			"type": "object",
			"properties": {
				"id": {"type": "integer", "description": "Id of the task to run now"}
			},
			"required": ["id"]
		}`),
	}
}

func tool(name, description, parameters string) llm.Tool {
	return llm.Tool{
		Type: "function",
		Function: llm.ToolFunction{
			Name:        name,
			Description: description,
			Parameters:  json.RawMessage(parameters),
		},
	}
}

// CallTool 按名称分发到对应的工具实现，复用 tools.go 中的私有函数。
// argsJSON 为模型生成的参数（JSON 编码）。run_task 在非管理员时返回 errAdminRequired。
// 未知工具名返回错误。
func CallTool(name string, argsJSON []byte, isAdmin bool) (any, error) {
	switch name {
	case "list_tasks":
		var in listTasksInput
		if err := unmarshalArgs(argsJSON, &in); err != nil {
			return nil, err
		}
		return listTasks(in)
	case "get_task":
		var in getTaskInput
		if err := unmarshalArgs(argsJSON, &in); err != nil {
			return nil, err
		}
		return getTask(in)
	case "query_task_logs":
		var in queryTaskLogsInput
		if err := unmarshalArgs(argsJSON, &in); err != nil {
			return nil, err
		}
		return queryTaskLogs(in)
	case "list_hosts":
		return listHosts()
	case "run_task":
		if !isAdmin {
			return nil, errAdminRequired
		}
		var in runTaskInput
		if err := unmarshalArgs(argsJSON, &in); err != nil {
			return nil, err
		}
		return runTask(in)
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

// unmarshalArgs 容忍空参数（模型对无参工具可能传空串），空时按零值处理。
func unmarshalArgs(argsJSON []byte, dst any) error {
	if len(argsJSON) == 0 {
		return nil
	}
	if err := json.Unmarshal(argsJSON, dst); err != nil {
		return fmt.Errorf("invalid tool arguments: %w", err)
	}
	return nil
}
