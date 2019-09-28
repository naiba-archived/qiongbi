package main

import "github.com/gin-gonic/gin"

func main() {
	engine := gin.Default()
	engine.Static("/asset", "resource/asset")
	engine.GET("/", func(c *gin.Context) {
		c.File("resource/template/index.html")
	})
	engine.Run(":8080")
}
