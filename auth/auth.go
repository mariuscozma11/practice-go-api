// Package auth
package auth

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/mariuscozma11/practice-go-api/db"
	"golang.org/x/crypto/bcrypt"
)

type loginDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	// Get context
	ctx := context.Background()

	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("couldnt load env variables bro", err)
	}

	// Bind request body json to loginDetails struct
	var loginDetails loginDetails
	err = c.BindJSON(&loginDetails)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON format",
		})
		return
	}

	// Struct for getting email and password from db
	var user struct {
		email    string `db:"email"`
		password string `db:"password_hash"`
	}

	// Get email and hashedpassword from db(for future I can get just the password and use the email from request body)
	err = db.DB.QueryRow(ctx, "select password_hash,email from users where email= $1", loginDetails.Email).Scan(&user.password, &user.email)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Compare hashed db password with request body password
	err = bcrypt.CompareHashAndPassword([]byte(user.password), []byte(loginDetails.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	// Create access token 30 minute ttl
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.email,
		"exp": time.Now().Add(time.Minute).Unix(),
	})

	// Sign access token  with secret and convert to string
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Println(err)
		return
	}

	// Create the refresh token 1 month ttl
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.email,
		"exp": time.Now().Add(time.Hour * 730).Unix(),
	})

	// Sign refresh token with secret and convert to string
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Println(err)
		return
	}

	// Update refresh token
	_, err = db.DB.Query(ctx, "update users set refresh_token=$1 where email=$2", refreshTokenString, user.email)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update refresh token",
		})
	}

	// Attach access token and refresh token to cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Refresh_Token", refreshTokenString, 3600*24*30, "", "", false, true)

	c.IndentedJSON(http.StatusOK, gin.H{
		"acess_token": accessTokenString,
	})
}

func ProtectedRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("couldnt load env variables bro", err)
		}

		auth := c.GetHeader("Authorization")
		prefix := "Bearer: "
		if !strings.HasPrefix(auth, prefix) {
			c.AbortWithStatusJSON(http.StatusBadRequest, "Missing auth header")
			return
		}
		tokenString := strings.TrimPrefix(auth, prefix)
		claims := jwt.StandardClaims{}

		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			if claims.ExpiresAt < time.Now().Unix() {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "token_expired",
				})
				return
			}

			if err == jwt.ErrSignatureInvalid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "unauthorized",
				})
				return
			}
			log.Printf("Failed to parse payload %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "bad request",
			})
			return
		}
		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}
		c.IndentedJSON(http.StatusOK, claims.ExpiresAt)
		c.Next()
	}
}
