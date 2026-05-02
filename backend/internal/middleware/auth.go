package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

var authClient *auth.Client

func InitFirebase() {
	ctx := context.Background()
	saPath := os.Getenv("FIREBASE_SERVICE_ACCOUNT_PATH")
	
	var app *firebase.App
	var err error

	if saPath != "" {
		opt := option.WithCredentialsFile(saPath)
		app, err = firebase.NewApp(ctx, nil, opt)
	} else {
		app, err = firebase.NewApp(ctx, nil)
	}

	if err != nil {
		log.Printf("error initializing firebase app: %v\n", err)
		return
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Printf("error getting firebase auth client: %v\n", err)
		return
	}

	authClient = client
	log.Println("Firebase Admin SDK initialized")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if authClient == nil {
			c.Next() // Allow in dev if not configured
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		idToken := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
		token, err := authClient.VerifyIDToken(c.Request.Context(), idToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store firebase UID in context
		c.Set("firebase_uid", token.UID)
		c.Next()
	}
}
