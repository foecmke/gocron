package diagnosis

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gocronx-team/gocron/internal/models"
	"github.com/gocronx-team/gocron/internal/modules/llm"
)

func TestParse(t *testing.T) {
	t.Run("clean json", func(t *testing.T) {
		r := Parse(`{"root_cause":"端口未监听","suggestions":["启动服务","检查防火墙"]}`)
		if r.RootCause != "端口未监听" || len(r.Suggestions) != 2 {
			t.Fatalf("unexpected: %+v", r)
		}
	})
	t.Run("fenced json", func(t *testing.T) {
		r := Parse("```json\n{\"root_cause\":\"x\",\"suggestions\":[\"y\"]}\n```")
		if r.RootCause != "x" || len(r.Suggestions) != 1 {
			t.Fatalf("unexpected: %+v", r)
		}
	})
	t.Run("non-json fallback", func(t *testing.T) {
		r := Parse("就是连不上而已")
		if r.RootCause != "就是连不上而已" || len(r.Suggestions) != 0 {
			t.Fatalf("expected raw fallback, got %+v", r)
		}
	})
}

func TestBuildPrompt(t *testing.T) {
	log := &models.TaskLog{Name: "backup", Protocol: models.TaskRPC, Command: "python3 b.py", Hostname: "node-1", Result: "boom"}
	p := BuildPrompt(log)
	for _, want := range []string{"backup", "RPC(Shell)", "python3 b.py", "node-1", "boom"} {
		if !strings.Contains(p, want) {
			t.Fatalf("prompt missing %q: %s", want, p)
		}
	}
}

func TestDiagnose_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"choices": []map[string]any{
				{"message": map[string]string{"role": "assistant", "content": `{"root_cause":"节点缺 python3","suggestions":["安装 python3"]}`}},
			},
		})
	}))
	defer srv.Close()

	client := llm.New(srv.URL, "sk-test", "gpt-test")
	log := &models.TaskLog{Name: "backup", Protocol: models.TaskRPC, Command: "python3 b.py", Result: "python3: not found"}
	r, err := Diagnose(context.Background(), client, log, false)
	if err != nil {
		t.Fatalf("Diagnose: %v", err)
	}
	if r.RootCause != "节点缺 python3" || len(r.Suggestions) != 1 {
		t.Fatalf("unexpected result: %+v", r)
	}
}
