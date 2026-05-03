package handlers

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/user/itstep-backend/internal/repositories"
	"google.golang.org/api/iterator"
)

// AdminGetUsers - Получение списка пользователей с пагинацией
func AdminGetUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit := 20
	fmt.Sscanf(limitStr, "%d", &limit)

	iter := repositories.FirestoreClient.Collection("users").Limit(limit).Documents(c.Request.Context())
	
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

	c.JSON(http.StatusOK, gin.H{"users": users, "total": len(users)})
}

// AdminUpdateUserRole - Смена роли пользователя
func AdminUpdateUserRole(c *gin.Context) {
	uid := c.Param("id")
	var input struct {
		Role string `json:"role"`
	}
	
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.Role != "admin" && input.Role != "student" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role value"})
		return
	}

	_, err := repositories.FirestoreClient.Collection("users").Doc(uid).Update(c.Request.Context(), []firestore.Update{
		{Path: "role", Value: input.Role},
	})

	if err != nil {
		// Если документ не существует, создаем его с ролью
		_, err = repositories.FirestoreClient.Collection("users").Doc(uid).Set(c.Request.Context(), map[string]interface{}{
			"role": input.Role,
		}, repositories.MergeAll)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}

// AdminGetPosts - Получение всех постов с пагинацией
func AdminGetPosts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit := 50
	fmt.Sscanf(limitStr, "%d", &limit)

	iter := repositories.FirestoreClient.Collection("posts").OrderBy("created_at", repositories.Descending).Limit(limit).Documents(c.Request.Context())
	
	var posts []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
			return
		}
		
		data := doc.Data()
		data["id"] = doc.Ref.ID
		posts = append(posts, data)
	}

	c.JSON(http.StatusOK, posts)
}

// AdminDeletePost - Удаление поста модератором
func AdminDeletePost(c *gin.Context) {
	id := c.Param("id")
	
	_, err := repositories.FirestoreClient.Collection("posts").Doc(id).Delete(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted by moderator"})
}

// AdminDeleteNews - Удаление новости
func AdminDeleteNews(c *gin.Context) {
	id := c.Param("id")
	_, err := repositories.FirestoreClient.Collection("news").Doc(id).Delete(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete news"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "News deleted"})
}

// AdminDeleteClub - Удаление клуба
func AdminDeleteClub(c *gin.Context) {
	id := c.Param("id")
	_, err := repositories.FirestoreClient.Collection("clubs").Doc(id).Delete(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete club"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Club deleted"})
}

// AdminUpdateNews - Редактирование новости
func AdminUpdateNews(c *gin.Context) {
	id := c.Param("id")
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	_, err := repositories.FirestoreClient.Collection("news").Doc(id).Set(c.Request.Context(), input, repositories.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update news"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "News updated"})
}

// AdminUpdateClub - Редактирование клуба
func AdminUpdateClub(c *gin.Context) {
	id := c.Param("id")
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	_, err := repositories.FirestoreClient.Collection("clubs").Doc(id).Set(c.Request.Context(), input, repositories.MergeAll)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update club"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Club updated"})
}
