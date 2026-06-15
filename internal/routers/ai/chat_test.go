package ai

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gocronx-team/gocron/internal/models"
	"github.com/gocronx-team/gocron/internal/modules/logger"
	"github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestMain(m *testing.M) {
	logger.InitLogger()
	gin.SetMode(gin.TestMode)
	m.Run()
}

func setupDb(t *testing.T) func() {
	t.Helper()
	db, err := gorm.Open(gormlite.Open(":memory:"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.Setting{}, &models.Task{}, &models.TaskLog{}, &models.Host{}, &models.TaskHost{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	original := models.Db
	models.Db = db
	return func() { models.Db = original }
}

func enableLLM(t *testing.T, baseURL string) {
	t.Helper()
	if err := new(models.Setting).UpdateLLM(true, baseURL, "sk-test", "gpt-test"); err != nil {
		t.Fatalf("UpdateLLM: %v", err)
	}
}

func doChat(r *gin.Engine, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/api/ai/chat", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func newRouter() *gin.Engine {
	r := gin.New()
	r.POST("/api/ai/chat", Chat)
	return r
}

// sseFrame 是从 SSE 响应体中解析出的一条事件。
type sseFrame struct {
	event string
	data  map[string]any
}

func parseSSE(t *testing.T, body string) []sseFrame {
	t.Helper()
	var frames []sseFrame
	var cur sseFrame
	scanner := bufio.NewScanner(strings.NewReader(body))
	flush := func() {
		if cur.event != "" {
			frames = append(frames, cur)
		}
		cur = sseFrame{}
	}
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "event: "):
			cur.event = strings.TrimPrefix(line, "event: ")
		case strings.HasPrefix(line, "data: "):
			data := strings.TrimPrefix(line, "data: ")
			m := map[string]any{}
			_ = json.Unmarshal([]byte(data), &m)
			cur.data = m
		case line == "":
			flush()
		}
	}
	flush()
	return frames
}

func TestChat_StreamsToolLoopAndReplies(t *testing.T) {
	defer setupDb(t)()
	seedFailingLog(t)

	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		w.Header().Set("Content-Type", "text/event-stream")
		if n == 1 {
			// 第一轮：流式请求调用 query_task_logs 工具（参数跨分片）
			_, _ = io.WriteString(w, "data: {\"choices\":[{\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"c1\",\"type\":\"function\",\"function\":{\"name\":\"query_task_logs\"}}]}}]}\n\n")
			_, _ = io.WriteString(w, "data: {\"choices\":[{\"delta\":{\"tool_calls\":[{\"index\":0,\"function\":{\"arguments\":\"{\\\"status\\\":0}\"}}]}}]}\n\n")
			_, _ = io.WriteString(w, "data: [DONE]\n\n")
			return
		}
		// 第二轮：收到 tool 结果后流式给出终答
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"role":"tool"`) {
			t.Errorf("expected tool result message in 2nd request, body=%s", body)
		}
		_, _ = io.WriteString(w, "data: {\"choices\":[{\"delta\":{\"content\":\"There was \"}}]}\n\n")
		_, _ = io.WriteString(w, "data: {\"choices\":[{\"delta\":{\"content\":\"1 failed task.\"}}]}\n\n")
		_, _ = io.WriteString(w, "data: [DONE]\n\n")
	}))
	defer srv.Close()
	enableLLM(t, srv.URL)

	w := doChat(newRouter(), `{"messages":[{"role":"user","content":"昨晚哪些任务失败了?"}]}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/event-stream") {
		t.Fatalf("content-type = %q", ct)
	}
	frames := parseSSE(t, w.Body.String())

	var (
		sawToolCall   bool
		sawToolResult bool
		sawDone       bool
		reply         strings.Builder
	)
	for _, f := range frames {
		switch f.event {
		case "tool_call":
			if f.data["name"] == "query_task_logs" {
				sawToolCall = true
			}
		case "tool_result":
			if f.data["name"] == "query_task_logs" && f.data["ok"] == true {
				sawToolResult = true
			}
		case "message":
			if s, ok := f.data["content"].(string); ok {
				reply.WriteString(s)
			}
		case "done":
			sawDone = true
		case "error":
			t.Fatalf("unexpected error event: %v", f.data)
		}
	}
	if !sawToolCall {
		t.Errorf("missing tool_call event for query_task_logs")
	}
	if !sawToolResult {
		t.Errorf("missing tool_result event with ok=true")
	}
	if reply.String() != "There was 1 failed task." {
		t.Errorf("assembled reply = %q", reply.String())
	}
	if !sawDone {
		t.Errorf("missing trailing done event")
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Fatalf("expected 2 llm calls, got %d", calls)
	}
}

func TestChat_EmptyMessages(t *testing.T) {
	defer setupDb(t)()
	enableLLM(t, "http://unused.example")

	w := doChat(newRouter(), `{"messages":[]}`)
	if ct := w.Header().Get("Content-Type"); strings.Contains(ct, "text/event-stream") {
		t.Fatalf("expected JSON error (no SSE), content-type = %q", ct)
	}
	var env map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &env)
	if code, _ := env["code"].(float64); code == 0 {
		t.Fatalf("expected failure code for empty messages, got %+v", env)
	}
}

func TestChat_LLMNotConfigured(t *testing.T) {
	defer setupDb(t)()
	// 不启用 LLM：FromSettings 返回 ErrNotConfigured

	w := doChat(newRouter(), `{"messages":[{"role":"user","content":"hi"}]}`)
	if ct := w.Header().Get("Content-Type"); strings.Contains(ct, "text/event-stream") {
		t.Fatalf("expected JSON error (no SSE), content-type = %q", ct)
	}
	var env map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &env)
	if code, _ := env["code"].(float64); code == 0 {
		t.Fatalf("expected failure code when llm not configured, got %+v", env)
	}
	if msg, _ := env["message"].(string); !strings.Contains(msg, "AI") {
		t.Fatalf("expected llm_not_configured message, got %q", msg)
	}
}

func seedFailingLog(t *testing.T) {
	t.Helper()
	log := models.TaskLog{TaskId: 1, Name: "nightly", Status: models.Disabled}
	if err := models.Db.Create(&log).Error; err != nil {
		t.Fatalf("seed log: %v", err)
	}
}
