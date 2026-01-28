// Package posts
package posts

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/mariuscozma11/practice-go-api/db"
)

type Post struct {
	PostID    string `db:"post_id"`
	Title     string `db:"title"`
	Content   string `db:"content"`
	CreatedAt string `db:"created_at"`
}

func GetPosts(c *gin.Context) {
	ctx := context.Background()
	rows, _ := db.DB.Query(ctx, "select post_id, title, content, created_at::text from posts")
	defer rows.Close()
	posts, err := pgx.CollectRows(rows, pgx.RowToStructByName[Post])
	if err != nil {
		log.Println(err)
		return
	}
	c.IndentedJSON(http.StatusOK, posts)
}

type CreatePost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func PostPost(c *gin.Context) {
	ctx := context.Background()
	var newPost CreatePost
	if err := c.BindJSON(&newPost); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	}
	var post Post
	err := db.DB.QueryRow(ctx, "insert into posts (title,content) values ($1, $2) returning post_id, title, content", newPost.Title, newPost.Content).Scan(&post.PostID, &post.Title, &post.Content)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
	c.IndentedJSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {
	ctx := context.Background()
	var post Post
	id := c.Param("id")

	err := db.DB.QueryRow(ctx, "delete from posts where post_id=$1 returning post_id, title, content", id).Scan(&post.PostID, &post.Title, &post.Content)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
	c.IndentedJSON(http.StatusOK, post)
}

func GetPostByID(c *gin.Context) {
	ctx := context.Background()
	var post Post
	id := c.Param("id")

	err := db.DB.QueryRow(ctx, "select * from posts where post_id=$1", id).Scan(&post.PostID, &post.Title, &post.Content)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
	c.IndentedJSON(http.StatusOK, post)
}

func UpdatePostByID(c *gin.Context) {
	ctx := context.Background()
	var post Post
	id := c.Param("id")
	var updatePost CreatePost
	if err := c.BindJSON(&updatePost); err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
	}

	err := db.DB.QueryRow(ctx, "update posts set title=$1, content=$2 where post_id=$3 returning post_id, title, content", updatePost.Title, updatePost.Content, id).Scan(&post.PostID, &post.Title, &post.Content)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
	c.IndentedJSON(http.StatusOK, post)

}
