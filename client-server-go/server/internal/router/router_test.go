package router

import (
	"testing"

	"github.com/SaiVikrantG/server/internal/models"
	"github.com/SaiVikrantG/server/internal/response"
)

type mockResponseWriter struct {
	statusCode int
	headers    map[string]string
	body       []byte
}

func newMockResponseWriter() *mockResponseWriter {
	return &mockResponseWriter{headers: make(map[string]string)}
}

func (m *mockResponseWriter) Header() map[string]string {
	return m.headers
}

func (m *mockResponseWriter) Write(statusCode int, body []byte) (int, error) {
	m.statusCode = statusCode
	m.body = body
	return len(body), nil
}

type mockHandler struct {
	called bool
}

func (h *mockHandler) ServeHTTP(w response.ResponseWriter, r *models.Request) {
	h.called = true
	w.Write(200, []byte("ok"))
}

func TestRouter_MatchesRoute(t *testing.T) {
	r := NewRouter()
	h := &mockHandler{}
	r.Handle("GET", "/hello", h)

	w := newMockResponseWriter()
	req := &models.Request{HTTPMethod: "GET", Path: "/hello"}
	r.ServeHTTP(w, req)

	if !h.called {
		t.Error("expected handler to be called")
	}
	if w.statusCode != 200 {
		t.Errorf("expected status 200, got %d", w.statusCode)
	}
}

func TestRouter_NoMatchReturns404(t *testing.T) {
	r := NewRouter()

	w := newMockResponseWriter()
	req := &models.Request{HTTPMethod: "GET", Path: "/notregistered"}
	r.ServeHTTP(w, req)

	if w.statusCode != 404 {
		t.Errorf("expected status 404, got %d", w.statusCode)
	}
}

func TestRouter_MethodMismatchReturns404(t *testing.T) {
	r := NewRouter()
	r.Handle("GET", "/hello", &mockHandler{})

	w := newMockResponseWriter()
	req := &models.Request{HTTPMethod: "POST", Path: "/hello"}
	r.ServeHTTP(w, req)

	if w.statusCode != 404 {
		t.Errorf("expected status 404, got %d", w.statusCode)
	}
}
