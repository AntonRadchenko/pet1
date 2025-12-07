package taskService

import (
	"errors"
	"strings"
)

// 3. service-слой (мозг)

// TaskService — слой бизнес-логики (есть ссылка на TaskRepo) --
// -- то есть вызывает методы репозитория, добавляя логику (валидация, проверка данных)
// то есть решается, что делать дальше
// этот слой не знает, как работает база — он использует TaskRepo для доступа к данным

// структура параметров метода CreateTask
type CreateTaskParams struct {
	Task   string
	IsDone *bool
	UserId uint
}

// структура параметров метода UpdateTask
type UpdateTaskParams struct {
	Task   *string
	IsDone *bool
	UserId *uint
}

// бизнес-модель, которую возвращает сервис
type Task struct {
	ID     uint
	Task   string
	IsDone *bool
	UserId uint
}

type TaskService struct {
	repo TaskRepoInterface // используем интерфейс
}

// конструктор NewTaskService - связывает сервис и репозиторий
func NewTaskService(r TaskRepoInterface) *TaskService {
	return &TaskService{repo: r}
}

// CreateTask - создает новую задачу (с проверкой что она не пустя)
func (s *TaskService) CreateTask(params CreateTaskParams) (*Task, error) {
	// проверка на пустой тип задачи
	if strings.TrimSpace(params.Task) == "" {
		return nil, errors.New("task is empty")
	}

	if params.UserId == 0 {
		return nil, errors.New("user_id is required")
	}

	// если isDone не был передан пользователем, то он будет по умолчанию false
	isDone := false
	if params.IsDone != nil { // если done не пустой, то есть был передан в бади
		isDone = *params.IsDone // то обновляем его по указателю
	}

	// создаем бд-модель
	dbTask := &TaskStruct{
		Task:   params.Task,
		IsDone: isDone,
		UserId: params.UserId,
	}

	createdTask, err := s.repo.Create(dbTask) // передаем данные в репозиторий
	if err != nil {
		return nil, err
	}

	// маппим бд-модель в бизнес-модель
	return &Task{
		ID:     createdTask.ID,
		Task:   createdTask.Task,
		IsDone: &createdTask.IsDone,
		UserId: createdTask.UserId,
	}, nil
}

// GetTasks - возвращает все задачи
func (s *TaskService) GetTasks() ([]Task, error) {
	dbTasks, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// маппим бд-модель в бизнес-модель
	tasks := make([]Task, 0, len(dbTasks))
	for _, dbTask := range dbTasks {
		tasks = append(tasks, Task{
			ID:     dbTask.ID,
			Task:   dbTask.Task,
			IsDone: &dbTask.IsDone,
			UserId: dbTask.UserId,
		})
	}
	return tasks, nil
}

func (s *TaskService) UpdateTask(id uint, params UpdateTaskParams) (*Task, error) {
	dbTask, err := s.repo.GetByID(id)
	if err != nil || dbTask.ID == 0 {
		return nil, errors.New("task not found")
	}

	updated := false

	if params.Task != nil {
		// проверяем что таска не nil перед ее обновлением
		if strings.TrimSpace(*params.Task) == "" {
			return nil, errors.New("task is empty")
		}
		// обновляем
		dbTask.Task = *params.Task // обновляем таску если она была передана для обновления
		updated = true
	}

	if params.IsDone != nil {
		// обновляем
		dbTask.IsDone = *params.IsDone // обновляем флаг is_done если он был передан для обновления
		updated = true
	}

	if params.UserId != nil {
		if *params.UserId == 0 {
			return nil, errors.New("user_id cannot be 0")
		}
		dbTask.UserId = *params.UserId
		updated = true
	}

	if !updated {
		return nil, errors.New("no fields to update")
	}

	// обновляем задачу
	updatedTask, err := s.repo.Update(&dbTask)
	if err != nil {
		return nil, err
	}

	// маппим бд-модель в бизнес-модель
	return &Task{
		ID:     updatedTask.ID,
		Task:   updatedTask.Task,
		IsDone: &updatedTask.IsDone,
		UserId: updatedTask.UserId,
	}, nil
}

func (s *TaskService) DeleteTask(id uint) error {
	// ищем задачу по ID
	task, err := s.repo.GetByID(id)
	if err != nil || task.ID == 0 {
		return errors.New("task not found")
	}
	// удаляем задачу
	err = s.repo.Delete(&task)
	if err != nil {
		return err
	}
	return nil
}
