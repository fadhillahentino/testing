package repository

import(
	"github.com/fadhilfcr/oren-service/src/modules/users/model"
)

// Interface User
type UserRepository interface{
	Save(*model.User) error
	Update(string, *model.User) error
	Delete(string) error
	FindById(string) (*model.User,error)
	FindAll() (*model.Users,error)
	CheckLogin(string,string)bool
}