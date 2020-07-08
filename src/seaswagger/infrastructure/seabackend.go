package infrastructure

import (
	"context"
	"errors"
	"strconv"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/pteich/gosea/src/seaswagger/infrastructure/client"
	"github.com/pteich/gosea/src/seaswagger/infrastructure/client/operations"
	"github.com/pteich/gosea/src/seaswagger/infrastructure/models"
)

//go:generate swagger generate client -f ../gosea.yaml -A "seabackendapi"

// SeaBackend bundles all function to access external json endpoint
type SeaBackend struct {
	endpoint string
	client   *client.Seabackendapi
}

type config struct {
	Endpoint string `inject:"config:seabackend.host"`
}

func (sb *SeaBackend) Inject(cfg *config) {
	if cfg != nil {
		sb.endpoint = cfg.Endpoint
	}

	transport := httptransport.New(cfg.Endpoint, "", nil)
	client := client.New(transport, strfmt.Default)
	sb.client = client
}

// LoadPosts loads all existing posts from external endpoint
func (sb *SeaBackend) LoadPosts(ctx context.Context) ([]*models.Post, error) {
	posts, err := sb.client.Operations.GetPosts(operations.NewGetPostsParamsWithContext(ctx))
	if err != nil {
		return nil, err
	}

	return posts.GetPayload(), nil
}

func (sb *SeaBackend) LoadUsers(ctx context.Context) ([]*models.User, error) {
	users, err := sb.client.Operations.GetUsers(operations.NewGetUsersParamsWithContext(ctx))
	if err != nil {
		return nil, err
	}

	return users.GetPayload(), nil
}

func (sb *SeaBackend) LoadUser(ctx context.Context, id string) (*models.User, error) {
	params := operations.NewGetUsersParamsWithContext(ctx)
	idArg, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	idArg64 := int64(idArg)
	params.ID = &idArg64

	users, err := sb.client.Operations.GetUsers(params)
	if err != nil {
		return nil, err
	}

	remoteUsers := users.GetPayload()
	if len(remoteUsers) <= 0 {
		return nil, errors.New("user not found")
	}

	return remoteUsers[0], nil
}
