package mcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gocronx-team/gocron/internal/models"
	"github.com/gocronx-team/gocron/internal/modules/diagnosis"
	"github.com/gocronx-team/gocron/internal/modules/docsearch"
	"github.com/gocronx-team/gocron/internal/modules/llm"
	"github.com/gocronx-team/gocron/internal/service"
)

const maxPageSize = 100

// ── list_tasks ──────────────────────────────────────────────────────────────

type listTasksInput struct {
	Name     string `json:"name,omitempty" jsonschema:"按任务名称模糊过滤"`
	Tag      string `json:"tag,omitempty" jsonschema:"按标签模糊过滤"`
	Status   *int   `json:"status,omitempty" jsonschema:"按状态过滤，0 表示禁用，1 表示启用，省略则返回全部"`
	Page     int    `json:"page,omitempty" jsonschema:"页码，从 1 开始，默认 1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"每页数量，默认 20，最大 100"`
}

type listTasksOutput struct {
	Total int64         `json:"total"`
	Tasks []models.Task `json:"tasks"`
}

func listTasks(in listTasksInput) (listTasksOutput, error) {
	params := models.CommonMap{
		"Page":     normalizePage(in.Page),
		"PageSize": normalizePageSize(in.PageSize),
	}
	if in.Name != "" {
		params["Name"] = in.Name
	}
	if in.Tag != "" {
		params["Tag"] = in.Tag
	}
	if in.Status != nil {
		params["Status"] = *in.Status
	}

	taskModel := new(models.Task)
	total, err := taskModel.Total(params)
	if err != nil {
		return listTasksOutput{}, err
	}
	tasks, err := taskModel.List(params)
	if err != nil {
		return listTasksOutput{}, err
	}
	for i, item := range tasks {
		tasks[i].NextRunTime = models.NextRunTime(service.ServiceTask.NextRunTime(item))
	}
	return listTasksOutput{Total: total, Tasks: tasks}, nil
}

// ── get_task ────────────────────────────────────────────────────────────────

type getTaskInput struct {
	Id int `json:"id" jsonschema:"任务 ID"`
}

func getTask(in getTaskInput) (models.Task, error) {
	taskModel := new(models.Task)
	task, err := taskModel.Detail(in.Id)
	if err != nil {
		return models.Task{}, err
	}
	if task.Id != 0 {
		task.NextRunTime = models.NextRunTime(service.ServiceTask.NextRunTime(task))
	}
	return task, nil
}

// ── query_task_logs ───────────────────────────────────────────────────────────

type queryTaskLogsInput struct {
	TaskId    int    `json:"task_id,omitempty" jsonschema:"按任务 ID 过滤，省略则返回全部任务的日志"`
	Status    *int   `json:"status,omitempty" jsonschema:"按执行状态过滤，0 表示失败，1 表示执行中，2 表示成功（已完成），3 表示取消"`
	Keyword   string `json:"keyword,omitempty" jsonschema:"按任务名或执行输出模糊搜索"`
	StartTime string `json:"start_time,omitempty" jsonschema:"只返回该时间(含)之后的日志，格式 2006-01-02 15:04:05 或 2006-01-02，按服务器时区解析"`
	EndTime   string `json:"end_time,omitempty" jsonschema:"只返回该时间(不含)之前的日志，格式同 start_time"`
	Page      int    `json:"page,omitempty" jsonschema:"页码，从 1 开始，默认 1"`
	PageSize  int    `json:"page_size,omitempty" jsonschema:"每页数量，默认 20，最大 100"`
}

type queryTaskLogsOutput struct {
	Total int64            `json:"total"`
	Logs  []models.TaskLog `json:"logs"`
}

func queryTaskLogs(in queryTaskLogsInput) (queryTaskLogsOutput, error) {
	params := models.CommonMap{
		"Page":     normalizePage(in.Page),
		"PageSize": normalizePageSize(in.PageSize),
	}
	if in.TaskId > 0 {
		params["TaskId"] = in.TaskId
	}
	if in.Status != nil {
		params["Status"] = *in.Status
	}
	if in.Keyword != "" {
		params["Keyword"] = in.Keyword
	}
	if t, ok := parseTimeArg(in.StartTime); ok {
		params["StartTime"] = t
	}
	if t, ok := parseTimeArg(in.EndTime); ok {
		params["EndTime"] = t
	}

	logModel := new(models.TaskLog)
	total, err := logModel.Total(params)
	if err != nil {
		return queryTaskLogsOutput{}, err
	}
	logs, err := logModel.List(params)
	if err != nil {
		return queryTaskLogsOutput{}, err
	}
	return queryTaskLogsOutput{Total: total, Logs: logs}, nil
}

// ── diagnose_task_log ─────────────────────────────────────────────────────────

type diagnoseTaskLogInput struct {
	Id int `json:"id" jsonschema:"要诊断的任务执行日志 ID"`
}

