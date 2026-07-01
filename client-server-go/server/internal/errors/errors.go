package errors

import "fmt"

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}

func NewBadRequest(message string) *HTTPError {
	return &HTTPError{StatusCode: 400, Message: message}
}

func NewNotFound(message string) *HTTPError {
	return &HTTPError{StatusCode: 404, Message: message}
}

func NewRequestTooLarge(message string) *HTTPError {
	return &HTTPError{StatusCode: 413, Message: message}
}

func NewHTTPVersionNotSupported(message string) *HTTPError {
	return &HTTPError{StatusCode: 505, Message: message}
}
