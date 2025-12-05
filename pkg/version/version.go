package version

import "strings"

// Version follows semantic versioning (vMAJOR.MINOR.PATCH).
const Version = "v0.2.0"

// UserAgent returns the default user agent string shared by the SDK.
func UserAgent() string {
	return "port-go-sdk/" + strings.TrimPrefix(Version, "v")
}
