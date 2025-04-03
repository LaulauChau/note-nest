package services

import "context"

type TokenService interface {
	GenerateToken(ctx context.Context) (string, error)
	HashToken(ctx context.Context, token string) (string, error)
}
