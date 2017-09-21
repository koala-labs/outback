package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UFOQuery Struct for passing query strings
type UFOQuery struct {
	Cluster string `form:"cluster"`
	Service string `form:"service"`
}

// UFOJson Struct for passing json body
type UFOJson struct {
	Cluster string `json:"cluster" binding:"required"`
	Service string `json:"service" binding:"required"`
	Version string `json:"version" binding:"required"`
}

func routes() *gin.Engine {
	var ufoQuery UFOQuery
	var ufoJSON UFOJson

	ok := func(c *gin.Context) {
		c.String(200, "")
	}

	cors := func(c *gin.Context) {
		c.Writer.Header().Add("access-control-allow-origin", "*")
		c.Writer.Header().Add("access-control-allow-headers", "accept, content-type")
		c.Writer.Header().Add("access-control-allow-methods", "GET,HEAD,POST,DELETE,OPTIONS,PUT,PATCH")
	}

	routes := gin.Default()
	routes.Use(cors)
	routes.OPTIONS("/ufo", ok)
	routes.OPTIONS("/ufo/deploy", ok)

	routes.StaticFile("/", "./app/dist/index.html")
	routes.StaticFS("/static", http.Dir("./app/dist/static"))

	routes.GET("/ufo/clusters", func(c *gin.Context) {
		c.JSON(200, listECSClusters())
	})

	routes.GET("/ufo/services", func(c *gin.Context) {
		if c.BindQuery(&ufoQuery) == nil {
			c.JSON(200, listECSServices(ufoQuery.Cluster))
		}
	})

	routes.GET("/ufo/versions", func(c *gin.Context) {
		if c.BindQuery(&ufoQuery) == nil {
			c.JSON(200, filterImages(listImages(ufoQuery.Service)))
		}
	})

	routes.GET("/ufo/repo", func(c *gin.Context) {
		if c.BindQuery(&ufoQuery) == nil {
			c.JSON(200, getRepoURI(ufoQuery.Service))
		}
	})

	routes.POST("/ufo/deploy", func(c *gin.Context) {
		c.BindJSON(&ufoJSON)
		service, taskDefinitionArn := registerNewDefinition(ufoJSON.Service, ufoJSON.Version)
		c.JSON(201, updateService(ufoJSON.Cluster, service, taskDefinitionArn))
	})

	return routes
}
