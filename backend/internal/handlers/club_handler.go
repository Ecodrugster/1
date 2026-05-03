package handlers

import (
	"log"
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
		log.Printf("[Clubs Debug] Club: %s, Members: %v", data["name"], data["members"])
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

	// Добавляем создателя в список участников автоматически
	userID := c.GetString("firebase_uid")
	club["members"] = []string{userID}
	club["created_by"] = userID

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

func UpdateClub(c *gin.Context) {
	clubID := c.Param("id")
	userID := c.GetString("firebase_uid")
	
	ref := repositories.FirestoreClient.Collection("clubs").Doc(clubID)
	doc, err := ref.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Club not found"})
		return
	}

	if doc.Data()["created_by"] != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the creator can edit this club"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err = ref.Set(c.Request.Context(), input, repositories.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update club"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Club updated successfully"})
}

func DeleteClub(c *gin.Context) {
	clubID := c.Param("id")
	userID := c.GetString("firebase_uid")
	
	ref := repositories.FirestoreClient.Collection("clubs").Doc(clubID)
	doc, err := ref.Get(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Club not found"})
		return
	}

	if doc.Data()["created_by"] != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only the creator can delete this club"})
		return
	}

	_, err = ref.Delete(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete club"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Club deleted successfully"})
}
func LeaveClub(c *gin.Context) {
	clubID := c.Param("id")
	userID := c.GetString("firebase_uid")
	
	ref := repositories.FirestoreClient.Collection("clubs").Doc(clubID)
	_, err := ref.Update(c.Request.Context(), []firestore.Update{
		{
			Path:  "members",
			Value: firestore.ArrayRemove(userID),
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave club"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Left successfully"})
}
