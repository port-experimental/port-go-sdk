package porter

import "testing"

func TestErrorString(t *testing.T) {
	err := &Error{StatusCode: 404, Message: "not found"}
	got := err.Error()
	if want := "port api: 404 Not Found: not found"; got != want {
		t.Fatalf("want %q got %q", want, got)
	}
}
