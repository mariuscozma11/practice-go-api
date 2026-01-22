// Package auth
package auth

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/mariuscozma11/practice-go-api/db"
	"golang.org/x/crypto/bcrypt"
)

type loginDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type dbPass struct {
	password string `db:"password_hash"`
}

func LoginUser(c *gin.Context) {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("couldnt load env variables bro", err)
	}
	var loginDetails loginDetails
	err = c.BindJSON(&loginDetails)
	if err != nil {
		c.Error(err)
	}
	var user struct {
		email    string
		password string
	}
	err = db.DB.QueryRow(ctx, "select password_hash,email from users where email=$1", loginDetails.Email).Scan(&user.password, &user.email)
	if err != nil {
		c.Error(err)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.password), []byte(loginDetails.Password))
	if err != nil {
		c.Error(err)
		return
	}
	// Create the access token 30 minute lived
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    user.email,
		"password": user.password,
		"exp":      time.Now().Add(time.Minute * 30).Unix(),
	})
	accessTokenString, err := accessToken.SignedString([]byte("123"))
	// Create the refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":    user.email,
		"password": user.password,
		"exp":      time.Now().Add(time.Hour * 730).Unix(),
	})
	refreshTokenString, err := refreshToken.SignedString([]byte("123"))
}
