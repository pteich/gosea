package seaswagger

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type Module struct {
}

func (m *Module) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(routes))
}
