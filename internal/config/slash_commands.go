package config

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// LoadSlashCommands reads .agent-layer/slash-commands/*.md in lexicographic order.
func LoadSlashCommands(dir string) ([]SlashCommand, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("missing slash commands directory %s: %w", dir, err)
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".md") {
			names = append(names, name)
		}
	}

	sort.Strings(names)

	commands := make([]SlashCommand, 0, len(names))
	for _, name := range names {
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed reading slash command %s: %w", path, err)
		}
		data = bytes.TrimPrefix(data, utf8BOM)
		description, body, err := parseSlashCommand(string(data))
		if err != nil {
			return nil, fmt.Errorf("invalid slash command %s: %w", path, err)
		}
		commands = append(commands, SlashCommand{
			Name:        strings.TrimSuffix(name, ".md"),
			Description: description,
			Body:        body,
			SourcePath:  path,
		})
	}

	return commands, nil
}

func parseSlashCommand(content string) (string, string, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	if !scanner.Scan() {
		return "", "", fmt.Errorf("missing content")
	}
	if strings.TrimSpace(scanner.Text()) != "---" {
		return "", "", fmt.Errorf("missing front matter")
	}

	var fmLines []string
	foundEnd := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			foundEnd = true
			break
		}
		fmLines = append(fmLines, line)
	}
	if !foundEnd {
		return "", "", fmt.Errorf("unterminated front matter")
	}

	var bodyBuilder strings.Builder
	for scanner.Scan() {
		bodyBuilder.WriteString(scanner.Text())
		bodyBuilder.WriteString("\n")
	}
	if err := scanner.Err(); err != nil {
		return "", "", fmt.Errorf("failed reading content: %w", err)
	}

	body := strings.TrimPrefix(bodyBuilder.String(), "\n")
	body = strings.TrimRight(body, "\n")

	description, err := parseDescription(fmLines)
	if err != nil {
		return "", "", err
	}

	return description, body, nil
}

func parseDescription(lines []string) (string, error) {
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if !strings.HasPrefix(line, "description:") {
			continue
		}
		value := strings.TrimSpace(strings.TrimPrefix(line, "description:"))
		if value == "" {
			return "", fmt.Errorf("description is empty")
		}
		if value == ">-" || value == ">" || value == "|" || value == "|+" || value == "|-" {
			var parts []string
			for j := i + 1; j < len(lines); j++ {
				if strings.HasPrefix(lines[j], "  ") {
					parts = append(parts, strings.TrimSpace(strings.TrimPrefix(lines[j], "  ")))
					continue
				}
				break
			}
			description := strings.TrimSpace(strings.Join(parts, " "))
			if description == "" {
				return "", fmt.Errorf("description is empty")
			}
			return description, nil
		}
		return strings.Trim(value, "\""), nil
	}

	return "", fmt.Errorf("missing description in front matter")
}
