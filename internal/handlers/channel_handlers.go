package handlers

import (
	"net/http"
	"rio/internal/models"
	"rio/internal/store"

	"github.com/gin-gonic/gin"
)

func CreateChannel(c *gin.Context) {
	var newChannel models.Channel
	if err := c.ShouldBindJSON(&newChannel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	newChannel.ID = store.GetNextChannelId()

	store.Channels = append(store.Channels, newChannel)
	c.JSON(http.StatusCreated, newChannel)
}

func GetChannels(c *gin.Context) {
	serverIDStr := c.Param("server_id")

	var filteredChannels []models.Channel
	for _, channel := range store.Channels {
		if channel.ServerID == serverIDStr {
			filteredChannels = append(filteredChannels, channel)
		}
	}

	if len(filteredChannels) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no channels found for this server"})
		return
	}

	c.JSON(http.StatusOK, filteredChannels)
}
