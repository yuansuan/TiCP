package monitor

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMonitor(t *testing.T) {
	r := gin.Default()
	p := newMonitor(&Config{
		Server: r,
	})
	p.runMetricsServer()

	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	result := w.Result()
	defer result.Body.Close()

	expectedCode := 200
	if result.StatusCode != expectedCode {
		t.Error("request /metrics can't get 200 code")
	}
}
