package routes

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type Routes interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Routes) {
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	api := r.Group("/api/v1")
	for _, route := range routes {
		route.Register(api)
	}
	//Throw the error if use wrong method for the routes
	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(404, gin.H{
			"Error": "Not Found",
			"path":  ctx.Request.URL.Path,
		})
	})
}
