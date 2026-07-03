package parser

import (
	"strings"
	"testing"
)

func TestParse_ValidGetRequest(t *testing.T) {
	raw := "GET /index.html HTTP/1.1\r\nHost: localhost\r\nConnection: close\r\n\r\n"
	req, err := Parse(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.HTTPMethod != "GET" {
		t.Errorf("expected GET, got %s", req.HTTPMethod)
	}
	if req.Path != "/index.html" {
		t.Errorf("expected /index.html, got %s", req.Path)
	}
	if req.HTTPVersion != "HTTP/1.1" {
		t.Errorf("expected HTTP/1.1, got %s", req.HTTPVersion)
	}
	if req.Headers["Host"] != "localhost" {
		t.Errorf("expected Host: localhost, got %s", req.Headers["Host"])
	}
}

func TestParse_WithBody(t *testing.T) {
	body := "hello"
	raw := "POST /upload HTTP/1.1\r\nContent-Length: 5\r\n\r\n" + body
	req, err := Parse(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(req.Body) != body {
		t.Errorf("expected body %q, got %q", body, string(req.Body))
	}
}

func TestParse_MalformedRequestLine(t *testing.T) {
	raw := "BADREQUEST\r\n\r\n"
	_, err := Parse(strings.NewReader(raw))
	if err == nil {
		t.Fatal("expected error for malformed request line, got nil")
	}
}

func TestParse_MalformedHeader(t *testing.T) {
	raw := "GET / HTTP/1.1\r\nBadHeader\r\n\r\n"
	_, err := Parse(strings.NewReader(raw))
	if err == nil {
		t.Fatal("expected error for malformed header, got nil")
	}
}

func TestParse_InvalidContentLength(t *testing.T) {
	raw := "POST / HTTP/1.1\r\nContent-Length: abc\r\n\r\n"
	_, err := Parse(strings.NewReader(raw))
	if err == nil {
		t.Fatal("expected error for invalid Content-Length, got nil")
	}
}
