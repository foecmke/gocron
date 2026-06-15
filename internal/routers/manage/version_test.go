package manage

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gocronx-team/gocron/internal/modules/app"
)

func TestVersion_ReturnsAppVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	original := app.AppVersion
	app.AppVersion = "9.9.9"
	defer func() { app.AppVersion = original }()

	r := gin.New()
	r.GET("/api/system/version", Version)

	req := httptest.NewRequest(http.MethodGet, "/api/system/version", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var env map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	data, _ := env["data"].(map[string]any)
	if data["version"] != "9.9.9" {
		t.Fatalf("version = %v, want 9.9.9", data["version"])
	}
}
