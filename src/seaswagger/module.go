package seaswagger

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/web"

	"github.com/pteich/gosea/src/seaswagger/application"
	"github.com/pteich/gosea/src/seaswagger/domain/service"
	"github.com/pteich/gosea/src/seaswagger/infrastructure"
	"github.com/pteich/gosea/src/seaswagger/interfaces/controllers"
)

type Module struct {
}

func (m *Module) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(routes))

	injector.Bind(new(controllers.PostsWithUserLoader)).To(new(application.PostsWithUsers))
	injector.Bind(new(service.SeaBackendLoader)).To(new(infrastructure.SeaBackend))
}
