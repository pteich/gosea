package controllers

import (
	"context"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"

	"github.com/pteich/gosea/src/seabackend/domain/entity"
)

type PostsWithUserLoader interface {
	RetrievePostsWithUsersFromBackend(ctx context.Context, filter string) ([]entity.Post, error)
}

type Api struct {
	responder     *web.Responder
	postsWithUser PostsWithUserLoader
	logger        flamingo.Logger
}

func (a *Api) Inject(postsWithUser PostsWithUserLoader, responder *web.Responder, logger flamingo.Logger) {
	a.postsWithUser = postsWithUser
	a.responder = responder
	a.logger = logger
}

func (a *Api) ShowPostsWithUsers(ctx context.Context, req *web.Request) web.Result {
	filter, _ := req.Query1("filter")

	a.logger.Info("posts with filter " + filter)
	responsePosts, err := a.postsWithUser.RetrievePostsWithUsersFromBackend(ctx, filter)
	if err != nil {
		a.logger.Error(err)
		return a.responder.ServerError(err)
	}

	return a.responder.Data(responsePosts)
}