// diagnoseTaskLog 对某条失败日志做 LLM 归因，复用 diagnosis 包（与 HTTP 接口同源）。
func diagnoseTaskLog(in diagnoseTaskLogInput) (diagnosis.Result, error) {
	logModel := new(models.TaskLog)
	if err := logModel.Find(int64(in.Id)); err != nil {
		return diagnosis.Result{}, err
	}
	if strings.TrimSpace(logModel.Result) == "" {
		return diagnosis.Result{}, fmt.Errorf("task log %d has no execution output to diagnose", in.Id)
	}
	client, err := llm.FromSettings()
	if err != nil {
		return diagnosis.Result{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), diagnosis.Timeout)
	defer cancel()
	return diagnosis.Diagnose(ctx, client, logModel, false)
}

// ── list_hosts ────────────────────────────────────────────────────────────────

type listHostsOutput struct {
	Hosts []models.Host `json:"hosts"`
}

func listHosts() (listHostsOutput, error) {
	hostModel := new(models.Host)
	hosts, err := hostModel.List(models.CommonMap{"Page": 1, "PageSize": maxPageSize})
	if err != nil {
		return listHostsOutput{}, err
	}
	return listHostsOutput{Hosts: hosts}, nil
}

// ── search_docs ───────────────────────────────────────────────────────────────

type searchDocsInput struct {
	Query string `json:"query" jsonschema:"要在 gocron 官方文档中检索的关键词或问题"`
	TopN  int    `json:"top_n,omitempty" jsonschema:"返回最相关的片段数，默认 4，最大 8"`
}

type searchDocsOutput struct {
	Results []docsearch.Chunk `json:"results"`
}

// searchDocs 在内嵌的 gocron 文档里做离线关键词检索，供模型回答"怎么用/是否支持"类问题。
func searchDocs(in searchDocsInput) (searchDocsOutput, error) {
	n := in.TopN
	if n > 8 {
		n = 8
	}
	return searchDocsOutput{Results: docsearch.Search(in.Query, n)}, nil
}

// ── list_templates ────────────────────────────────────────────────────────────

type listTemplatesInput struct {
	Category string `json:"category,omitempty" jsonschema:"按分类过滤，省略则返回全部"`
	Page     int    `json:"page,omitempty" jsonschema:"页码，从 1 开始，默认 1"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"每页数量，默认 20，最大 100"`
}

// templateSummary 是任务模板的精简视图（不含 command 正文等大字段，控制 token）。
type templateSummary struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Protocol    int8   `json:"protocol"` // 1=HTTP, 2=Shell
	Spec        string `json:"spec"`     // 默认 cron 表达式，可能为空
	IsBuiltin   int8   `json:"is_builtin"`
	UsageCount  int    `json:"usage_count"`
}

type listTemplatesOutput struct {
	Total     int64             `json:"total"`
	Templates []templateSummary `json:"templates"`
}

func listTemplates(in listTemplatesInput) (listTemplatesOutput, error) {
	params := models.CommonMap{
		"Page":     normalizePage(in.Page),
		"PageSize": normalizePageSize(in.PageSize),
	}
	if in.Category != "" {
		params["Category"] = in.Category
	}
	tmplModel := new(models.TaskTemplate)
	total, err := tmplModel.Total(params)
	if err != nil {
		return listTemplatesOutput{}, err
	}
	list, err := tmplModel.List(params)
	if err != nil {
		return listTemplatesOutput{}, err
	}
	out := make([]templateSummary, 0, len(list))
	for _, tpl := range list {
		out = append(out, templateSummary{
			Id:          tpl.Id,
			Name:        tpl.Name,
			Description: tpl.Description,
			Category:    tpl.Category,
			Protocol:    tpl.Protocol,
			Spec:        tpl.Spec,
			IsBuiltin:   tpl.IsBuiltin,
			UsageCount:  tpl.UsageCount,
		})
	}
	return listTemplatesOutput{Total: total, Templates: out}, nil
}

// ── run_task ──────────────────────────────────────────────────────────────────

type runTaskInput struct {
	Id int `json:"id" jsonschema:"要立即执行的任务 ID"`
}

type runTaskOutput struct {
	Started bool   `json:"started"`
	Message string `json:"message"`
}

func runTask(in runTaskInput) (runTaskOutput, error) {
	taskModel := new(models.Task)
	task, err := taskModel.Detail(in.Id)
	if err != nil {
		return runTaskOutput{}, err
	}
	if task.Id <= 0 {
		return runTaskOutput{Started: false, Message: "task not found"}, nil
	}
	task.Spec = "MCP manual run"
	service.ServiceTask.Run(task)
	return runTaskOutput{Started: true, Message: "task started, check logs for result"}, nil
}

func normalizePage(page int) int {
	if page <= 0 {
		return 1
	}
	return page
}

func normalizePageSize(size int) int {
	if size <= 0 {
		return 20
	}
	if size > maxPageSize {
		return maxPageSize
	}
	return size
}

// parseTimeArg 解析模型传入的时间字符串，按服务器本地时区处理（日志 start_time 为本地时间）。
// 支持 RFC3339、"2006-01-02 15:04:05"、"2006-01-02" 三种格式；无法解析时返回 false 表示不过滤。
func parseTimeArg(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false
	}
	layouts := []string{time.RFC3339, "2006-01-02 15:04:05", "2006-01-02"}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
