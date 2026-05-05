package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/itstep-backend/internal/repositories"
	"google.golang.org/api/iterator"
	"cloud.google.com/go/firestore"
)

type NewsItem struct {
	ID          string    `json:"id" firestore:"-"`
	Title       string    `json:"title" firestore:"title"`
	Description string    `json:"description" firestore:"description"`
	Category    string    `json:"category" firestore:"category"` // news, announcement, event, deadline
	CreatedAt   time.Time `json:"created_at" firestore:"created_at"`
}

func CreateNews(c *gin.Context) {
	var item NewsItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item.CreatedAt = time.Now()

	ref, _, err := repositories.FirestoreClient.Collection("news").Add(c.Request.Context(), item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create news"})
		return
	}
	item.ID = ref.ID
	c.JSON(http.StatusCreated, item)
}

func GetNews(c *gin.Context) {
	category := c.Query("category")
	query := repositories.FirestoreClient.Collection("news").OrderBy("created_at", firestore.Desc)
	
	if category != "" {
		query = query.Where("category", "==", category)
	}

	iter := query.Documents(c.Request.Context())

	var news []NewsItem
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch news"})
			return
		}
		var item NewsItem
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		news = append(news, item)
	}

	c.JSON(http.StatusOK, news)
}
