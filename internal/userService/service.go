package userService

import (
	"errors"
	"strings"

	"github.com/AntonRadchenko/WebPet1/openapi"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo UserRepoInterface
}

func NewUserService(r UserRepoInterface) *UserService {
	return &UserService{repo: r}
}

// функция хеширования пароля
func hashPass(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (s *UserService) CreateUser(userRequest openapi.PostUsersJSONRequestBody) (*UserStruct, error) {
	if strings.TrimSpace(string(userRequest.Email)) == "" {
		return nil, errors.New("email is empty")
	}

	if strings.TrimSpace(userRequest.Password) == "" {
		return nil, errors.New("password is empty")
	}

	// хешируем пароль 
	hashedPassword, err := hashPass(userRequest.Password)
	if err != nil {
		return nil, errors.New("fail to hash password")
	}

	user := &UserStruct{
		Email: string(userRequest.Email),
		Password: hashedPassword, // передаю в модель бд захешировнный пароль
	}

	createdUser, err := s.repo.Create(user)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (s *UserService) GetUsers() ([]UserStruct, error) {
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) UpdateUser(id uint, userRequest openapi.PatchUsersIdJSONRequestBody) (*UserStruct, error) {
	user, err := s.repo.GetByID(id)
	if err != nil || user.ID == 0 {
		return nil, errors.New("user not found")
	}

	updated := false

	if userRequest.Email != nil {
		if strings.TrimSpace(string(*userRequest.Email)) == "" {
			return nil, errors.New("email is empty")
		}
		user.Email = string(*userRequest.Email)
		updated = true
	}

	if userRequest.Password != nil {
		if strings.TrimSpace(*userRequest.Password) == "" {
			return nil, errors.New("password is empty")
		}
		// хешируем пароль с реквеста
		hashed, err := hashPass(*userRequest.Password)
		if err != nil {
			return nil, errors.New("fail to hash password")
		}
		user.Password = hashed
		updated = true
	}

	if !updated {
		return nil, errors.New("no fields to update")
	}

	updatedUser, err := s.repo.Update(&user)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (s *UserService) DeleteUser(id uint) error {
	user, err := s.repo.GetByID(id)
	if err != nil || user.ID == 0 {
		return errors.New("user not found")
	}

	err = s.repo.Delete(&user)
	if err != nil {
		return err
	}
	return nil
}