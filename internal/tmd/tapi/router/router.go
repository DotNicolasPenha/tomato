package router

import (
	"com.dotvinci.tm/internal/core/distros"
	"com.dotvinci.tm/internal/tmd/tapi/router/declarator"
	"com.dotvinci.tm/internal/tmd/tapi/router/renderer"
)

func Router(cwd string, ctx distros.DistroExecContext) {
	routes := declarator.DeclareRoutes(cwd)
	renderer.RenderRoutes(routes, ctx)
}
