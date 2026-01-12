package handlers

import (
	"net/http"
	"rio/internal/models"
	"rio/internal/store"

	"github.com/gin-gonic/gin"
)

func CreateServer(c *gin.Context) {
	var newServer models.Server
	if err := c.ShouldBindJSON(&newServer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	newServer.ID = store.GetNextServerId()
	store.Servers = append(store.Servers, newServer)
	c.JSON(http.StatusCreated, newServer)
}

func GetServers(c *gin.Context) {
	c.JSON(http.StatusOK, store.Servers)
}
