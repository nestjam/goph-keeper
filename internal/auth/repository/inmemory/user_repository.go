package inmemory

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/nestjam/goph-keeper/internal/auth"
	"github.com/nestjam/goph-keeper/internal/auth/model"
)

type userRepository struct {
	users map[string]*model.User
	ids   map[uuid.UUID]struct{}
	mu    sync.Mutex
}

func NewUserRepository() auth.UserRepository {
	return &userRepository{
		users: make(map[string]*model.User),
		ids:   make(map[uuid.UUID]struct{}),
	}
}

func (r *userRepository) Register(ctx context.Context, user *model.User) (uuid.UUID, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user.Password == "" {
		return uuid.Nil, auth.ErrUserPasswordIsEmpty
	}

	if _, ok := r.users[user.Email]; ok {
		return uuid.Nil, auth.ErrUserWithEmailIsRegistered
	}

	id := uuid.New()
	createdUser := &model.User{
		ID:       id,
		Email:    user.Email,
		Password: user.Password,
	}
	r.users[createdUser.Email] = createdUser

	return id, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user, ok := r.users[email]; ok {
		return user, nil
	}

	return nil, auth.ErrUserIsNotRegistered
}
