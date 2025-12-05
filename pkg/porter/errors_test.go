package porter

import (
	"errors"
	"testing"
)

func TestIsNotFound(t *testing.T) {
	err := &Error{StatusCode: 404}
	if !IsNotFound(err) {
		t.Error("expected IsNotFound to return true for 404")
	}
	if IsNotFound(&Error{StatusCode: 200}) {
		t.Error("expected IsNotFound to return false for 200")
	}
}

func TestIsUnauthorized(t *testing.T) {
	err := &Error{StatusCode: 401}
	if !IsUnauthorized(err) {
		t.Error("expected IsUnauthorized to return true for 401")
	}
}

func TestIsRateLimited(t *testing.T) {
	err := &Error{StatusCode: 429}
	if !IsRateLimited(err) {
		t.Error("expected IsRateLimited to return true for 429")
	}
}

func TestIsServerError(t *testing.T) {
	tests := []struct {
		code int
		want bool
	}{
		{500, true},
		{502, true},
		{503, true},
		{400, false},
		{404, false},
		{200, false},
	}
	for _, tt := range tests {
		err := &Error{StatusCode: tt.code}
		if got := IsServerError(err); got != tt.want {
			t.Errorf("IsServerError(%d) = %v, want %v", tt.code, got, tt.want)
		}
	}
}

func TestStatusCode(t *testing.T) {
	err := &Error{StatusCode: 404}
	if got := StatusCode(err); got != 404 {
		t.Errorf("StatusCode() = %d, want 404", got)
	}
	if got := StatusCode(errors.New("other error")); got != 0 {
		t.Errorf("StatusCode() = %d, want 0", got)
	}
}

func TestErrorMessage(t *testing.T) {
	err := &Error{StatusCode: 404, Message: "not found"}
	if got := ErrorMessage(err); got != "not found" {
		t.Errorf("ErrorMessage() = %q, want %q", got, "not found")
	}
}
