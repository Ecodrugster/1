package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/itstep-backend/internal/repositories"
	"google.golang.org/api/iterator"
	"cloud.google.com/go/firestore"
)

type NewsItem struct {
	ID          string `json:"id" firestore:"-"`
	Title       string `json:"title" firestore:"title"`
	Description string `json:"description" firestore:"description"`
	CreatedAt   string `json:"created_at" firestore:"created_at"`
}

func GetNews(c *gin.Context) {
	iter := repositories.FirestoreClient.Collection("news").OrderBy("created_at", firestore.Desc).Documents(c.Request.Context())

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
