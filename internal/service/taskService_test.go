package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	// создаем слайс структур, в каждой из которых описан тестовый случай
	tests := []struct {
		name      string                                   // имя теста
		input     *TaskStruct                              // входные данные
		mockSetup func(m *MockTaskRepo, input *TaskStruct) // функция настройки мок репозитория
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
			result, err := service.CreateTask(tt.input.Task, &tt.input.IsDone) // вызывается метод из сервисного слоя

			if tt.wantErr { // если ожидается ошибка, то проверяется что ошибка произошла
				assert.Error(t, err)
			} else { // а если ошибки НЕ ожидается, то проверяется что ее нет, и что результат соответствует ожидаемому входному значению
				assert.NoError(t, err)
				assert.Equal(t, tt.input, result)
			}

			mockRepo.AssertExpectations(t) // проверяем что все ожидаемые вызовы методов мока были выполнены
		})
	}
}

func TestGetTasks(* *testing.T) {
	tests := []struct {
		name string
		mockSetup func(m *MockTaskRepo)
		wantErr bool
		wantTasks []TaskStruc
	}{
		{
			// понять в каком виде должны быть поля в примерах тестов
		}
	}
}