package service

import (
	"errors"
	"fmt"

	"github.com/ulvinamazow/taskManager/model"
)

type Storage interface {
	Save(tasks []model.Task) error
	Load() ([]model.Task, error)
	ToJSON(tasks []model.Task) (string, error)
}

type Service struct {
	storage Storage
	tasks   []model.Task
}

func New(storage Storage) (*Service, error) {
	tasks, err := storage.Load()
	if err != nil {
		return nil, fmt.Errorf("service can not created: %w", err)
	}

	return &Service{
		storage: storage,
		tasks:   tasks,
	}, nil
}

func (service *Service) Add(title string) error {
	if title == "" {
		return errors.New("The tasks title con not be empty")
	}

	nextID := 0

	for _, t := range service.tasks {
		if t.ID > nextID {
			nextID = t.ID
		}
	}
	nextID++

	newTask := model.Task{
		ID:        nextID,
		Title:     title,
		Completed: false,
	}

	service.tasks = append(service.tasks, newTask)

	if err := service.storage.Save(service.tasks); err != nil {
		return fmt.Errorf("The task could not be saved: %w", err)
	}
	return nil
}

func (service *Service) Complete(id int) error {
	for i, t := range service.tasks {
		if t.ID == id {
			service.tasks[i].Completed = true
			if err := service.storage.Save(service.tasks); err != nil {
				return fmt.Errorf("The task could not be saved: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("Task with ID=%d was not found", id)
}

func (service *Service) Delete(id int) error {
	for i, t := range service.tasks {
		if t.ID == id {
			// Slice-dan çıxar.
			service.tasks = append(service.tasks[:i], service.tasks[i+1:]...)
			if err := service.storage.Save(service.tasks); err != nil {
				return fmt.Errorf("The task could not be saved: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("Task with ID=%d was not found", id)
}

func (service *Service) List() []model.Task {
	return service.tasks
}

func (service *Service) ToJSON() (string, error) {
	return service.storage.ToJSON(service.tasks)
}
