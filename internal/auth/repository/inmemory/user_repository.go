package inmemory

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/nestjam/goph-keeper/internal/auth"
	"github.com/nestjam/goph-keeper/internal/auth/model"
)

type userRepository struct {
	users map[string]model.User
	ids   map[uuid.UUID]struct{}
	mu    sync.Mutex
}

func NewUserRepository() auth.UserRepository {
	return &userRepository{
		users: make(map[string]model.User),
		ids:   make(map[uuid.UUID]struct{}),
	}
}

func (r *userRepository) Register(ctx context.Context, user model.User) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[user.Email]; ok {
		return model.User{}, auth.ErrUserWithEmailIsRegistered
	}

	id := uuid.New()
	createdUser := model.User{
		ID:       id,
		Email:    user.Email,
		Password: user.Password,
	}
	r.users[createdUser.Email] = createdUser

	return createdUser, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user, ok := r.users[email]; ok {
		return user, nil
	}

	return model.User{}, auth.ErrUserIsNotRegistered
}
