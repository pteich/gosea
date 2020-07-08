package controllers

import (
	"context"
	"fmt"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"

	"github.com/pteich/gosea/src/seaswagger/domain/entity"
)

type PostsWithUserLoader interface {
	RetrievePostsWithUsersFromBackend(ctx context.Context, filter string) ([]entity.Post, error)
}

type Api struct {
	responder     *web.Responder
	logger        flamingo.Logger
	postsWithUser PostsWithUserLoader
}

func (a *Api) Inject(postsWithUser PostsWithUserLoader, responder *web.Responder, logger flamingo.Logger) {
	a.postsWithUser = postsWithUser
	a.responder = responder
	a.logger = logger
}

func (a *Api) ShowPostsWithUsers(ctx context.Context, req *web.Request) web.Result {
	filter, _ := req.Query1("filter")

	a.logger.Info(fmt.Sprintf("posts with filter %s", filter))
	responsePosts, err := a.postsWithUser.RetrievePostsWithUsersFromBackend(ctx, filter)
	if err != nil {
		a.logger.Error(err)
		return a.responder.ServerError(err)
	}

	return a.responder.Data(responsePosts)
}
