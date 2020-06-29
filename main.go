package main

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/core/healthcheck"
	"flamingo.me/flamingo/v3/core/requestlogger"

	"github.com/pteich/gosea/src/seabackend"
)

var Version = "latest"

func main() {
	flamingo.App([]dingo.Module{
		new(requestlogger.Module),
		new(healthcheck.Module),
		new(seabackend.Module),
	})
}
