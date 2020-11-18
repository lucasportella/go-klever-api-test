package routes

import (
	"github.com/gin-gonic/gin"
)

type Routes struct {
}

func (c Routes) StartGin() *gin.Engine {
	r := gin.New()

	r.POST("/posts", handlers.CreatePost)
	// r.GET("/posts", handlers.ListPosts)
	// r.GET("/posts/:id", handlers.GetPost)
	// r.PUT("/posts/:id", handlers.UpdatePost)
	// r.DELETE("/posts/:id", handlers.DeletePost)

	return r
}
