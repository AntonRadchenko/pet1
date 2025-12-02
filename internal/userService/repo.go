package userService

type UserRepoInterface interface {
	Create(user *UserStruct) (*UserStruct, error)
	GetAll() ([]UserStruct, error)
	GetByID(id uint) (UserStruct, error)
	GetByEmail(email string) (UserStruct, error) // пока неполнятно зачем данный метод
	Update(user *UserStruct) (*UserStruct, error)
	Delete(user *UserStruct) error
}

type UserRepo struct{}

func Create(user *UserStruct) (*UserStruct, error) {

}