package service

import (
	"errors"
	"fmt"

	"github.com/ulvinamazow/taskManager/model"
	"github.com/ulvinamazow/taskManager/storage"
)

type Service struct {
	tasks []model.Task
}

func NewService() (*Service, error) {
	tasks, err := storage.Load()
	if err != nil {
		return nil, fmt.Errorf("service can not created: %w", err)
	}

	return &Service{tasks: tasks}, nil
}

func (service *Service) Add(title string) error {
	if title == "" {
		return errors.New("The tasks title con not be empty")
	}

	nextID := 0

	newTask := model.Task{
		ID:        nextID,
		Title:     title,
		Completed: false,
	}
	service.tasks = append(service.tasks, newTask)

	if err := storage.Save(service.tasks); err != nil {
		return fmt.Errorf("The task could not be saved: %w", err)
	}
	return nil
}

func (s *Service) Complete(id int) error {
	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks[i].Completed = true
			if err := storage.Save(s.tasks); err != nil {
				return fmt.Errorf("The task could not be saved: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("Task with ID=%d was not found", id)
}

func (s *Service) Delete(id int) error {
	for i, t := range s.tasks {
		if t.ID == id {
			// Slice-dan çıxar.
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			if err := storage.Save(s.tasks); err != nil {
				return fmt.Errorf("The task could not be saved: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("Task with ID=%d was not found", id)
}

func (s *Service) List() []model.Task {
	return s.tasks
}

func (s *Service) ToJSON() (string, error) {
	return storage.ToJSON(s.tasks)
}
