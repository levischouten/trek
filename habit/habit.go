package habit

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Habit struct {
	Name        string
	Description string
	Active      bool
}

var logs map[string][]string

func Init() error {
	logs = make(map[string][]string)
	dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user config directory:", err)
		os.Exit(1)
	}

	filePath := filepath.Join(dir, ".config/trek", "log")

	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, " ", 2)

		if len(parts) == 2 {
			date := parts[0]
			habits := strings.Split(parts[1], " ")

			logs[date] = habits
		}
	}

	return scanner.Err()
}

func Save() error {
	dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user config directory:", err)
		os.Exit(1)
	}

	filePath := filepath.Join(dir, ".config/trek", "log")
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file for writing: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for date, loggedHabits := range logs {
		_, err := writer.WriteString(fmt.Sprintf("%s %s\n", date, strings.Join(loggedHabits, " ")))
		if err != nil {
			return fmt.Errorf("failed to write to log file: %w", err)
		}
	}
	return writer.Flush()
}

func LogHabit(habit string, date string) error {
	habits, ok := logs[date]
	if !ok {
		habits = []string{}
	}

	found := false

	for _, h := range habits {
		if h == habit {
			found = true
			break
		}
	}

	if !found {
		logs[date] = append(logs[date], habit)
	}

	return nil
}
