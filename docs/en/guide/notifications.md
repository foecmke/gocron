# Notifications

Task execution can notify you via **Email, Slack, and Webhook**. Setup has two steps: ① configure the channels under **System → Notification**; ② on each task, set a **notification policy** (when to notify, which channel, and who receives it).

## 1. Global channel config (System → Notification)

The page has three tabs:

### Email
- **SMTP host / port / user / password**: the sending mailbox's SMTP server
- **Recipients**: one or more email addresses
- **Template**: the email body template (see template variables below)

### Slack
- **Webhook URL**: a Slack Incoming Webhook URL
- **Channels**: one or more channels
- **Template**: the message template

### Webhook
- **Webhook URL list**: each entry has a **name** and a **URL**; add as many as you need
- **Template**: the rendered template is sent as the **HTTP request body, POSTed to each URL**

## 2. Template variables

All three channels use Go template syntax with these variables:

| Variable | Meaning |
|----------|---------|
| <code v-pre>{{.TaskId}}</code> | Task ID |
| <code v-pre>{{.TaskName}}</code> | Task name |
| <code v-pre>{{.Status}}</code> | Execution status |
| <code v-pre>{{.Result}}</code> | Execution output |
| <code v-pre>{{.Remark}}</code> | Task remark |

Webhook template example (JSON):

```json
{
  "task": "{{.TaskName}}",
  "status": "{{.Status}}",
  "output": "{{.Result}}"
}
```

## 3. Per-task notification policy (New / Edit task → Notification)

Once the channels are configured, set the policy on a task:

- **When to notify**: Disabled / On failure only / Always / On keyword match
- **Channel**: Email / Slack / Webhook
- **Recipients**: pick from the configured recipients / channels / Webhook URLs
- **Keyword**: for "on keyword match", only notify when the execution output contains this keyword

## Full steps to configure a Webhook

1. **System → Notification → Webhook** tab
2. Edit the **template** (using the variables above, typically a JSON body)
3. Add an entry to the **Webhook URL list**: a **name** and the receiving **URL**, then save
4. New / Edit task → **Notification**:
   - When to notify: "On failure only" or "Always" (or "On keyword match" with a keyword)
   - Channel: **Webhook**
   - Recipients: select the URL you just configured
5. Save. After the task runs, the rendered template is **POSTed** to the matching Webhook URL per your policy.
