package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ulvinamazow/taskManager/model"
	"github.com/ulvinamazow/taskManager/service"
)

func printTasks(tasks []model.Task) {
	if len(tasks) == 0 {
		fmt.Println("There are no tasks.")
		return
	}

	for _, task := range tasks {
		fmt.Printf("ID=%d | %s | completed=%t\n",
			task.ID, task.Title, task.Completed)
	}
}

func printUsage() {
	fmt.Println("Usage: task <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  add \"title\"      — add a new task")
	fmt.Println("  list               — show all tasks")
	fmt.Println("  complete <id>      — mark the task as completed")
	fmt.Println("  delete <id>        — delete task")
}

func parseID(s string) (int, error) {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("unreliable ID %q: %w", s, err)
	}
	return id, nil
}

func main() {

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	service, err := service.NewService()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task add \"title\"")
			return
		}

		title := strings.Join(os.Args[2:], " ")

		if err := service.Add(title); err != nil {
			fmt.Println("Error: ", err)
			return
		}

		fmt.Printf("Task added: %s\n", os.Args[2])

	case "list":
		printTasks(service.List())

	case "json":
		out, err := service.ToJSON()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println(out)

	case "complete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task complete <id>")
			return
		}

		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if err := service.Complete(id); err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("ID=%d completed.\n", id)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: task delete <id>")
			return
		}

		id, err := parseID(os.Args[2])
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if err := service.Delete(id); err != nil {
			fmt.Println("Error: ", err)
			return
		}

		fmt.Printf("ID=%d deleted.\n", id)

	default:
		fmt.Printf("Unknown command: %q\n\n", command)
		printUsage()
	}
}
