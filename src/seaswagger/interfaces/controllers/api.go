package controllers

import (
	"context"
	"errors"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type Api struct {
	responder *web.Responder
	logger    flamingo.Logger
}

func (a *Api) Inject(responder *web.Responder, logger flamingo.Logger) {
	a.responder = responder
	a.logger = logger
}

func (a *Api) ShowPostsWithUsers(ctx context.Context, req *web.Request) web.Result {
	return a.responder.ServerError(errors.New("not implemented"))
}
