package router

import (
	"com.dotvinci.tm/internal/common/logger"
	"com.dotvinci.tm/internal/core/distros"
	"com.dotvinci.tm/internal/tmd/tapi/router/declarator"
	"com.dotvinci.tm/internal/tmd/tapi/router/renderer"
)

func Router(cwd string, ctx distros.DistroExecContext) {
	routes, err := declarator.DeclareRoutes(cwd)
	if err != nil {
		logger.Error("erro to router routes in tapi")
	}
	renderer.RenderRoutes(routes, ctx)
}
