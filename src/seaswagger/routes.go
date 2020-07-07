package seaswagger

import (
	"flamingo.me/flamingo/v3/framework/web"

	"github.com/pteich/gosea/src/seaswagger/interfaces/controllers"
)

type routes struct {
	apiController *controllers.Api
}

func (r *routes) Inject(controller *controllers.Api) {
	r.apiController = controller
}

func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/swaggerposts", "seaSwagger.Posts")
	registry.HandleGet("seaSwagger.Posts", r.apiController.ShowPostsWithUsers)
}
