package ai

import (
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

func doChat(r *gin.Engine, body string) (int, map[string]any) {
	req := httptest.NewRequest(http.MethodPost, "/api/ai/chat", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var env map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &env)
	return w.Code, env
}

func newRouter() *gin.Engine {
	r := gin.New()
	r.POST("/api/ai/chat", Chat)
	return r
}

func TestChat_RunsToolLoopAndReplies(t *testing.T) {
	defer setupDb(t)()
	seedFailingLog(t)

	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n == 1 {
			// 第一轮：请求调用 query_task_logs 工具
			_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"",` +
				`"tool_calls":[{"id":"c1","type":"function","function":{"name":"query_task_logs","arguments":"{\"status\":0}"}}]}}]}`))
			return
		}
		// 第二轮：收到 tool 结果后给出终答
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"role":"tool"`) {
			t.Errorf("expected tool result message in 2nd request, body=%s", body)
		}
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"There was 1 failed task last night."}}]}`))
	}))
	defer srv.Close()
	enableLLM(t, srv.URL)

	code, env := doChat(newRouter(), `{"messages":[{"role":"user","content":"昨晚哪些任务失败了?"}]}`)
	if code != http.StatusOK {
		t.Fatalf("status = %d", code)
	}
	data, _ := env["data"].(map[string]any)
	if data == nil {
		t.Fatalf("missing data: %+v", env)
	}
	if reply, _ := data["reply"].(string); !strings.Contains(reply, "failed task") {
		t.Fatalf("unexpected reply: %v", data["reply"])
	}
	used, _ := data["tools_used"].([]any)
	if len(used) != 1 || used[0] != "query_task_logs" {
		t.Fatalf("unexpected tools_used: %v", data["tools_used"])
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Fatalf("expected 2 llm calls, got %d", calls)
	}
}

func TestChat_EmptyMessages(t *testing.T) {
	defer setupDb(t)()
	enableLLM(t, "http://unused.example")

	_, env := doChat(newRouter(), `{"messages":[]}`)
	if env["code"].(float64) == 0 {
		t.Fatalf("expected failure code for empty messages, got %+v", env)
	}
}

func TestChat_LLMNotConfigured(t *testing.T) {
	defer setupDb(t)()
	// 不启用 LLM：FromSettings 返回 ErrNotConfigured

	_, env := doChat(newRouter(), `{"messages":[{"role":"user","content":"hi"}]}`)
	if env["code"].(float64) == 0 {
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
