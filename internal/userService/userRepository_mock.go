package userService

import "github.com/stretchr/testify/mock"

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(user *UserStruct) (*UserStruct, error) {
	args := m.Called(user)
	var u *UserStruct
	if res := args.Get(0); res != nil {
		u = res.(*UserStruct)
	}
	return u, args.Error(1)
}

func (m *MockUserRepo) GetAll() ([]UserStruct, error) {
	args := m.Called() 
	var users []UserStruct
	if res := args.Get(0); res != nil {
		users = res.([]UserStruct)
	}
	return users, args.Error(1)
}

func (m *MockUserRepo) GetByID(id uint) (UserStruct, error) {
    args := m.Called(id) 
    var user UserStruct
    if res := args.Get(0); res != nil {
        user = res.(UserStruct)
    }
    return user, args.Error(1) 	
}

func (m *MockUserRepo) Update(user *UserStruct) (*UserStruct, error) {
    args := m.Called(user) 
    var updatedUser *UserStruct
    if res := args.Get(0); res != nil {
        updatedUser = res.(*UserStruct)
    }
    return updatedUser, args.Error(1) 
}

func (m *MockUserRepo) Delete(user *UserStruct) error {
	args := m.Called(user)
	return args.Error(0)
}
