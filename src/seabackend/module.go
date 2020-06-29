package seabackend

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/web"

	"github.com/pteich/gosea/src/seabackend/application"
	"github.com/pteich/gosea/src/seabackend/domain/service"
	"github.com/pteich/gosea/src/seabackend/infrastructure"
	"github.com/pteich/gosea/src/seabackend/interfaces/controllers"
)

type Module struct {
}

func (m *Module) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(routes))

	injector.Bind(new(controllers.PostsWithUserLoader)).To(new(application.PostsWithUsers))
	injector.Bind(new(service.SeaBackendLoader)).To(new(infrastructure.SeaBackend))
	injector.Bind(new(infrastructure.Cacher)).To(new(infrastructure.RequestCache))
}
