package service

import (
	"errors"
	"strings"
)

// service-слой (мозг)

// 3. TaskService — слой бизнес-логики (есть ссылка на TaskRepo) --
// -- то есть вызывает методы репозитория, добавляя логику (валидация, проверка данных)
// то есть решается, что делать дальше
// этот слой не знает, как работает база — он использует TaskRepo для доступа к данным

type TaskService struct {
	repo *TaskRepo
}
 
// конструктор NewService - связывает сервис и репозиторий
func NewService(r *TaskRepo) *TaskService {
	return &TaskService{repo: r}
}

// CreateTask - создает новую задачу (с проверкой что она не пустя)
func (s *TaskService) CreateTask(text string) (*TaskStruct, error) {
	if strings.TrimSpace(text) == "" {
		return nil, errors.New("task is empty")
	}
	task := &TaskStruct{Task: text}
	if err := s.repo.Create(task); err != nil {
		return nil, err
	}
	return task, nil
}

// GetTasks - возвращает все задачи
func (s *TaskService) GetTasks() ([]TaskStruct, error) {
	var tasks []TaskStruct
	if err := s.repo.GetAll(&tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}


func (s *TaskService) UpdateTask(id uint, text string) (*TaskStruct, error) {
	var task TaskStruct
	if err := s.repo.GetByID(&task, id); err != nil {
		return nil, errors.New("task not found") // или просто err
	}

	task.Task = strings.TrimSpace(text)
	if task.Task == "" {
		return nil, errors.New("task is empty")
	}

	if err := s.repo.Update(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *TaskService) DeleteTask(id uint) error {
	var task TaskStruct
	if err := s.repo.GetByID(&task, id); err != nil {
		return errors.New("task not fount")
	}
	return s.repo.Delete(&task, id)
}