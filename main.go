// Main package
// Endpoint definition
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mariuscozma11/practice-go-api/auth"
	"github.com/mariuscozma11/practice-go-api/db"
	"github.com/mariuscozma11/practice-go-api/posts"
)

func main() {
	var err error
	router := gin.Default()
	db.DB, err = db.NewDBPool()
	defer db.DB.Close()
	if err != nil {
		log.Fatal("Error creating DB Pool:", err)
	}
	// Unprotected routes:
	router.GET("/posts/:id", posts.GetPostByID)
	router.POST("/login", auth.Login)
	router.GET("/posts", posts.GetPosts)
	router.GET("/refresh", auth.RefreshToken)

	// Protected routes:
	router.Use(auth.ProtectedRoute())
	{
		router.DELETE("/posts/:id", posts.DeletePost)
		router.POST("/posts", posts.PostPost)
		router.PATCH("/posts/:id", posts.UpdatePostByID)

	}
	router.Run()
}
