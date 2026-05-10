package handlers

import (
	"net/http"
	"strings"
	"time"

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

	c.JSON(http.StatusOK, normalizeUserDocument(doc.Data(), doc.Ref.ID))
}

func UpdateUserProfile(c *gin.Context) {
	firebaseUID := c.GetString("firebase_uid")

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	userRef := repositories.FirestoreClient.Collection("users").Doc(firebaseUID)

	existingData := map[string]interface{}{}
	existingDoc, err := userRef.Get(c.Request.Context())
	if err == nil {
		existingData = existingDoc.Data()
	}

	// Users can update only own public profile fields.
	// Role/group and other privileged fields are controlled by admin only.
	update := map[string]interface{}{
		"updated_at": now,
		"role":       normalizeRole(asString(existingData["role"])),
	}

	if update["role"] == "student" {
		// role already defaults to student, keep explicit for new docs as well.
		update["role"] = "student"
	}

	if createdAt, ok := existingData["created_at"]; ok {
		update["created_at"] = createdAt
	} else {
		update["created_at"] = now
	}

	if groupName := asString(existingData["group_name"]); groupName != "" {
		update["group_name"] = groupName
		update["group"] = groupName
	} else if groupName := asString(existingData["group"]); groupName != "" {
		update["group_name"] = groupName
		update["group"] = groupName
	}

	displayName := strings.TrimSpace(asString(input["display_name"]))
	if displayName == "" {
		displayName = strings.TrimSpace(asString(input["displayName"]))
	}
	if displayName != "" {
		update["display_name"] = displayName
		update["displayName"] = displayName
	}

	photoURL := strings.TrimSpace(asString(input["photo_url"]))
	if photoURL == "" {
		photoURL = strings.TrimSpace(asString(input["photoURL"]))
	}
	if photoURL != "" {
		update["photo_url"] = photoURL
		update["photoURL"] = photoURL
	}

	email := strings.TrimSpace(asString(input["email"]))
	if email != "" {
		update["email"] = email
	} else if existingEmail := asString(existingData["email"]); existingEmail != "" {
		update["email"] = existingEmail
	}

	_, err = userRef.Set(c.Request.Context(), update, repositories.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	savedDoc, err := userRef.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, normalizeUserDocument(update, firebaseUID))
		return
	}

	c.JSON(http.StatusOK, normalizeUserDocument(savedDoc.Data(), firebaseUID))
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

	// Count comments
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

		users = append(users, normalizeUserDocument(doc.Data(), doc.Ref.ID))
	}

	c.JSON(http.StatusOK, users)
}
