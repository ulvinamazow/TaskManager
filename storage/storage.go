package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ulvinamazow/taskManager/model"
)

const TasksFile = "tasks.json"

func Save(tasks []model.Task) error {

	jsonStr, err := tasksToJSON(tasks)
	if err != nil {
		return fmt.Errorf("JSON creating error: %w", err)
	}

	if err := os.WriteFile(TasksFile, []byte(jsonStr), 0644); err != nil {
		return fmt.Errorf("file write error: %w", err)
	}

	return nil
}

func Load() ([]model.Task, error) {
	data, err := os.ReadFile(TasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []model.Task{}, nil
		}
		return nil, fmt.Errorf("file reading error: %w", err)
	}

	if len(data) == 0 {
		return []model.Task{}, nil
	}

	var tasks []model.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("JSON parse error: %w", err)
	}

	return tasks, nil
}

func tasksToJSON(tasks []model.Task) (string, error) {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return "", fmt.Errorf("Error converting to JSON: %w", err)
	}
	return string(data), nil
}

func ToJSON(tasks []model.Task) (string, error) {
	return tasksToJSON(tasks)
}
