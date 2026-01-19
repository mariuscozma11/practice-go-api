// Main package
// Endpoint definition
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mariuscozma11/practice-go-api/db"
	"github.com/mariuscozma11/practice-go-api/posts"
)

func main() {
	db.Greet()
	router := gin.Default()
	router.GET("/posts", posts.GetPosts)
	router.GET("posts/:id", posts.GetPostById)
	router.POST("/posts", posts.PostPosts)
	router.Run("localhost:8080")
}
