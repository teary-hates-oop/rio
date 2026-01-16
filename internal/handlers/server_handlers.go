package handlers

import (
	"net/http"
	"rio/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type ServerHandler struct {
	service *service.ServerService
}

func NewServerHandler(svc *service.ServerService) *ServerHandler {
	return &ServerHandler{service: svc}
}

func (h *ServerHandler) CreateServer(c *gin.Context) {
	currentUserID := c.GetString("user_id")
	if currentUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	server, err := h.service.CreateServer(currentUserID, input.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, server)
}

func (h *ServerHandler) GetServers(c *gin.Context) {
	currentUserID := c.GetString("user_id")
	if currentUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	servers, err := h.service.ListUserServers(currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, servers)
}

func (h *ServerHandler) GetServer(c *gin.Context) {
	currentUserID := c.GetString("user_id")
	if currentUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	serverID := c.Param("id")
	if serverID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "server ID is required"})
		return
	}

	server, err := h.service.GetServer(currentUserID, serverID)
	if err != nil {
		if err.Error() == "user is not a member of server" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "server not found"})
		return
	}

	c.JSON(http.StatusOK, server)
}

func (h *ServerHandler) UpdateServer(c *gin.Context) {
	currentUserID := c.GetString("user_id")
	if currentUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	serverID := c.Param("id")
	if serverID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "server ID is required"})
		return
	}

	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.UpdateServerName(currentUserID, serverID, input.Name)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "permissions") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *ServerHandler) DeleteServer(c *gin.Context) {
	currentUserID := c.GetString("user_id")
	if currentUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	serverID := c.Param("id")
	if serverID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "server ID is required"})
		return
	}

	err := h.service.DeleteServer(currentUserID, serverID)
	if err != nil {
		if strings.Contains(err.Error(), "permissions") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "server not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ServerHandler) AddMember(c *gin.Context) {
	currentUserID := c.GetString("user_id")
	if currentUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	serverID := c.Param("id")
	if serverID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "server ID is required"})
		return
	}

	var input struct {
		UserID string `json:"userId" binding:"required"`
		Role   string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.AddMember(currentUserID, serverID, input.UserID, input.Role)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "permissions") || strings.Contains(err.Error(), "insufficient") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *ServerHandler) RemoveMember(c *gin.Context) {
	currentUserID := c.GetString("user_id")
	if currentUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	serverID := c.Param("id")
	targetUserID := c.Param("userId")

	if serverID == "" || targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "server ID and user ID are required"})
		return
	}

	err := h.service.RemoveMember(currentUserID, serverID, targetUserID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "permissions") || strings.Contains(err.Error(), "insufficient") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ServerHandler) ChangeMemberRole(c *gin.Context) {
	currentUserID := c.GetString("user_id")
	if currentUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	serverID := c.Param("id")
	targetUserID := c.Param("userId")

	if serverID == "" || targetUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "server ID and user ID are required"})
		return
	}

	var input struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.ChangeMemberRole(currentUserID, serverID, targetUserID, input.Role)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "permissions") || strings.Contains(err.Error(), "insufficient") {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
