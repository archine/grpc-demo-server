package service

import (
	"grpc-demo-server/internal/entity"
	"grpc-demo-server/internal/mapper"

	"github.com/archine/gin-plus/v4/component/ioc"
	"github.com/archine/grpc-demo-proto/user"
)

func init() {
	ioc.RegisterBeanDef(&UserService{})
}

type UserService struct {
	userMapper *mapper.UserMapper `autowire:""`
}

// Create 创建用户
func (us *UserService) Create(arg *user.CreateUserRequest) (int, error) {
	return us.userMapper.CreateUser(&entity.User{
		Name:  arg.Name,
		Email: arg.Email,
		Age:   int(arg.Age),
	})
}

// GetById 根据ID获取用户
func (us *UserService) GetById(id int) (*user.User, error) {
	return us.userMapper.GetUserById(id)
}

// List 获取用户列表
func (us *UserService) List() ([]*user.User, error) {
	return us.userMapper.FindUserList()
}
