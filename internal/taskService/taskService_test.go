package taskService

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateTask(t *testing.T) {
	// вспомогательные функции
	boolPtr := func(b bool) *bool { return &b }

	// создаем слайс структур, в каждой из которых описан тестовый случай
	tests := []struct {
		name      string                                                     // имя теста
		params    CreateTaskParams                                           // входные данные
		want      *Task                                                      // ошидаемая бизнес-модель
		mockSetup func(m *MockTaskRepo, params CreateTaskParams, want *Task) // функция настройки мок-репы
		wantErr   bool                                                       // флаг который говорит ожидать ли ошибку
	}{
		{
			name: "успешное создание задачи",
			params: CreateTaskParams{
				Task:   "Test",
				IsDone: boolPtr(false),
				UserId: 1,
			},
			want: &Task{
				Task:   "Test",
				IsDone: boolPtr(false),
				UserId: 1,
			},
			mockSetup: func(m *MockTaskRepo, params CreateTaskParams, want *Task) {
				// Конвертируем params в TaskStruct для мока
				dbTask := &TaskStruct{
					Task:   params.Task,
					IsDone: *params.IsDone,
					UserId: params.UserId,
				}
				m.On("Create", dbTask).Return(dbTask, nil)
			},
			wantErr: false,
		},
        {
            name: "ошибка при создании в БД",
            params: CreateTaskParams{
                Task:   "Bad task",
                IsDone: boolPtr(false),
                UserId: 1,
            },
            want:    nil,
            wantErr: true,
            mockSetup: func(m *MockTaskRepo, params CreateTaskParams, want *Task) {
                dbTask := &TaskStruct{
                    Task:   params.Task,
                    IsDone: *params.IsDone,
                    UserId: params.UserId,
                }
                m.On("Create", dbTask).Return(&TaskStruct{}, errors.New("db error"))
            },
        },
	}

	for _, tt := range tests { // проходимся по тестам
		t.Run(tt.name, func(t *testing.T) { // запускаем каждый тест
			mockRepo := new(MockTaskRepo)    // создаем для каждого теста мок-репозиоторий
			tt.mockSetup(mockRepo, tt.params, tt.want) // настройка мока

			service := NewTaskService(mockRepo)
			result, err := service.CreateTask(tt.params)

			if tt.wantErr { // если ожидается ошибка, то проверяется что ошибка произошла
				assert.Error(t, err)
				assert.Nil(t, result)
			} else { // а если ошибки НЕ ожидается, то проверяется что ее нет, и что результат соответствует ожидаемому входному значению
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Task, result.Task)
				assert.NotNil(t, result.IsDone)
				assert.Equal(t, *tt.want.IsDone, *result.IsDone)
				assert.Equal(t, tt.want.UserId, result.UserId)
			}

			mockRepo.AssertExpectations(t) // проверяем что все ожидаемые вызовы методов мока были выполнены
		})
	}
}

