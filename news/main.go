package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	v1 := r.Group("/v1")
	v1.GET("/sources", GetSources)
	v1.GET("/news/:source", GetNews)

	r.Run(":7891")
}
