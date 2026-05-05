package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/itstep-backend/internal/repositories"
	"google.golang.org/api/iterator"
)

type GradeInput struct {
	StudentID string `json:"student_id" binding:"required"`
	Subject   string `json:"subject" binding:"required"`
	Value     int    `json:"value" binding:"required"`
	Comment   string `json:"comment"`
}

// AddGrade - Выставление оценки учителем
func AddGrade(c *gin.Context) {
	teacherID := c.GetString("firebase_uid")
	var input GradeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grade := map[string]interface{}{
		"student_id": input.StudentID,
		"teacher_id": teacherID,
		"subject":    input.Subject,
		"value":      input.Value,
		"comment":    input.Comment,
		"created_at": time.Now(),
	}

	_, _, err := repositories.FirestoreClient.Collection("grades").Add(c.Request.Context(), grade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add grade"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Grade added successfully"})
}

// GetUserGrades - Получение оценок конкретного пользователя
func GetUserGrades(c *gin.Context) {
	studentID := c.Query("student_id")
	if studentID == "" {
		studentID = c.GetString("firebase_uid")
	}

	iter := repositories.FirestoreClient.Collection("grades").Where("student_id", "==", studentID).OrderBy("created_at", repositories.Descending).Documents(c.Request.Context())
	
	var grades []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch grades"})
			return
		}
		
		data := doc.Data()
		data["id"] = doc.Ref.ID
		grades = append(grades, data)
	}

	c.JSON(http.StatusOK, grades)
}
