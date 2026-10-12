package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ulvinamazow/taskManager/model"
)

type FileStorage struct {
	path string
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{path: path}
}

func (fs *FileStorage) Save(tasks []model.Task) error {

	jsonStr, err := fs.toJSON(tasks)
	if err != nil {
		return fmt.Errorf("JSON creating error: %w", err)
	}

	if err := os.WriteFile(fs.path, []byte(jsonStr), 0644); err != nil {
		return fmt.Errorf("file write error: %w", err)
	}

	return nil
}

func (fs *FileStorage) Load() ([]model.Task, error) {
	data, err := os.ReadFile(fs.path)
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

func (fs *FileStorage) toJSON(tasks []model.Task) (string, error) {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return "", fmt.Errorf("Json conversion error: %w", err)
	}
	return string(data), nil
}

func (fs *FileStorage) ToJSON(tasks []model.Task) (string, error) {
	return fs.toJSON(tasks)
}
