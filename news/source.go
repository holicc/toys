package main

import "github.com/gin-gonic/gin"

type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetSources(c *gin.Context) {
	c.JSON(200, []Source{
		{ID: "cctv", Name: "CCTV-News"},
		{ID: "cctv1", Name: "CCTV-News1"},
		{ID: "cctv2", Name: "CCTV-News2"},
		{ID: "cctv3", Name: "CCTV-News3"},
	})
}
