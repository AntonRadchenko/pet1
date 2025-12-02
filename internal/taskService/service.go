package taskService

import (
	"errors"
	"strings"

	openapi "github.com/AntonRadchenko/WebPet1/openapi"
)

// 3. service-слой (мозг)

// TaskService — слой бизнес-логики (есть ссылка на TaskRepo) --
// -- то есть вызывает методы репозитория, добавляя логику (валидация, проверка данных)
// то есть решается, что делать дальше
// этот слой не знает, как работает база — он использует TaskRepo для доступа к данным

type TaskServiceInterface interface {
	CreateTask(taskRequest openapi.PostTasksJSONRequestBody) (*TaskStruct, error)
	GetTasks() ([]TaskStruct, error)
	UpdateTask(id uint, taskRequest openapi.PatchTasksIdJSONRequestBody) (*TaskStruct, error)
	DeleteTask(id uint) error
}

type TaskService struct {
	repo TaskRepoInterface // используем интерфейс
}
	
// конструктор NewService - связывает сервис и репозиторий
func NewService(r TaskRepoInterface) *TaskService {
	return &TaskService{repo: r}
}

// CreateTask - создает новую задачу (с проверкой что она не пустя)
func (s *TaskService) CreateTask(taskRequest openapi.PostTasksJSONRequestBody) (*TaskStruct, error) {
	// проверка на пустой тип задачи
	if strings.TrimSpace(taskRequest.Task) == "" {
		return nil, errors.New("task is empty")
	}

	// если isDone не был передан пользователем, то он будет по умолчанию false
	isDone := false
	if taskRequest.IsDone != nil { // если done не пустой, то есть был передан в бади
		isDone = *taskRequest.IsDone // то обновляем его по указателю
	}

	task := &TaskStruct{
		Task: taskRequest.Task,
		IsDone: isDone,
	}

	createdTask, err := s.repo.Create(task) // передаем данные в репозиторий
	if err != nil {
		return nil, err
	}
	return createdTask, nil // передаем объект задачи (который вернул репозиторий) обратно в хендлер 
}

// GetTasks - возвращает все задачи
func (s *TaskService) GetTasks() ([]TaskStruct, error) {
	tasks, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) UpdateTask(id uint, taskRequest openapi.PatchTasksIdJSONRequestBody) (*TaskStruct, error) {
	// проверка на пустой тип задачи
    if strings.TrimSpace(taskRequest.Task) == "" {
        return nil, errors.New("task is empty") // Если задача пустая, возвращаем ошибку
    }

	task, err := s.repo.GetByID(id); 
	if err != nil {
		return nil, err 
	}

	// обновляем поля задачи
	task.Task = taskRequest.Task
	if taskRequest.IsDone != nil {
		task.IsDone = *taskRequest.IsDone // обновляем IsDone если он был передан для обновления
	}

	// обновляем задачу
	updatedTask, err := s.repo.Update(&task)
	if err != nil {
		return nil, err
	}
	return updatedTask, nil
}

func (s *TaskService) DeleteTask(id uint) error {
	// ищем задачу по ID
	task, err := s.repo.GetByID(id) 
	if err != nil {
		return errors.New("task not found")
	}
	// удаляем задачу
	err = s.repo.Delete(&task, id)
	if err != nil {
		return err
	}
	return nil
}

