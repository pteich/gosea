package service

//go:generate mockery --all

import (
	"context"

	"github.com/pteich/gosea/src/seabackend/domain/entity"
)

type SeaBackendLoader interface {
	LoadPosts(ctx context.Context) ([]entity.RemotePost, error)
	LoadUsers(ctx context.Context) ([]entity.RemoteUser, error)
	LoadUser(ctx context.Context, id string) (entity.RemoteUser, error)
}
