package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

const tasksFile = "tasks.json"

func addTask(tasks []Task, title string) ([]Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return tasks, errors.New("The task title cannot be empty")
	}
	nextID := 0

	for _, t := range tasks {
		if t.ID > nextID {
			nextID = t.ID
		}
	}

	nextID++

	newTask := Task{
		ID:        nextID,
		Title:     title,
		Completed: false,
	}

	return append(tasks, newTask), nil
}

func completeTask(tasks []Task, id int) error {
	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Completed = true
			return nil
		}
	}
	return fmt.Errorf("Task with ID=%d was not found", id)
}

func deleteTask(tasks []Task, id int) ([]Task, error) {
	for i, t := range tasks {
		if t.ID == id {
			newTask := append(tasks[:i], tasks[i+1:]...)
			return newTask, nil
		}
	}

	return tasks, fmt.Errorf("Task with ID=%d was not found", id)
}

func printTasks(tasks []Task) {
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

func taskToJSON(tasks []Task) (string, error) {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return "", fmt.Errorf("Error converting to JSON: %w", err)
	}

	return string(data), nil
}

func saveTasksToFile(tasks []Task) error {
	jsonStr, err := taskToJSON(tasks)
	if err != nil {
		return fmt.Errorf("Error creating JSON: %w", err)
	}

	if err := os.WriteFile("tasks.json", []byte(jsonStr), 0644); err != nil {
		return fmt.Errorf("File writing error: %w", err)
	}
	return nil
}

func loadTasksFromFile() ([]Task, error) {
	data, err := os.ReadFile(tasksFile)
	if err != nil {
		// xetanin fayl yoxdur tipinde olub olmadigini yoxlayir
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, fmt.Errorf("file reading error: %w", err)
	}

	if len(data) == 0 {
		return []Task{}, nil
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("JSON parse error: %w", err)
	}

	return tasks, nil
}

func main() {

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	tasks, err := loadTasksFromFile()
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

		newTask, err := addTask(tasks, os.Args[2])
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		tasks = newTask

		if err := saveTasksToFile(tasks); err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Printf("Task added: %s\n", os.Args[2])

	case "list":
		printTasks(tasks)

	case "json":
		out, err := taskToJSON(tasks)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Println(out)

	case "save":
		if err := saveTasksToFile(tasks); err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Printf("Tasks were written to the %s file.\n", tasksFile)

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

		if err := completeTask(tasks, id); err != nil {
			fmt.Println("Error:", err)
			return
		}
		if err := saveTasksToFile(tasks); err != nil {
			fmt.Println("Error: ", err)
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

		newTask, err := deleteTask(tasks, id)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		tasks = newTask

		if err := saveTasksToFile(tasks); err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Printf("ID=%d deleted.\n", id)

	default:
		fmt.Printf("Unknown command: %q\n\n", command)
		printUsage()
	}
}
