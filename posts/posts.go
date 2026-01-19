package posts

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type post struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var posts = []post{
	{ID: "1", Title: "My first post", Content: "This is my first post"},
	{ID: "2", Title: "My second post", Content: "This is my second post"},
	{ID: "3", Title: "My third post", Content: "This is my third post"},
}

func GetPosts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, posts)
}

func PostPosts(c *gin.Context) {
	var newPost post
	if err := c.BindJSON(&newPost); err != nil {
		return
	}
	posts = append(posts, newPost)
	c.IndentedJSON(http.StatusCreated, newPost)
}

func GetPostById(c *gin.Context) {
	id := c.Param("id")
	for _, post := range posts {
		if post.ID == id {
			c.IndentedJSON(http.StatusOK, post)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "post not found"})
}
