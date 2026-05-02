package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"cloud.google.com/go/firestore"
	"github.com/user/itstep-backend/internal/models"
	"github.com/user/itstep-backend/internal/repositories"
	"google.golang.org/api/iterator"
)

func CreatePost(c *gin.Context) {
	authorID := c.GetString("firebase_uid")
	
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.AuthorID = authorID
	post.CreatedAt = time.Now()
	post.Likes = 0

	// Save to Firestore
	_, _, err := repositories.FirestoreClient.Collection("posts").Add(c.Request.Context(), post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func GetPosts(c *gin.Context) {
	iter := repositories.FirestoreClient.Collection("posts").OrderBy("created_at", firestore.Desc).Limit(20).Documents(c.Request.Context())
	
	var posts []models.Post
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
			return
		}
		
		var post models.Post
		doc.DataTo(&post)
		post.ID = doc.Ref.ID
		posts = append(posts, post)
	}

	c.JSON(http.StatusOK, posts)
}

func LikePost(c *gin.Context) {
	postID := c.Param("id")
	
	// Atomic increment of likes in Firestore
	ref := repositories.FirestoreClient.Collection("posts").Doc(postID)
	err := repositories.FirestoreClient.RunTransaction(c.Request.Context(), func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(ref)
		if err != nil {
			return err
		}
		
		likes, _ := doc.Data()["likes"].(int64)
		return tx.Update(ref, []firestore.Update{{Path: "likes", Value: likes + 1}})
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post liked"})
}

func AddComment(c *gin.Context) {
	postID := c.Param("id")
	authorID := c.GetString("firebase_uid")

	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment.PostID = postID
	comment.AuthorID = authorID
	comment.CreatedAt = time.Now()

	_, _, err := repositories.FirestoreClient.Collection("posts").Doc(postID).Collection("comments").Add(c.Request.Context(), comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func GetComments(c *gin.Context) {
	postID := c.Param("id")
	iter := repositories.FirestoreClient.Collection("posts").Doc(postID).Collection("comments").OrderBy("created_at", firestore.Asc).Documents(c.Request.Context())

	var comments []models.Comment
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
			return
		}

		var comment models.Comment
		doc.DataTo(&comment)
		comment.ID = doc.Ref.ID
		comments = append(comments, comment)
	}

	c.JSON(http.StatusOK, comments)
}
