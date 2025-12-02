package taskService

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

func (m *MockTaskRepo) GetAll() ([]TaskStruct, error) {
	args := m.Called() // Проверяем, что метод вызван с правильным параметром
	var tasks []TaskStruct
	if res := args.Get(0); res != nil {
		tasks = res.([]TaskStruct)
	}
	return tasks, args.Error(1)
}

func (m *MockTaskRepo) GetByID(id uint) (TaskStruct, error) {
    args := m.Called(id) // Проверяем, что метод вызван с правильным параметром
    var task TaskStruct
    if res := args.Get(0); res != nil {
        task = res.(TaskStruct)
    }
    return task, args.Error(1) 
}

func (m *MockTaskRepo) Update(task *TaskStruct) (*TaskStruct, error) {
    args := m.Called(task) // Проверяем, что метод вызван с правильным параметром
    var updatedTask *TaskStruct
    if res := args.Get(0); res != nil {
        updatedTask = res.(*TaskStruct)
    }
    return updatedTask, args.Error(1) 
}

func (m *MockTaskRepo) Delete(task *TaskStruct) error {
    args := m.Called(task) // Проверяем, что метод бы9л вызван с правильными параметрами
    return args.Error(0)
}