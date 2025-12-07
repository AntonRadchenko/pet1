package userService

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name      string
		params    CreateUserParams
		want      *User
		mockSetup func(m *MockUserRepo, params CreateUserParams, want *User)
		wantErr   bool
	}{
		{
			name: "успешное создание пользователя",
			params: CreateUserParams{
				Email:    "test@example.com",
				Password: "123",
			},
			want: &User{
				Email: "test@example.com",
			},
			mockSetup: func(m *MockUserRepo, params CreateUserParams, want *User) {
				dbUser := &UserStruct{
					ID:       41,
					Email:    params.Email,
					Password: "$2a$10$hashed123",
				}
				m.On("Create", mock.Anything).Return(dbUser, nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка - email пустой",
			params: CreateUserParams{
				Email:    "",
				Password: "password123",
			},
			want: nil,
			mockSetup: func(m *MockUserRepo, params CreateUserParams, want *User) {
				// Мок не вызывается
			},
			wantErr: true,
		},
		{
			name: "ошибка - пароль пустой",
			params: CreateUserParams{
				Email:    "test@example.com",
				Password: "",
			},
			want: nil,
			mockSetup: func(m *MockUserRepo, params CreateUserParams, want *User) {
				// Мок не вызывается
			},
			wantErr: true,
		},
		{
			name: "ошибка при создании в БД",
			params: CreateUserParams{
				Email:    "test@example.com",
				Password: "password123",
			},
			want:    nil,
			wantErr: true,
			mockSetup: func(m *MockUserRepo, params CreateUserParams, want *User) {
				m.On("Create", mock.Anything).Return(nil, errors.New("db error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.mockSetup(mockRepo, tt.params, tt.want)

			service := NewUserService(mockRepo)
			result, err := service.CreateUser(tt.params)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotZero(t, result.ID) // ID > 0
				assert.Equal(t, tt.want.Email, result.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetUsers(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(m *MockUserRepo)
		wantErr   bool
		want      []User
	}{
		{
			name: "успешное получение всех пользователей",
			mockSetup: func(m *MockUserRepo) {
				m.On("GetAll").Return([]UserStruct{
					{Email: "user1@example.com"},
					{Email: "user2@example.com"},
				}, nil)
			},
			wantErr: false,
			want: []User{
				{Email: "user1@example.com"},
				{Email: "user2@example.com"},
			},
		},
		{
			name: "ошибка при получении пользователей",
			mockSetup: func(m *MockUserRepo) {
				m.On("GetAll").Return(nil, errors.New("db error"))
			},
			wantErr: true,
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.mockSetup(mockRepo)

			service := NewUserService(mockRepo)
			result, err := service.GetUsers()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Сравниваем длину слайсов, чтобы убедиться что они одинаковые
				assert.Equal(t, len(tt.want), len(result))

				// Если слайсы не пустые, то проходим по ним и сравниваем только важные поля
				for i := range result {
					assert.Equal(t, tt.want[i].Email, result[i].Email)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	stringPtr := func(s string) *string { return &s }

	tests := []struct {
		name    string
		id      uint
		params  UpdateUserParams
		want    *User
		wantErr bool
		mockSetup func(m *MockUserRepo, id uint, params UpdateUserParams, want *User)
	}{
		{
			name: "успешное обновление пользователя",
			id:   1,
			params: UpdateUserParams{
				Email:    stringPtr("newemail@example.com"),
				Password: stringPtr("newpassword123"),
			},
			want: &User{
				ID:    1,
				Email: "newemail@example.com",
			},
			wantErr: false,
			mockSetup: func(m *MockUserRepo, id uint, params UpdateUserParams, want *User) {
				// 1. Существующий пользователь в БД (для GetByID)
				existingUser := UserStruct{
					ID:       id,
					Email:    "oldemail@example.com",
					Password: "hashed_old123",
				}
				m.On("GetByID", id).Return(existingUser, nil)

				// 2. Обновлённый пользователь (для Update)
				updatedUser := &UserStruct{
					ID:       id,
					Email:    *params.Email,
					Password: "hashed_new123", 
				}
				m.On("Update", mock.Anything).Return(updatedUser, nil)
			},
		},	
		{
			name: "обновление только email",
			id: 2,
			params: UpdateUserParams{
				Email: stringPtr("newemail@example.com"),
			},
			want: &User{
				ID: 2,
				Email: "newemail@example.com",
			},
			wantErr: false,
			mockSetup: func(m *MockUserRepo, id uint, params UpdateUserParams, want *User) {
				existingUser := UserStruct{
					ID: id,
					Email: "oldemail@example.com",
					Password: "hashed123",
				}
				m.On("GetByID", id).Return(existingUser, nil)

				updatedUser := &UserStruct{
					ID: id,
					Email: *params.Email,
					Password: "hashed123",
				}
				m.On("Update", mock.Anything).Return(updatedUser, nil)
			},
		},
		{
			name: "обновление только пароля",
			id:   3,
			params: UpdateUserParams{
				Password: stringPtr("newpassword123"),
			},
			want: &User{
				ID:    3,
				Email: "existing@example.com", // email остается прежним
			},
			wantErr: false,
			mockSetup: func(m *MockUserRepo, id uint, params UpdateUserParams, want *User) {
				existingUser := UserStruct{
					ID:       id,
					Email:    "existing@example.com",
					Password: "hashed_old123",
				}
				m.On("GetByID", id).Return(existingUser, nil)

				updatedUser := &UserStruct{
					ID:       id,
					Email:    "existing@example.com", // email не меняется
					Password: "hashed_new123", // новый хэш
				}
				m.On("Update", mock.Anything).Return(updatedUser, nil)
			},
		},
		{
			name: "ошибка - пользователь не найден",
			id:   999,
			params: UpdateUserParams{
				Email: stringPtr("newemail@example.com"),
			},
			want:    nil,
			wantErr: true,
			mockSetup: func(m *MockUserRepo, id uint, params UpdateUserParams, want *User) {
				m.On("GetByID", id).Return(UserStruct{}, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "ошибка - пустой email",
			id:   4,
			params: UpdateUserParams{
				Email: stringPtr(""), // пустая строка
			},
			want:    nil,
			wantErr: true,
			mockSetup: func(m *MockUserRepo, id uint, params UpdateUserParams, want *User) {
				existingUser := UserStruct{
					ID:       id,
					Email:    "existing@example.com",
					Password: "hashed123",
				}
				m.On("GetByID", id).Return(existingUser, nil)
			},
		},
		{
			name:   "все поля nil - нет полей для обновления",
			id:     5,
			params: UpdateUserParams{}, // все поля nil
			want:    nil,
			wantErr: true,
			mockSetup: func(m *MockUserRepo, id uint, params UpdateUserParams, want *User) {
				existingUser := UserStruct{
					ID:       id,
					Email:    "existing@example.com",
					Password: "hashed123",
				}
				m.On("GetByID", id).Return(existingUser, nil)
			},
		},
		{
			name: "ошибка - пустой пароль",
			id:   6,
			params: UpdateUserParams{
				Password: stringPtr(""), // пустая строка
			},
			want:    nil,
			wantErr: true,
			mockSetup: func(m *MockUserRepo, id uint, params UpdateUserParams, want *User) {
				existingUser := UserStruct{
					ID:       id,
					Email:    "existing@example.com",
					Password: "hashed123",
				}
				m.On("GetByID", id).Return(existingUser, nil)
			},
		},
		{
			name: "ошибка при обновлении в БД",
			id:   7,
			params: UpdateUserParams{
				Email: stringPtr("newemail@example.com"),
			},
			want:    nil,
			wantErr: true,
			mockSetup: func(m *MockUserRepo, id uint, params UpdateUserParams, want *User) {
				existingUser := UserStruct{
					ID:       id,
					Email:    "existing@example.com",
					Password: "hashed123",
				}
				m.On("GetByID", id).Return(existingUser, nil)
				m.On("Update", mock.Anything).Return(nil, errors.New("db error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			tt.mockSetup(mockRepo, tt.id, tt.params, tt.want)

			service := NewUserService(mockRepo)
			result, err := service.UpdateUser(tt.id, tt.params)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.want.ID, result.ID)
				assert.Equal(t, tt.want.Email, result.Email)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteUser(t *testing.T) {
    tests := []struct {
        name      string
        id        uint
        mockSetup func(m *MockUserRepo, id uint)
        wantErr   bool
    }{
        {
            name: "успешное удаление пользователя",
            id:   1,
            mockSetup: func(m *MockUserRepo, id uint) {
                existingUser := UserStruct{
                    ID:       id,
                    Email:    "user@example.com",
                    Password: "hashed_password",
                }
                m.On("GetByID", id).Return(existingUser, nil)
                m.On("Delete", &existingUser).Return(nil)
            },
            wantErr: false,
        },
        {
            name: "пользователь не найден",
            id:   999,
            mockSetup: func(m *MockUserRepo, id uint) {
                m.On("GetByID", id).Return(UserStruct{}, errors.New("not found"))
            },
            wantErr: true,
        },
        {
            name: "ошибка при удалении в бд",
            id:   2,
            mockSetup: func(m *MockUserRepo, id uint) {
                existingUser := UserStruct{
                    ID:       id,
                    Email:    "user2@example.com",
                    Password: "hashed_password2",
                }
                m.On("GetByID", id).Return(existingUser, nil)
                m.On("Delete", &existingUser).Return(errors.New("db error"))
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := new(MockUserRepo)
            tt.mockSetup(mockRepo, tt.id)

            service := NewUserService(mockRepo)
            err := service.DeleteUser(tt.id)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }

            mockRepo.AssertExpectations(t)
        })
    }
}