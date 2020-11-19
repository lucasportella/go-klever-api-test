package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	postpb "github.com/roneycharles/klever/third_party/gen"
	"google.golang.org/grpc"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := postpb.NewPostServiceClient(conn)

	// Set up a http server.
	r := gin.Default()

	r.POST("/posts", func(c *gin.Context) {
		post := postpb.Post{}

		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Missing values"})
			return
		}

		req := &postpb.CreatePostRequest{
			Post: &post,
		}
		res, err := client.CreatePost(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, res.Post)
	})

	r.GET("/posts/:id", func(c *gin.Context) {
		id := c.Param("id")

		req := &postpb.GetPostRequest{Id: id}
		res, err := client.GetPost(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, res.Post)
	})

	r.GET("/posts/", func(c *gin.Context) {

		req := &postpb.ListPostsRequest{}
		_, err := client.ListPosts(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, &postpb.ListPostsRequest{})
	})

	r.DELETE("/posts/:id", func(c *gin.Context) {
		id := c.Param("id")

		req := &postpb.DeletePostRequest{Id: id}
		_, err := client.DeletePost(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"This post was deleted": fmt.Sprint(id),
		})
	})

	r.PUT("/posts/", func(c *gin.Context) {
		post := postpb.Post{}

		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Missing values"})
			return
		}

		req := &postpb.UpdatePostRequest{
			Post: &post,
		}
		res, err := client.UpdatePost(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, res.Post)
	})

	if err := r.Run(":8052"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
