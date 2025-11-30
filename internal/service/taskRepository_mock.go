package service

import (
	"github.com/stretchr/testify/mock"
)

// MockTaskRepo - вместо настоящего - поддельный репозиторий (для тестирования сервисного слоя)
type MockTaskRepo struct {
	mock.Mock
}

func (m *MockTaskRepo) Create(task *TaskStruct) (*TaskStruct, error) {
	args := m.Called(task) // Called() проверяет что метод вызван с правильными параметрами
	var t *TaskStruct
	if res := args.Get(0); res != nil {
		t = res.(*TaskStruct)
	}
	return t, args.Error(1)
}

func (m *MockTaskRepo) GetAll(tasks *[]TaskStruct) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTaskRepo) GetByID(task *TaskStruct, id uint) error {
	args := m.Called(task, id)
	return args.Error(0)
}

func (m *MockTaskRepo) Update(task *TaskStruct) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepo) Delete(task *TaskStruct, id uint) error {
	args := m.Called(task, id)
	return args.Error(0)
}
