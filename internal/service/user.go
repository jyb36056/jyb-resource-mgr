package service

import (
	"context"
	"fmt"
	"sync"

	userv1 "jyb-resource-mgr/api/user/v1"
)

type UserService struct {
	mu     sync.RWMutex           // 读写锁
	store  map[int32]*userv1.User // 内存存储，key为用户id，value为用户信息
	nextID int32                  // 下一个用户id
}

// NewUserService 依赖注入
func NewUserService() *UserService {
	// 初始化UserService
	svc := &UserService{
		store:  make(map[int32]*userv1.User),
		nextID: 1,
	}
	// 初始化默认用户
	svc.store[1] = &userv1.User{Id: 1, Name: "张三", Email: "zhangsan@example.com"}
	return svc
}

// GetUser 获取用户信息
func (s *UserService) GetUser(ctx context.Context, id int32) (*userv1.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock() // 在函数结束时/发生panic时自动解锁
	u, ok := s.store[id]
	if !ok {
		return nil, fmt.Errorf("user %d not found", id)
	}
	return u, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, name, email string) (*userv1.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	u := &userv1.User{Id: s.nextID, Name: name, Email: email}
	s.store[u.Id] = u
	return u, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(ctx context.Context, id int32, name, email string) (*userv1.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.store[id]
	if !ok {
		return nil, fmt.Errorf("user %d not found", id)
	}
	u.Name = name
	u.Email = email
	return u, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.store[id]; !ok {
		return fmt.Errorf("user %d not found", id)
	}
	delete(s.store, id)
	return nil
}

// ListUsers 获取所有用户
func (s *UserService) ListUsers(ctx context.Context) ([]*userv1.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	users := make([]*userv1.User, 0, len(s.store))
	for _, u := range s.store {
		users = append(users, u)
	}
	return users, nil
}
