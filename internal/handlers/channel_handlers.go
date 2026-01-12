package handlers

import (
	"net/http"
	"rio/internal/models"
	"rio/internal/store"
	"strconv"

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
	server_id, err := strconv.Atoi(c.Param("server_id"))
	if err != nil || server_id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "server id must be a valid non-negative integer"})
		return
	}

	var filteredChannels []models.Channel
	for _, channel := range store.Channels {
		if channel.ServerID == server_id {
			filteredChannels = append(filteredChannels, channel)
		}
	}

	if len(filteredChannels) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no channels with id"})
		return
	}

	c.JSON(http.StatusOK, filteredChannels)
}
