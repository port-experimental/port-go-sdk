package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// loadDotEnvFile best-effort loads environment variables from a .env-style file.
// Missing files are silently ignored, matching godotenv.Load behavior.
func loadDotEnvFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(line[len("export "):])
		}
		idx := strings.IndexRune(line, '=')
		if idx == -1 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		if key == "" {
			continue
		}
		val := strings.TrimSpace(line[idx+1:])
		if len(val) >= 2 && isQuoted(val) {
			val = val[1 : len(val)-1]
		}
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		if err := os.Setenv(key, val); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func isQuoted(s string) bool {
	if len(s) < 2 {
		return false
	}
	start, end := s[0], s[len(s)-1]
	return (start == '"' && end == '"') || (start == '\'' && end == '\'')
}
