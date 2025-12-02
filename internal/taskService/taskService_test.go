package taskService

import (
	"errors"
	"testing"

	"github.com/AntonRadchenko/WebPet1/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateTask(t *testing.T) {
	// создаем слайс структур, в каждой из которых описан тестовый случай
	tests := []struct {
		name      string                                   // имя теста
		input     *TaskStruct                              // входные данные
		mockSetup func(m *MockTaskRepo, input *TaskStruct) // функция настройки мок-репы
		wantErr   bool                                     // флаг который говорит ожидать ли ошибку
	}{
		{
			name:  "успешное создание задачи",
			input: &TaskStruct{Task: "Test", IsDone: false},
			mockSetup: func(m *MockTaskRepo, input *TaskStruct) {
				// настраиваем мок, чтобы при вызове CreateTask с параметром input
				// он возвращал заранее заданные данные (например, объект или ошибку).
				m.On("Create", input).Return(input, nil)
			},
			wantErr: false,
		},
		{
			name:  "ошибка при создании",
			input: &TaskStruct{Task: "Bad task", IsDone: false},
			mockSetup: func(m *MockTaskRepo, input *TaskStruct) {
				m.On("Create", input).Return(&TaskStruct{}, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests { // проходимся по тестам
		t.Run(tt.name, func(t *testing.T) { // запускаем каждый тест
			mockRepo := new(MockTaskRepo)    // создаем для каждого теста мок-репозиоторий
			tt.mockSetup(mockRepo, tt.input) // настройка мока

			service := NewService(mockRepo)
			result, err := service.CreateTask(openapi.PostTasksJSONRequestBody{ // вызывается метод из сервисного слоя
				Task: tt.input.Task,
				IsDone: &tt.input.IsDone,
			})			

			if tt.wantErr { // если ожидается ошибка, то проверяется что ошибка произошла
				assert.Error(t, err)
			} else { // а если ошибки НЕ ожидается, то проверяется что ее нет, и что результат соответствует ожидаемому входному значению
				assert.NoError(t, err)
				// сравниваем только поля Task, и IsDone так как ID и остальные поля сравнивать не нужно (они все ровно разные так как генерятся по новой)
				assert.Equal(t, tt.input.Task, result.Task)
				assert.Equal(t, tt.input.IsDone, result.IsDone)
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
		wantTasks []TaskStruct
	}{
		{
			name: "успешное получение всех задач",
			mockSetup: func(m *MockTaskRepo) {
				m.On("GetAll").Return([]TaskStruct{
					{Task: "Task 1", IsDone: true},
					{Task: "Task 2", IsDone: false},
				}, nil)
			},
			wantErr:   false,
			wantTasks: []TaskStruct{{Task: "Task 1", IsDone: true}, {Task: "Task 2", IsDone: false}},
		},
		{
			name: "ошибка при получении задач",
			mockSetup: func(m *MockTaskRepo) {
				m.On("GetAll").Return(nil, errors.New("db error"))
			},
			wantErr:   true,
			wantTasks: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepo)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo)
			result, err := service.GetTasks()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Сравниваем длину слайсов, чтобы убедиться что они одинаковые
				assert.Equal(t, len(tt.wantTasks), len(result))

				// Если слайсы не пустые, то проходим по ним и сравниваем только важные поля
				if len(result) > 0 {
					for i := range result {
						assert.Equal(t, tt.wantTasks[i].Task, result[i].Task)
						assert.Equal(t, tt.wantTasks[i].IsDone, result[i].IsDone)
					}
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateTask(t *testing.T) {
    tests := []struct {
        name      string
        id        uint
        input     *TaskStruct
        wantTask  *TaskStruct
        wantErr   bool
        mockSetup func(m *MockTaskRepo, id uint, input, wantTask *TaskStruct)
    }{
        {
            name: "успешное обновление задачи",
            id:   1,
            input: &TaskStruct{
                Task:   "Updated task",
                IsDone: true,
            },
            wantTask: &TaskStruct{
                ID:     1,
                Task:   "Updated task",
                IsDone: true,
            },
            wantErr: false,
            mockSetup: func(m *MockTaskRepo, id uint, input, wantTask *TaskStruct) {
                existingTask := TaskStruct{
                    ID:     id,
                    Task:   "Old task",
                    IsDone: false,
                }
                m.On("GetByID", id).Return(existingTask, nil)
                // используем mock.Anything
                m.On("Update", mock.Anything).Return(wantTask, nil)
            },
        },
        {
            name: "обновление только текста задачи",
            id:   2,
            input: &TaskStruct{
                Task:   "New text",
                IsDone: false, // Не будет передано в API
            },
            wantTask: &TaskStruct{
                ID:     2,
                Task:   "New text",
                IsDone: false, // Старое значение сохранится
            },
            wantErr: false,
            mockSetup: func(m *MockTaskRepo, id uint, input, wantTask *TaskStruct) {
                existingTask := TaskStruct{
                    ID:     id,
                    Task:   "Old task",
                    IsDone: false,
                }
                m.On("GetByID", id).Return(existingTask, nil)
                m.On("Update", mock.Anything).Return(wantTask, nil)
            },
        },
        {
            name: "ошибка - задача не найдена",
            id:   999,
            input: &TaskStruct{
                Task:   "Some task",
                IsDone: true,
            },
            wantTask: nil,
            wantErr:  true,
            mockSetup: func(m *MockTaskRepo, id uint, input, wantTask *TaskStruct) {
                m.On("GetByID", id).Return(TaskStruct{}, errors.New("task not found"))
            },
        },
        {
            name: "ошибка - пустая задача",
            id:   1,
            input: &TaskStruct{
                Task:   "", // Пустая строка
                IsDone: true,
            },
            wantTask: nil,
            wantErr:  true,
            mockSetup: func(m *MockTaskRepo, id uint, input, wantTask *TaskStruct) {
                // Ничего не настраиваем - валидация сработает раньше
            },
        },
        {
            name: "ошибка при обновлении в БД",
            id:   1,
            input: &TaskStruct{
                Task:   "Updated task",
                IsDone: true,
            },
            wantTask: nil,
            wantErr:  true,
            mockSetup: func(m *MockTaskRepo, id uint, input, wantTask *TaskStruct) {
                existingTask := TaskStruct{
                    ID:     id,
                    Task:   "Old task",
                    IsDone: false,
                }
                m.On("GetByID", id).Return(existingTask, nil)
                m.On("Update", mock.Anything).Return(nil, errors.New("db error"))
            },
		},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := new(MockTaskRepo)
            tt.mockSetup(mockRepo, tt.id, tt.input, tt.wantTask)

            service := NewService(mockRepo)
            
            apiRequest := openapi.PatchTasksIdJSONRequestBody{
                Task: &tt.input.Task,
            }
            
			if tt.name == "успешное обновление задачи" || tt.name == "ошибка при обновлении в БД" {
				apiRequest.IsDone = &tt.input.IsDone
			}
            
            result, err := service.UpdateTask(tt.id, apiRequest)

            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, result)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, result)
                assert.Equal(t, tt.wantTask.Task, result.Task)
                assert.Equal(t, tt.wantTask.IsDone, result.IsDone)
            }

            mockRepo.AssertExpectations(t)
        })
    }
}

func TestDeleteTask(t *testing.T) {
	tasks := []struct {
		name string
		id uint
		mockSetup func(m *MockTaskRepo, id uint)
		wantErr bool
	}{
		{
			name: "успешное удаление задачи",
			id: 1,
			mockSetup: func(m *MockTaskRepo, id uint) {
				existingTask := TaskStruct{
					ID: id,
					Task: "Task 1",
					IsDone: false,
				}
				m.On("GetByID", id).Return(existingTask, nil)
				m.On("Delete", &existingTask, id).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "задача не найдена",
			id: 999,
			mockSetup: func(m *MockTaskRepo, id uint) {
				// GetByID сразу вернет ошибку так как не найдет id
				m.On("GetByID", id).Return(TaskStruct{}, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name: "ошибка при удалении в бд",
			id: 2,
			mockSetup: func(m *MockTaskRepo, id uint) {
                existingTask := TaskStruct{
                    ID:     id,
                    Task:   "Task 2",
                    IsDone: false,
                }		
				m.On("GetByID", id).Return(existingTask, nil)
				// ошибка возникает при удалении из бд		
				m.On("Delete", &existingTask, id).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tasks {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepo)
			tt.mockSetup(mockRepo, tt.id)

			service := NewService(mockRepo)
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