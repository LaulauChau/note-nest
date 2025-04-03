package repositories

import (
	"context"

	"github.com/LaulauChau/note-nest/internal/domain/entities"
	"github.com/LaulauChau/note-nest/internal/domain/repositories"
	"github.com/google/uuid"
)

type UserRepositoryImpl struct {
	q *Queries
}

func NewUserRepository(q *Queries) repositories.UserRepository {
	return &UserRepositoryImpl{q: q}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	// Generate UUID if not provided
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Parse the ID into UUID
	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return err
	}

	// Use the manual query instead of CreateUser since it doesn't include ID
	_, err = r.q.db.Exec(ctx,
		"INSERT INTO users (id, email, name, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, user.Email, user.Name, user.Password, user.CreatedAt, user.UpdatedAt)

	return err
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	user, err := r.q.GetUserByID(ctx, userID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return &entities.User{
		ID:       user.ID.String(),
		Email:    user.Email,
		Name:     user.Name,
		Password: user.Password,
	}, nil
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	user, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	return &entities.User{
		ID:       user.ID.String(),
		Email:    user.Email,
		Name:     user.Name,
		Password: user.Password,
	}, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	userID, err := uuid.Parse(user.ID)
	if err != nil {
		return err
	}

	params := UpdateUserParams{
		ID:       userID,
		Email:    user.Email,
		Name:     user.Name,
		Password: user.Password,
	}

	return r.q.UpdateUser(ctx, params)
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.q.DeleteUser(ctx, userID)
}
