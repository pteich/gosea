package seabackend

import (
	"flamingo.me/flamingo/v3/framework/web"

	"github.com/pteich/gosea/src/seabackend/interfaces/controllers"
)

type routes struct {
	apiController *controllers.Api
}

func (r *routes) Inject(controller *controllers.Api) {
	r.apiController = controller
}

func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/posts", "seaBackend.Posts")
	registry.HandleGet("seaBackend.Posts", r.apiController.ShowPostsWithUsers)
}
