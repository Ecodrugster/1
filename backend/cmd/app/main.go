package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/user/itstep-backend/internal/handlers"
	"github.com/user/itstep-backend/internal/middleware"
	"github.com/user/itstep-backend/internal/repositories"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}

	// Initialize Firestore
	repositories.InitFirestore()

	// Initialize Firebase
	middleware.InitFirebase()

	r := gin.Default()

	// CORS Setup
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // Or specify config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API v1 group
	v1 := r.Group("/api/v1")
	v1.Use(middleware.AuthMiddleware())
	{
		v1.GET("/profile", handlers.GetUserProfile)
		v1.PUT("/profile", handlers.UpdateUserProfile)
		
		// Posts
		v1.POST("/posts", handlers.CreatePost)
		v1.GET("/posts", handlers.GetPosts)
		v1.POST("/posts/:id/like", handlers.LikePost)
		v1.POST("/posts/:id/comments", handlers.AddComment)
		v1.GET("/posts/:id/comments", handlers.GetComments)

		// Users
		v1.GET("/users", handlers.GetAllUsers)

		// Clubs
		v1.GET("/clubs", handlers.GetClubs)
		v1.POST("/clubs", handlers.CreateClub)
		v1.POST("/clubs/:id/join", handlers.JoinClub)
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
