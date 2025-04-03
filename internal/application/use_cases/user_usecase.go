package use_cases

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/LaulauChau/note-nest/internal/application/services"
	"github.com/LaulauChau/note-nest/internal/domain/entities"
	"github.com/LaulauChau/note-nest/internal/domain/repositories"
)

type UserUseCase struct {
	userRepo    repositories.UserRepository
	hashService services.HashService
}

func NewUserUseCase(
	userRepo repositories.UserRepository,
	hashService services.HashService,
) *UserUseCase {
	return &UserUseCase{
		userRepo:    userRepo,
		hashService: hashService,
	}
}

func (uc *UserUseCase) RegisterUser(ctx context.Context, email, name, password string) (*entities.User, error) {
	// Check if username is already taken
	existingUser, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		log.Printf("error getting user by email: %v", err)
		return nil, err
	}
	if existingUser != nil {
		log.Printf("email already taken")
		return nil, errors.New("email already taken")
	}

	// Hash the password
	hashedPassword, err := uc.hashService.HashPassword(ctx, password)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		return nil, err
	}

	// Create a new user
	now := time.Now()
	user := &entities.User{
		ID:        uuid.New().String(),
		Email:     email,
		Name:      name,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save the user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		log.Printf("error creating user: %v", err)
		return nil, err
	}

	// Don't return the password hash
	user.Password = ""
	return user, nil
}

func (uc *UserUseCase) AuthenticateUser(ctx context.Context, email, password string) (*entities.User, error) {
	// Get the user by email
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify the password
	valid, err := uc.hashService.VerifyPassword(ctx, user.Password, password)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("invalid credentials")
	}

	// Don't return the password hash
	user.Password = ""
	return user, nil
}
