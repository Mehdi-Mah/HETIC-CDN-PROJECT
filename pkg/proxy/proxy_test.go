package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewProxyHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	handler := NewProxyHandler()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadGateway {
		t.Errorf("expected status %v, got %v", http.StatusBadGateway, resp.StatusCode)
	}
}
