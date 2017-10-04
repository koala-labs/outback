package main

import (
	"net/http"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{Asset, AssetDir, AssetInfo, root}
	return &binaryFileSystem{
		fs,
	}
}

// UFOQuery Struct for passing query strings
type UFOQuery struct {
	Cluster        string `form:"cluster"`
	Service        string `form:"service"`
	TaskDefinition string `form:"definition"`
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

	routes.Use(static.Serve("/", BinaryFileSystem("../app/dist/")))

	routes.GET("/ufo/clusters", func(c *gin.Context) {
		c.JSON(200, listECSClusters())
	})

	routes.GET("/ufo/service", func(c *gin.Context) {
		if c.BindQuery(&ufoQuery) == nil {
			c.JSON(200, describeService(ufoQuery.Cluster, ufoQuery.Service))
		}
	})

	routes.GET("/ufo/commit", func(c *gin.Context) {
		if c.BindQuery(&ufoQuery) == nil {
			c.JSON(200, getLastDeployedCommit(ufoQuery.TaskDefinition))
		}
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
		service, taskDefinitionArn := registerNewTaskDefinition(ufoJSON.Service, ufoJSON.Version)
		c.JSON(201, updateService(ufoJSON.Cluster, service, taskDefinitionArn))
	})

	return routes
}
