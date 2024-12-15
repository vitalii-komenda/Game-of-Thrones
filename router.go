package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/vitalii-komenda/got/docs"
)

// @title Game of Thrones API
// @version 1.0
// @description Game of Thrones API

// @host localhost:8080
// @BasePath /
func setupRouter(allControllers AllControllers) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.GET("/characters", allControllers.CharactersController.GetAll)
	r.GET("/characters/:name", allControllers.CharactersController.Get)
	r.POST("/characters", allControllers.CharactersController.Post)
	r.DELETE("/characters/:name", allControllers.CharactersController.Delete)
	r.PUT("/characters/:name", allControllers.CharactersController.Put)

	r.GET("/elastic/search", allControllers.SearchController.GetFromElastic)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
