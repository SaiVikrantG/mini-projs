package models

type Request struct {
	HTTPMethod  string
	Path        string
	HTTPVersion string
	Headers     map[string]string
	Body        []byte
}
