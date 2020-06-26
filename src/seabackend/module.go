package seabackend

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type Module struct {
}

func (m *Module) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(routes))
}

type routes struct {
	seaBackendController *SeaBackend
}

func (r *routes) Inject(controller *SeaBackend) *routes {
	r.seaBackendController = controller
	return r
}

func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.HandleGet("seaBackend.Posts", r.seaBackendController.Posts)
	registry.MustRoute("/api", "seaBackend.Posts")
	registry.MustRoute("/posts", "seaBackend.Posts")
}
