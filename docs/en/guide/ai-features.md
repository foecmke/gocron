# AI Features

gocron ships a set of AI capabilities, all powered by **your own OpenAI-compatible LLM** (cloud or local): natural-language-to-cron, failure-log diagnosis, an AI ops chat assistant, and MCP protocol integration.

## Prerequisite: configure the LLM

Go to **System → AI Config**, fill in and enable:

| Field | Description |
|-------|-------------|
| Enable | AI features are available only when enabled |
| Base URL | OpenAI-compatible endpoint, e.g. `https://api.openai.com/v1` or local `http://localhost:8000/v1` |
| API Key | Access key — **stored server-side only, never returned to the frontend** |
| Model | Model name, e.g. `gpt-4o`, `qwen2.5` |

Any **OpenAI-compatible** endpoint works: a cloud API, or a local model server (Ollama, vLLM, LM Studio, MLX, …). A local model gives you **fully offline** AI.

> Note: "thinking/reasoning" models are slower on multi-step tasks — this is expected; the chat panel streams the reasoning live and you can stop it anytime.

## Natural-language to Cron

When creating / editing a task, describe the schedule in natural language (e.g. "every day at 9:30 am", "hourly on weekdays"). The AI generates gocron's **6-field second-level cron expression** and previews the next run times for you to confirm.

## AI failure-log diagnosis

After a task run fails, click **Diagnose** on the task log. Based on the task config and the failed output, the AI returns:

- a one-sentence **root cause**
- actionable **fix suggestions**

## AI ops assistant (chat)

Click the robot icon in the top bar to open the chat panel and ask ops questions in natural language:

- "Which tasks failed recently?" / "Why did task X fail?" / "What task templates do I have?"
- **Streaming** output, with the model's reasoning shown live
- The AI calls **read-only tools** to fetch real data (tasks, logs, hosts, templates) and answers from the results — it does not fabricate
- Triggering a task run (`run_task`) **requires your explicit confirmation** and is audited — the AI never runs a task automatically
- Supports **Markdown rendering**, message copy, a resizable panel, and stop-anytime
- Reasoning and answers **follow your UI language**

## MCP integration

gocron has a built-in **MCP (Model Context Protocol)** server, so external MCP clients (e.g. Claude Desktop) can connect and use gocron as a tool to manage tasks.

- **Endpoint**: `https://<your-host>/mcp`
- Create an access token under **System → MCP Keys** (the token is shown in full only once at creation — save it)
- Client config example:

```json
{
  "mcpServers": {
    "gocron": {
      "url": "https://<your-host>/mcp",
      "headers": {
        "Authorization": "Bearer <your MCP token>"
      }
    }
  }
}
```

Available tools:

| Tool | Description |
|------|-------------|
| `list_tasks` | List tasks (filter by name / tag / status) |
| `get_task` | Get a single task's details |
| `query_task_logs` | Query execution logs (filter by task / status / keyword / time range) |
| `list_hosts` | List execution nodes |
| `list_templates` | List task templates |
| `diagnose_task_log` | Diagnose a failed log |
| `run_task` | Run a task now (**requires an admin token**) |

> Security: `run_task` is the only write operation and requires an admin token; the rest are read-only. Expose the `/mcp` endpoint over TLS (HTTPS).