func TestGetTasks(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(m *MockTaskRepo)
		wantErr   bool
		want []Task
	}{
		{
			name: "успешное получение всех задач",
			mockSetup: func(m *MockTaskRepo) {
				m.On("GetAll").Return([]TaskStruct{
					{Task: "Task 1", IsDone: true, UserId: 1},
					{Task: "Task 2", IsDone: false, UserId: 2},
				}, nil)
			},
			wantErr: false,
			want: []Task{
				{Task: "Task 1", IsDone: &[]bool{true}[0], UserId: 1},
				{Task: "Task 2", IsDone: &[]bool{false}[0], UserId: 2},
			},
		},
		{
			name: "ошибка при получении задач",
			mockSetup: func(m *MockTaskRepo) {
				m.On("GetAll").Return(nil, errors.New("db error"))
			},
			wantErr:   true,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepo)
			tt.mockSetup(mockRepo)

			service := NewTaskService(mockRepo)
			result, err := service.GetTasks()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Сравниваем длину слайсов, чтобы убедиться что они одинаковые
				assert.Equal(t, len(tt.want), len(result))

				// Если слайсы не пустые, то проходим по ним и сравниваем только важные поля
				for i := range result {
					assert.Equal(t, tt.want[i].Task, result[i].Task)
					assert.Equal(t, *tt.want[i].IsDone, *result[i].IsDone)
					assert.Equal(t, tt.want[i].UserId, result[i].UserId)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateTask(t *testing.T) {
	// Вспомогательные функции
	boolPtr := func(b bool) *bool { return &b }
	stringPtr := func(s string) *string { return &s }
	uintPtr := func(u uint) *uint { return &u }

	tests := []struct {
		name      string
		id        uint
		params    UpdateTaskParams
		want      *Task
		wantErr   bool
		mockSetup func(m *MockTaskRepo, id uint, params UpdateTaskParams, want *Task)
	}{
		{
			name: "успешное обновление задачи",
			id:   1,
			params: UpdateTaskParams{
				Task:   stringPtr("Updated task"),
				IsDone: boolPtr(true),
				UserId: uintPtr(2),
			},
			want: &Task{
				ID:     1,
				Task:   "Updated task",
				IsDone: boolPtr(true),
				UserId: 2,
			},
			wantErr: false,
			mockSetup: func(m *MockTaskRepo, id uint, params UpdateTaskParams, want *Task) {
				// 1. Существующая задача в БД (для GetByID)
				existingTask := TaskStruct{
					ID:     id,
					Task:   "Old task",
					IsDone: false,
					UserId: 1,
				}
				m.On("GetByID", id).Return(existingTask, nil)

				// 2. Обновлённая задача (для Update)
				updatedTask := &TaskStruct{
					ID:     id,
					Task:   *params.Task,
					IsDone: *params.IsDone,
					UserId: *params.UserId,
				}
				m.On("Update", mock.Anything).Return(updatedTask, nil)
			},
		},
		{
			name: "обновление только текста",
			id:   2,
			params: UpdateTaskParams{
				Task: stringPtr("New text"),
			},
			want: &Task{
				ID:     2,
				Task:   "New text",
				IsDone: boolPtr(false), // старое значение
				UserId: 1,              // старое значение
			},
			wantErr: false,
			mockSetup: func(m *MockTaskRepo, id uint, params UpdateTaskParams, want *Task) {
				existingTask := TaskStruct{
					ID:     id,
					Task:   "Old task",
					IsDone: false,
					UserId: 1,
				}
				m.On("GetByID", id).Return(existingTask, nil)

				updatedTask := &TaskStruct{
					ID:     id,
					Task:   *params.Task,
					IsDone: false, // не меняли
					UserId: 1,     // не меняли
				}
				m.On("Update", mock.Anything).Return(updatedTask, nil)
			},
		},

		{
			name: "обновление только статуса",
			id:   6,
			params: UpdateTaskParams{
				IsDone: boolPtr(true),
				// Task и UserID - nil
			},
			want: &Task{
				ID:     6,
				Task:   "Existing task", // старое значение
				IsDone: boolPtr(true),
				UserId: 1, // старое значение
			},
			wantErr: false,
			mockSetup: func(m *MockTaskRepo, id uint, params UpdateTaskParams, want *Task) {
				existingTask := TaskStruct{
					ID:     id,
					Task:   "Existing task",
					IsDone: false,
					UserId: 1,
				}
				m.On("GetByID", id).Return(existingTask, nil)

				updatedTask := &TaskStruct{
					ID:     id,
					Task:   "Existing task", // не меняли
					IsDone: true,            // обновили
					UserId: 1,               // не меняли
				}
				m.On("Update", mock.Anything).Return(updatedTask, nil)
			},
		},
		{
			name: "обновление только UserID",
			id:   7,
			params: UpdateTaskParams{
				UserId: uintPtr(3),
			},
			want: &Task{
				ID:     7,
				Task:   "Existing task",
				IsDone: boolPtr(false),
				UserId: 3,
			},
			wantErr: false,
			mockSetup: func(m *MockTaskRepo, id uint, params UpdateTaskParams, want *Task) {
				existingTask := TaskStruct{
					ID:     id,
					Task:   "Existing task",
					IsDone: false,
					UserId: 1,
				}
				m.On("GetByID", id).Return(existingTask, nil)

				updatedTask := &TaskStruct{
					ID:     id,
					Task:   "Existing task",
					IsDone: false,
					UserId: 3,
				}
				m.On("Update", mock.Anything).Return(updatedTask, nil)
			},
		},
		{
			name: "ошибка - задача не найдена",
			id:   999,
			params: UpdateTaskParams{
				Task: stringPtr("Some task"),
			},
			want:    nil,
			wantErr: true,
			mockSetup: func(m *MockTaskRepo, id uint, params UpdateTaskParams, want *Task) {
				m.On("GetByID", id).Return(TaskStruct{}, gorm.ErrRecordNotFound)
			},
		},

		{
			name: "ошибка - пустая задача",
			id:   3,
			params: UpdateTaskParams{
				Task: stringPtr(""), // пустая строка
			},
			want:    nil,
			wantErr: true,
			mockSetup: func(m *MockTaskRepo, id uint, params UpdateTaskParams, want *Task) {
				existingTask := TaskStruct{
					ID:     id,
					Task:   "Old task",
					IsDone: false,
					UserId: 1,
				}
				m.On("GetByID", id).Return(existingTask, nil)
			},
		},

		{
			name:    "все поля nil - нет полей для обновления",
			id:      4,
			params:  UpdateTaskParams{}, // все поля nil
			want:    nil,
			wantErr: true,
			mockSetup: func(m *MockTaskRepo, id uint, params UpdateTaskParams, want *Task) {
				existingTask := TaskStruct{
					ID:     id,
					Task:   "Existing task",
					IsDone: true,
					UserId: 1,
				}
				m.On("GetByID", id).Return(existingTask, nil)
			},
		},

		{
			name: "ошибка при обновлении в БД",
			id:   5,
			params: UpdateTaskParams{
				Task:   stringPtr("Updated task"),
				IsDone: boolPtr(true),
			},
			want:    nil,
			wantErr: true,
			mockSetup: func(m *MockTaskRepo, id uint, params UpdateTaskParams, want *Task) {
				existingTask := TaskStruct{
					ID:     id,
					Task:   "Old task",
					IsDone: false,
					UserId: 1,
				}
				m.On("GetByID", id).Return(existingTask, nil)
				m.On("Update", mock.Anything).Return(nil, errors.New("db error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepo)
			tt.mockSetup(mockRepo, tt.id, tt.params, tt.want)

			service := NewTaskService(mockRepo)
			result, err := service.UpdateTask(tt.id, tt.params)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.want.Task, result.Task)
				assert.Equal(t, *tt.want.IsDone, *result.IsDone)
				assert.Equal(t, tt.want.UserId, result.UserId)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteTask(t *testing.T) {
	tasks := []struct {
		name      string
		id        uint
		mockSetup func(m *MockTaskRepo, id uint)
		wantErr   bool
	}{
		{
			name: "успешное удаление задачи",
			id:   1,
			mockSetup: func(m *MockTaskRepo, id uint) {
				existingTask := TaskStruct{
					ID:     id,
					Task:   "Task 1",
					IsDone: false,
					UserId: 1,
				}
				m.On("GetByID", id).Return(existingTask, nil)
				m.On("Delete", &existingTask).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "задача не найдена",
			id:   999,
			mockSetup: func(m *MockTaskRepo, id uint) {
				// GetByID сразу вернет ошибку так как не найдет id
				m.On("GetByID", id).Return(TaskStruct{}, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name: "ошибка при удалении в бд",
			id:   2,
			mockSetup: func(m *MockTaskRepo, id uint) {
				existingTask := TaskStruct{
					ID:     id,
					Task:   "Task 2",
					IsDone: false,
					UserId: 2,
				}
				m.On("GetByID", id).Return(existingTask, nil)
				// ошибка возникает при удалении из бд
				m.On("Delete", &existingTask).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tasks {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepo)
			tt.mockSetup(mockRepo, tt.id)

			service := NewTaskService(mockRepo)
			err := service.DeleteTask(tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
