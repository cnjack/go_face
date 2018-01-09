package main

import (
	"github.com/cnjack/go_auto_clip/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.Default()

	c := handler.NewImage("resource/haarcascade_frontalface_alt.xml")
	e.POST("/rectangles", c.Rectangles)

	e.LoadHTMLGlob("resource/*.html")
	e.GET("/", c.Html)
	e.Run(":1815")
}
