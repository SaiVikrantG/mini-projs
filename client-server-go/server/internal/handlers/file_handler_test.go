package handlers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/SaiVikrantG/server/internal/models"
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

func TestFileHandler_ServesFile(t *testing.T) {
	dir := t.TempDir()
	content := []byte("<h1>Hello</h1>")
	os.WriteFile(filepath.Join(dir, "index.html"), content, 0644)

	h := &FileHandler{RootDir: dir}
	w := newMockResponseWriter()
	req := &models.Request{HTTPMethod: "GET", Path: "/index.html"}

	h.ServeHTTP(w, req)

	if w.statusCode != 200 {
		t.Errorf("expected status 200, got %d", w.statusCode)
	}
	if string(w.body) != string(content) {
		t.Errorf("expected body %q, got %q", content, w.body)
	}
	if w.headers["Content-Type"] != "text/html" {
		t.Errorf("expected Content-Type text/html, got %s", w.headers["Content-Type"])
	}
}

func TestFileHandler_NotFound(t *testing.T) {
	dir := t.TempDir()

	h := &FileHandler{RootDir: dir}
	w := newMockResponseWriter()
	req := &models.Request{HTTPMethod: "GET", Path: "/missing.html"}

	h.ServeHTTP(w, req)

	if w.statusCode != 404 {
		t.Errorf("expected status 404, got %d", w.statusCode)
	}
}

func TestFileHandler_PathTraversalBlocked(t *testing.T) {
	dir := t.TempDir()
	secret := filepath.Join(filepath.Dir(dir), "secret.txt")
	os.WriteFile(secret, []byte("secret"), 0644)

	h := &FileHandler{RootDir: dir}
	w := newMockResponseWriter()
	req := &models.Request{HTTPMethod: "GET", Path: "/../secret.txt"}

	h.ServeHTTP(w, req)

	if w.statusCode == 200 {
		t.Error("path traversal should not return 200")
	}
}
