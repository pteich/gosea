package service

//go:generate mockery --all

import (
	"context"

	"github.com/pteich/gosea/src/seaswagger/infrastructure/models"
)

type SeaBackendLoader interface {
	LoadPosts(ctx context.Context) ([]*models.Post, error)
	LoadUsers(ctx context.Context) ([]*models.User, error)
	LoadUser(ctx context.Context, id string) (*models.User, error)
}
