package porter

import (
	"fmt"
)

// Error wraps Port API error responses with HTTP metadata.
type Error struct {
	StatusCode int
	Message    string
	Body       []byte
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Message != "" {
		return fmt.Sprintf("port api: %d %s: %s", e.StatusCode, httpStatusText(e.StatusCode), e.Message)
	}
	return fmt.Sprintf("port api: %d %s", e.StatusCode, httpStatusText(e.StatusCode))
}

// Unwrap allows the error to work with errors.Is and errors.As.
func (e *Error) Unwrap() error {
	return nil
}

func httpStatusText(code int) string {
	if text := statusText[code]; text != "" {
		return text
	}
	return "status"
}

var statusText = map[int]string{
	400: "Bad Request",
	401: "Unauthorized",
	403: "Forbidden",
	404: "Not Found",
	409: "Conflict",
	422: "Unprocessable Entity",
	429: "Too Many Requests",
	500: "Internal Server Error",
	502: "Bad Gateway",
	503: "Service Unavailable",
}
