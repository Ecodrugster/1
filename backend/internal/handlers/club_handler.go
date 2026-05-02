package handlers

import (
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/user/itstep-backend/internal/repositories"
	"google.golang.org/api/iterator"
)

func GetClubs(c *gin.Context) {
	iter := repositories.FirestoreClient.Collection("clubs").Documents(c.Request.Context())
	
	var clubs []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch clubs"})
			return
		}
		
		data := doc.Data()
		data["id"] = doc.Ref.ID
		clubs = append(clubs, data)
	}

	c.JSON(http.StatusOK, clubs)
}

func CreateClub(c *gin.Context) {
	var club map[string]interface{}
	if err := c.ShouldBindJSON(&club); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, _, err := repositories.FirestoreClient.Collection("clubs").Add(c.Request.Context(), club)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create club"})
		return
	}

	c.JSON(http.StatusCreated, club)
}

func JoinClub(c *gin.Context) {
	clubID := c.Param("id")
	userID := c.GetString("firebase_uid")
	
	ref := repositories.FirestoreClient.Collection("clubs").Doc(clubID)
	_, err := ref.Update(c.Request.Context(), []firestore.Update{
		{
			Path:  "members",
			Value: firestore.ArrayUnion(userID),
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join club"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined successfully"})
}
