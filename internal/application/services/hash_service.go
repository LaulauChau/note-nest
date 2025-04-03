package services

import "context"

type HashService interface {
	HashPassword(ctx context.Context, password string) (string, error)
	VerifyPassword(ctx context.Context, hashedPassword, password string) (bool, error)
}
