package main

import (
	"net/http"
	"strings"

	"fmt"
	"runtime"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.fuzzhq.com/Web-Ops/ufo/pkg/ufo"
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

type AppState struct {
	c       *ecs.Cluster
	s       *ecs.Service
	oldT    *ecs.TaskDefinition
	newT    *ecs.TaskDefinition
	version string
}

func routes(UFO *ufo.UFO) *gin.Engine {
	var ufoQuery UFOQuery
	var ufoJSON UFOJson

	s := &AppState{}

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

	routes.GET("/ufo/status", func(c *gin.Context) {
		if c.BindQuery(&ufoQuery) == nil {
			PollForStatus(c.Writer, c.Request, UFO, s)
		}
	})

	routes.GET("/ufo/clusters", func(c *gin.Context) {
		clusters, err := UFO.Clusters()

		HandleError(err)

		c.JSON(200, clusters)
	})

	routes.GET("/ufo/services", func(c *gin.Context) {
		if c.BindQuery(&ufoQuery) != nil {
			return
		}

		cluster, err := UFO.GetCluster(ufoQuery.Cluster)

		HandleError(err)

		s.c = cluster

		services, err := UFO.Services(s.c)

		HandleError(err)

		c.JSON(200, services)
	})

	routes.GET("/ufo/service", func(c *gin.Context) {
		if c.BindQuery(&ufoQuery) != nil {
			return
		}

		service, err := UFO.GetService(s.c, ufoQuery.Service)

		HandleError(err)

		s.s = service

		c.JSON(200, service)
	})

	routes.GET("/ufo/versions", func(c *gin.Context) {
		if c.BindQuery(&ufoQuery) != nil {
			return
		}

		service, err := UFO.GetService(s.c, ufoQuery.Service)

		HandleError(err)

		s.s = service

		t, err := UFO.GetTaskDefinition(s.c, s.s)

		HandleError(err)

		s.oldT = t

		images, err := UFO.GetImages(t)

		HandleError(err)

		c.JSON(200, images)
	})

	routes.GET("/ufo/commit", func(c *gin.Context) {
		if c.BindQuery(&ufoQuery) != nil {
			return
		}

		commit, err := UFO.GetLastDeployedCommit(*s.s.TaskDefinition)

		HandleError(err)

		c.JSON(200, commit)
	})

	routes.POST("/ufo/deploy", func(c *gin.Context) {
		c.BindJSON(&ufoJSON)

		t, err := UFO.Deploy(s.c, s.s, ufoJSON.Version)

		s.newT = t

		HandleError(err)

		c.JSON(201, t)
	})

	return routes
}

func HandleError(err error) {
	if err == nil {
		return
	}

	parsed, ok := err.(awserr.Error)

	if !ok {
		log.Fatalf("Unable to parse error: %v.\n", err)
	}

	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	log.WithFields(log.Fields{
		"code":  parsed.Code(),
		"error": parsed.Error(),
		"frame": fmt.Sprintf("%s,:%d %s\n", frame.File, frame.Line, frame.Function),
	}).Fatal("Received an error from AWS.")
}
