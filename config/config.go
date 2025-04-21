package config

import (
	"fmt"
	"os"
	"strings"

	"trek/habit"
)

type Config struct {
	Habits map[string]habit.Habit
}

func ParseINI(filePath string) (Config, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open file: %w", err)
	}

	config := Config{
		Habits: make(map[string]habit.Habit),
	}
	lines := strings.Split(string(content), "\n")
	var current string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			current = strings.TrimPrefix(strings.TrimSuffix(line, "]"), "[")
			continue
		}

		if current == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			if strings.HasPrefix(current, "habit.") {
				name := strings.TrimPrefix(current, "habit.")
				h, ok := config.Habits[name]
				if !ok {
					h = habit.Habit{}
				}

				switch key {
				case "name":
					h.Name = value
				case "description":
					h.Description = strings.Trim(value, "\"")
				case "active":
					h.Active = strings.ToLower(value) == "true"
				}

				config.Habits[name] = h
			}
		}
	}

	return config, nil
}
