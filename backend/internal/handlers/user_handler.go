package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/itstep-backend/internal/repositories"
	"google.golang.org/api/iterator"
)

func GetUserProfile(c *gin.Context) {
	firebaseUID := c.GetString("firebase_uid")
	
	doc, err := repositories.FirestoreClient.Collection("users").Doc(firebaseUID).Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found or database error"})
		return
	}

	data := doc.Data()
	data["uid"] = doc.Ref.ID
	c.JSON(http.StatusOK, data)
}

func UpdateUserProfile(c *gin.Context) {
	firebaseUID := c.GetString("firebase_uid")
	
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := repositories.FirestoreClient.Collection("users").Doc(firebaseUID).Set(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, input)
}


func GetUserStats(c *gin.Context) {
	firebaseUID := c.GetString("firebase_uid")

	// Count posts
	postIter := repositories.FirestoreClient.Collection("posts").Where("author_id", "==", firebaseUID).Documents(c.Request.Context())
	postCount := 0
	for {
		_, err := postIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count posts"})
			return
		}
		postCount++
	}

	// Count comments (across all posts)
	commentCount := 0
	postDocs := repositories.FirestoreClient.Collection("posts").Documents(c.Request.Context())
	for {
		postDoc, err := postDocs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate posts for comments"})
			return
		}
		commentsIter := postDoc.Ref.Collection("comments").Where("author_id", "==", firebaseUID).Documents(c.Request.Context())
		for {
			_, err := commentsIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count comments"})
				return
			}
			commentCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{"posts": postCount, "comments": commentCount})
}

func GetAllUsers(c *gin.Context) {
	iter := repositories.FirestoreClient.Collection("users").Limit(50).Documents(c.Request.Context())
	
	var users []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
			return
		}
		
		data := doc.Data()
		data["uid"] = doc.Ref.ID
		users = append(users, data)
	}

	c.JSON(http.StatusOK, users)
}
