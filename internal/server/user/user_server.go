package user

import (
	"context"
	"grpc-demo-server/internal/service"

	"github.com/archine/gin-plus/v4/component/ioc"
	"github.com/archine/grpc-demo-proto/user"
	"google.golang.org/grpc"
)

func init() {
	ioc.RegisterBeanDef(&Server{})
}

type Server struct {
	user.UnimplementedUserServiceServer
	server      *grpc.Server         `autowire:""`
	userService *service.UserService `autowire:""`
}

func (s *Server) BeanPostConstruct() {
	user.RegisterUserServiceServer(s.server, s)
}

// CreateUser 创建用户
func (s *Server) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	userId, err := s.userService.Create(req)
	if err != nil {
		return nil, err
	}

	return &user.CreateUserResponse{Id: int32(userId)}, nil
}

// GetUser 根据id获取用户
func (s *Server) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	u, err := s.userService.GetById(int(req.Id))
	if err != nil {
		return nil, err
	}

	return &user.GetUserResponse{User: u}, nil

}

// FindUserList 获取用户列表
func (s *Server) FindUserList(ctx context.Context, req *user.FindUserListRequest) (*user.FindUserListResponse, error) {
	users, err := s.userService.List()
	if err != nil {
		return nil, err
	}
	
	return &user.FindUserListResponse{Users: users}, nil
}
