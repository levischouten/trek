package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"trek/config"
	"trek/habit"
)

func parseArgs(args []string) (string, []string, map[string]string) {
	if len(args) < 2 {
		fmt.Println("Usage: cli <command> [arguments...] [--flags...]")
		os.Exit(1)
	}

	command := args[1]
	var arguments []string
	flags := make(map[string]string)

	for _, arg := range args[2:] {
		if strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(arg[2:], "=", 2)
			if len(parts) != 2 {
				fmt.Println("Usage: cli <command> [arguments...] [--flags...]")
				os.Exit(1)
			}
			key, value := parts[0], parts[1]
			flags[key] = value
		} else {
			arguments = append(arguments, arg)
		}
	}

	return command, arguments, flags
}

func main() {
	command, arguments, flags := parseArgs(os.Args)

	dir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user config directory:", err)
		return
	}
	filePath := filepath.Join(dir, ".config/trek", "config.ini")

	cfg, err := config.ParseINI(filePath)
	if err != nil {
		fmt.Println("Error parsing INI file:", err)
		return
	}

	switch command {
	case "log":
		err = habit.Init()
		if err != nil {
			fmt.Println("Error initializing habit:", err)
			return
		}

		date, ok := flags["date"]
		if ok {
			_, err = time.Parse("2006-01-02", date)
			if err != nil {
				fmt.Println("Error parsing date:", err)
				os.Exit(1)
			}
		} else {
			date = time.Now().Format("2006-01-02")
		}

		if len(arguments) == 0 {
			reader := bufio.NewReader(os.Stdin)

			for name, h := range cfg.Habits {
				fmt.Printf("Did you '%s' today? (y/n): ", h.Name)
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(strings.ToLower(input))

				if input == "y" {
					err = habit.LogHabit(name, date)
					if err != nil {
						fmt.Println("Error logging Habit:", err)
						os.Exit(1)
					}
				}
			}
		} else {
			for _, name := range arguments {
				_, ok = cfg.Habits[name]

				if !ok {
					fmt.Printf("Habit [%s] does not exist\n", name)
					os.Exit(1)
				}

				err = habit.LogHabit(name, date)
				if err != nil {
					fmt.Println("Error logging Habit:", err)
					os.Exit(1)
				}
			}

			fmt.Println("Habit(s) logged.")
		}

		err = habit.Save()
		if err != nil {
			fmt.Println("Error Marking Habit:", err)
			os.Exit(1)
		}

		os.Exit(0)
	default:
		fmt.Printf("Command \"%s\" not recognized\n", command)
		os.Exit(1)
	}
}
