package service

import "github.com/nestjam/goph-keeper/internal/auth/model"

type AuthServiceMock struct {
	RegisterFunc func(user *model.User) (*model.User, error)
}

func (s *AuthServiceMock) Register(user *model.User) (*model.User, error) {
	return s.RegisterFunc(user)
}
