package service

import "github.com/nestjam/goph-keeper/internal/auth/model"

type FakeAuthService struct {
	RegisterFunc func(user *model.User) (*model.User, error)
}

func (s *FakeAuthService) Register(user *model.User) (*model.User, error) {
	return s.RegisterFunc(user)
}
