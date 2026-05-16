package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/itstep-backend/internal/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	defaultChatLimit = 100
	maxChatLimit     = 300
	maxMessageLength = 10000
)

type chatMessageDocument struct {
	ID         bson.ObjectID `bson:"_id,omitempty"`
	ChatID     string        `bson:"chat_id"`
	SenderID   string        `bson:"sender_id"`
	ReceiverID string        `bson:"receiver_id"`
	Text       string        `bson:"text"`
	Read       bool          `bson:"read"`
	CreatedAt  time.Time     `bson:"created_at"`
}

func ensureMongoChatReady(c *gin.Context) bool {
	if repositories.IsMongoChatReady() {
		return true
	}

	c.JSON(http.StatusServiceUnavailable, gin.H{
		"error": "Chat service is unavailable. MONGO_URL is not configured or MongoDB is unreachable.",
	})
	return false
}

func buildChatID(a, b string) string {
	users := []string{strings.TrimSpace(a), strings.TrimSpace(b)}
	sort.Strings(users)
	return users[0] + "_" + users[1]
}

func parseChatPeerID(c *gin.Context) string {
	userID := strings.TrimSpace(c.Query("user_id"))
	if userID == "" {
		userID = strings.TrimSpace(c.Query("uid"))
	}
	return userID
}

func toChatMessageResponse(doc chatMessageDocument) gin.H {
	return gin.H{
		"id":         doc.ID.Hex(),
		"chatId":     doc.ChatID,
		"senderId":   doc.SenderID,
		"receiverId": doc.ReceiverID,
		"text":       doc.Text,
		"read":       doc.Read,
		"createdAt":  doc.CreatedAt,
	}
}

func GetChatMessages(c *gin.Context) {
	if !ensureMongoChatReady(c) {
		return
	}

	requesterID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	otherUserID := parseChatPeerID(c)
	if otherUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	if otherUserID == requesterID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot open chat with yourself"})
		return
	}

	limit := defaultChatLimit
	if rawLimit := strings.TrimSpace(c.Query("limit")); rawLimit != "" {
		n, convErr := strconv.Atoi(rawLimit)
		if convErr != nil || n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a positive integer"})
			return
		}
		if n > maxChatLimit {
			n = maxChatLimit
		}
		limit = n
	}

	chatID := buildChatID(requesterID, otherUserID)

	filter := bson.M{
		"chat_id": chatID,
	}

	findOpts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := repositories.ChatMessagesCollection.Find(c.Request.Context(), filter, findOpts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat messages"})
		return
	}
	defer cursor.Close(c.Request.Context())

	var docs []chatMessageDocument
	if err := cursor.All(c.Request.Context(), &docs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode chat messages"})
		return
	}

	result := make([]gin.H, 0, len(docs))
	for i := len(docs) - 1; i >= 0; i-- {
		result = append(result, toChatMessageResponse(docs[i]))
	}

	c.JSON(http.StatusOK, result)
}

func SendChatMessage(c *gin.Context) {
	if !ensureMongoChatReady(c) {
		return
	}

	senderID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	receiverID := strings.TrimSpace(asString(input["receiver_id"]))
	if receiverID == "" {
		receiverID = strings.TrimSpace(asString(input["receiverId"]))
	}

	text := strings.TrimSpace(asString(input["text"]))

	if receiverID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "receiver_id is required"})
		return
	}
	if receiverID == senderID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot send message to yourself"})
		return
	}
	if text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message text cannot be empty"})
		return
	}
	if len([]rune(text)) > maxMessageLength {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Message is too long. Maximum is %d characters", maxMessageLength)})
		return
	}

	now := time.Now().UTC()
	doc := chatMessageDocument{
		ChatID:     buildChatID(senderID, receiverID),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Text:       text,
		Read:       false,
		CreatedAt:  now,
	}

	insertResult, err := repositories.ChatMessagesCollection.InsertOne(c.Request.Context(), doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
		return
	}

	insertedID, ok := insertResult.InsertedID.(bson.ObjectID)
	if ok {
		doc.ID = insertedID
	}

	c.JSON(http.StatusCreated, toChatMessageResponse(doc))
}

func MarkChatAsRead(c *gin.Context) {
	if !ensureMongoChatReady(c) {
		return
	}

	requesterID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	otherUserID := strings.TrimSpace(asString(input["user_id"]))
	if otherUserID == "" {
		otherUserID = strings.TrimSpace(asString(input["userId"]))
	}

	if otherUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	if otherUserID == requesterID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot mark own messages as read"})
		return
	}

	filter := bson.M{
		"chat_id":     buildChatID(requesterID, otherUserID),
		"sender_id":   otherUserID,
		"receiver_id": requesterID,
		"read":        false,
	}
	update := bson.M{
		"$set": bson.M{
			"read": true,
		},
	}

	result, err := repositories.ChatMessagesCollection.UpdateMany(c.Request.Context(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark chat messages as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"matched":  result.MatchedCount,
		"modified": result.ModifiedCount,
	})
}

func GetChatUnreadCount(c *gin.Context) {
	if !ensureMongoChatReady(c) {
		return
	}

	requesterID, err := getRequesterUID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	filter := bson.M{
		"receiver_id": requesterID,
		"read":        false,
	}

	count, err := repositories.ChatMessagesCollection.CountDocuments(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate unread messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
