package porter

import (
	"errors"
	"fmt"
)

// IsNotFound returns true if the error is a 404 Not Found error.
func IsNotFound(err error) bool {
	var perr *Error
	return errors.As(err, &perr) && perr.StatusCode == 404
}

// IsUnauthorized returns true if the error is a 401 Unauthorized error.
func IsUnauthorized(err error) bool {
	var perr *Error
	return errors.As(err, &perr) && perr.StatusCode == 401
}

// IsForbidden returns true if the error is a 403 Forbidden error.
func IsForbidden(err error) bool {
	var perr *Error
	return errors.As(err, &perr) && perr.StatusCode == 403
}

// IsRateLimited returns true if the error is a 429 Too Many Requests error.
func IsRateLimited(err error) bool {
	var perr *Error
	return errors.As(err, &perr) && perr.StatusCode == 429
}

// IsConflict returns true if the error is a 409 Conflict error.
func IsConflict(err error) bool {
	var perr *Error
	return errors.As(err, &perr) && perr.StatusCode == 409
}

// IsServerError returns true if the error is a 5xx server error.
func IsServerError(err error) bool {
	var perr *Error
	return errors.As(err, &perr) && perr.StatusCode >= 500 && perr.StatusCode < 600
}

// StatusCode returns the HTTP status code from a Port API error, or 0 if not a Port error.
func StatusCode(err error) int {
	var perr *Error
	if errors.As(err, &perr) {
		return perr.StatusCode
	}
	return 0
}

// ErrorMessage returns a user-friendly error message.
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	var perr *Error
	if errors.As(err, &perr) {
		if perr.Message != "" {
			return perr.Message
		}
		return fmt.Sprintf("HTTP %d: %s", perr.StatusCode, httpStatusText(perr.StatusCode))
	}
	// err is guaranteed to be non-nil here due to the nil check above
	return err.Error() //nolint:gocritic // nilValReturn: err is checked for nil above
}
