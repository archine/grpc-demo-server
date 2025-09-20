package mapper

import (
	"errors"
	"grpc-demo-server/internal/entity"

	"github.com/archine/gin-plus/v4/component/ioc"
	"github.com/archine/grpc-demo-proto/user"
	"gorm.io/gorm"
)

func init() {
	ioc.RegisterBeanDef(&UserMapper{})
}

// UserMapper 用户数据操作层
type UserMapper struct {
	db *gorm.DB `autowire:""`
}

func (um *UserMapper) CreateUser(user *entity.User) (int, error) {
	err := um.db.Create(user).Error
	if err != nil {
		return 0, err
	}

	return user.Id, nil
}

func (um *UserMapper) GetUserById(id int) (*user.User, error) {
	var u user.User
	err := um.db.Take(&u, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 用户未找到
		}
		return nil, err
	}

	return &u, nil
}

func (um *UserMapper) FindUserList() ([]*user.User, error) {
	var users []*user.User
	err := um.db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
