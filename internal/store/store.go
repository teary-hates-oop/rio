package store

import "rio/internal/models"

var (
	Users       = []models.User{}
	Servers     = []models.Server{}
	Channels    = []models.Channel{}
	Messages    = []models.Message{}
	UserServers = []models.UserServer{}

	nextUserID    = 1
	nextServerID  = 1
	nextChannelID = 1
	nextMessageID = 1
)

func GetNextUserId() int {
	id := nextUserID
	nextUserID++
	return id
}

func GetNextServerId() int {
	id := nextServerID
	nextServerID++
	return id
}

func GetNextChannelId() uint {
	id := uint(nextChannelID)
	nextChannelID++
	return id
}

func GetNextMessageId() uint {
	id := uint(nextMessageID)
	nextMessageID++
	return id
}
