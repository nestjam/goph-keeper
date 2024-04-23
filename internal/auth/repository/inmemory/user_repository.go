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

func (r *userRepository) Register(ctx context.Context, user *model.User) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[user.Email]; ok {
		return nil, auth.ErrUserWithEmailIsRegistered
	}

	id := generateID(r.ids)
	createdUser := &model.User{
		ID:       id,
		Email:    user.Email,
		Password: user.Password,
	}
	r.users[createdUser.Email] = createdUser
	return createdUser, nil
}

func generateID(ids map[uuid.UUID]struct{}) uuid.UUID {
	for {
		id := uuid.New()
		if _, ok := ids[id]; !ok {
			ids[id] = struct{}{}
			return id
		}
	}
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	panic("unimplemented")
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if foundUser, ok := r.users[email]; ok {
		return foundUser, nil
	}

	return nil, auth.ErrUserIsNotRegisteredAtEmail
}
