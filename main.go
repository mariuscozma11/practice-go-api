// Main package
// Endpoint definition
package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
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
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://mariuscozma.co"}, // Your Vite dev server
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Unprotected routes:
	router.GET("/posts/:id", posts.GetPostByID)
	router.POST("/login", auth.Login)
	router.GET("/posts", posts.GetPosts)
	router.GET("/refresh", auth.RefreshToken)
	router.GET("/", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{
			"message": "Success",
		})
	})

	// Protected routes:
	router.Use(auth.ProtectedRoute())
	{
		router.DELETE("/posts/:id", posts.DeletePost)
		router.POST("/posts", posts.PostPost)
		router.PATCH("/posts/:id", posts.UpdatePostByID)

	}
	router.Run()
}
