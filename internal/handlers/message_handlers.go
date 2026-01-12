package handlers

import (
	"net/http"
	"rio/internal/models"
	"rio/internal/store"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func SendMessage(c *gin.Context) {
	channel_id, err := strconv.Atoi(c.Param("channel_id"))
	if err != nil || channel_id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "channel id must be a positive integer"})
	}

	var newMessage models.Message
	if err := c.ShouldBindJSON(&newMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	newMessage.ChannelID = channel_id

	if strings.TrimSpace(newMessage.Content) == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "message content must not be empty"})
		return
	}

	newMessage.ID = store.GetNextMessageId()
	store.Messages = append(store.Messages, newMessage)
	c.JSON(http.StatusCreated, newMessage)
}

func GetMessages(c *gin.Context) {
	channel_id, err := strconv.Atoi(c.Param("channel_id"))

	if err != nil || channel_id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "server id must be a valid non-negative integer"})
		return
	}

	var filteredMessages []models.Message
	for _, message := range store.Messages {
		if message.ChannelID == channel_id {
			filteredMessages = append(filteredMessages, message)
		}
	}

	if len(filteredMessages) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no messages in that channel"})
		return
	}

	c.JSON(http.StatusOK, filteredMessages)
}
