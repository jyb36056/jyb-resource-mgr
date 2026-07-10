package handler

import (
	"context"

	userv1 "jyb-resource-mgr/api/user/v1"
	"jyb-resource-mgr/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	// 向后兼容，嵌入默认实现，未实现的RPC方法默认返回Unimplemented
	userv1.UnimplementedUserServiceServer
	svc *service.UserService
}

// NewGRPCServer 依赖注入
func NewGRPCServer(svc *service.UserService) *GRPCServer {
	return &GRPCServer{svc: svc}
}

// GetUser 获取用户
func (s *GRPCServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	// 调用业务逻辑svc.GetUser方法
	user, err := s.svc.GetUser(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user %d: %v", req.Id, err)
	}
	return &userv1.GetUserResponse{
		User: &userv1.User{Id: user.Id, Name: user.Name, Email: user.Email},
	}, nil
}

// CreateUser 创建用户
func (s *GRPCServer) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	// 调用业务逻辑svc.CreateUser方法
	user, err := s.svc.CreateUser(ctx, req.Name, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create user: %v", err)
	}
	return &userv1.CreateUserResponse{
		User: &userv1.User{Id: user.Id, Name: user.Name, Email: user.Email},
	}, nil
}

// UpdateUser 更新用户
func (s *GRPCServer) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	// 调用业务逻辑svc.UpdateUser方法
	user, err := s.svc.UpdateUser(ctx, req.Id, req.Name, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "update user %d: %v", req.Id, err)
	}
	return &userv1.UpdateUserResponse{
		User: &userv1.User{Id: user.Id, Name: user.Name, Email: user.Email},
	}, nil
}

// DeleteUser 删除用户
func (s *GRPCServer) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	// 调用业务逻辑svc.DeleteUser方法
	err := s.svc.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "delete user %d: %v", req.Id, err)
	}
	return &userv1.DeleteUserResponse{}, nil
}

// ListUsers 获取用户列表
func (s *GRPCServer) ListUsers(ctx context.Context, req *userv1.ListUsersRequest) (*userv1.ListUsersResponse, error) {
	// 调用业务逻辑svc.ListUsers方法
	users, err := s.svc.ListUsers(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list users: %v", err)
	}
	resp := &userv1.ListUsersResponse{}
	for _, u := range users {
		resp.Users = append(resp.Users, &userv1.User{Id: u.Id, Name: u.Name, Email: u.Email})
	}
	return resp, nil
}
