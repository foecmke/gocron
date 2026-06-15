package llm

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChat_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer sk-test" {
			t.Errorf("auth header = %q", got)
		}
		body, _ := io.ReadAll(r.Body)
		var req chatRequest
		_ = json.Unmarshal(body, &req)
		if req.Model != "gpt-test" {
			t.Errorf("model = %q", req.Model)
		}
		if len(req.Messages) != 2 || req.Messages[0].Role != "system" {
			t.Errorf("unexpected messages: %+v", req.Messages)
		}
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"  0 9 * * 1-5  "}}]}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "sk-test", "gpt-test")
	out, err := c.Chat(context.Background(), "sys", "user")
	if err != nil {
		t.Fatalf("Chat: %v", err)
	}
	if out != "0 9 * * 1-5" {
		t.Errorf("expected trimmed content, got %q", out)
	}
}

func TestChat_BaseURLTrailingSlash(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}]}`))
	}))
	defer srv.Close()

	// baseURL 带尾斜杠也应正确拼接，不出现双斜杠
	c := New(srv.URL+"/", "k", "m")
	if _, err := c.Chat(context.Background(), "s", "u"); err != nil {
		t.Fatalf("Chat: %v", err)
	}
	if gotPath != "/chat/completions" {
		t.Errorf("path = %q", gotPath)
	}
}

func TestChat_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"message":"invalid api key"}}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "bad", "m")
	_, err := c.Chat(context.Background(), "s", "u")
	if err == nil || !strings.Contains(err.Error(), "invalid api key") {
		t.Fatalf("expected error with upstream message, got %v", err)
	}
}

func TestChat_EmptyChoices(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"choices":[]}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "k", "m")
	if _, err := c.Chat(context.Background(), "s", "u"); err != ErrEmptyResponse {
		t.Fatalf("expected ErrEmptyResponse, got %v", err)
	}
}

func TestChatWithTools_ToolCallResponse(t *testing.T) {
	var gotTools int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer sk-test" {
			t.Errorf("auth header = %q", got)
		}
		body, _ := io.ReadAll(r.Body)
		var req toolChatRequest
		_ = json.Unmarshal(body, &req)
		gotTools = len(req.Tools)
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"",` +
			`"tool_calls":[{"id":"call_1","type":"function","function":{"name":"query_task_logs","arguments":"{\"status\":0}"}}]}}]}`))
	}))
	defer srv.Close()

	tools := []Tool{{Type: "function", Function: ToolFunction{Name: "query_task_logs", Parameters: json.RawMessage(`{"type":"object"}`)}}}
	c := New(srv.URL, "sk-test", "gpt-test")
	msg, err := c.ChatWithTools(context.Background(), []Message{{Role: "user", Content: "what failed?"}}, tools)
	if err != nil {
		t.Fatalf("ChatWithTools: %v", err)
	}
	if gotTools != 1 {
		t.Errorf("expected 1 tool in request body, got %d", gotTools)
	}
	if len(msg.ToolCalls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(msg.ToolCalls))
	}
	tc := msg.ToolCalls[0]
	if tc.ID != "call_1" || tc.Function.Name != "query_task_logs" || tc.Function.Arguments != `{"status":0}` {
		t.Errorf("unexpected tool call: %+v", tc)
	}
}

func TestChatWithTools_ContentResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"all good"}}]}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "sk-test", "gpt-test")
	msg, err := c.ChatWithTools(context.Background(), []Message{{Role: "user", Content: "status?"}}, nil)
	if err != nil {
		t.Fatalf("ChatWithTools: %v", err)
	}
	if len(msg.ToolCalls) != 0 {
		t.Errorf("expected no tool calls, got %d", len(msg.ToolCalls))
	}
	if msg.Content != "all good" {
		t.Errorf("content = %q", msg.Content)
	}
}

func TestChatWithTools_EmptyChoices(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"choices":[]}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "k", "m")
	if _, err := c.ChatWithTools(context.Background(), []Message{{Role: "user", Content: "x"}}, nil); err != ErrEmptyResponse {
		t.Fatalf("expected ErrEmptyResponse, got %v", err)
	}
}
