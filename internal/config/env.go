package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadEnv reads .agent-layer/.env into a key-value map.
func LoadEnv(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("missing env file %s: %w", path, err)
	}

	env := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		idx := strings.Index(line, "=")
		if idx <= 0 {
			return nil, fmt.Errorf("invalid env file %s:%d: expected KEY=VALUE", path, lineNo)
		}
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}
		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed reading env file %s: %w", path, err)
	}

	return env, nil
}
