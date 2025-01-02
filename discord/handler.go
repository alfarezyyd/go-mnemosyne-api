package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type Handler struct {
	discordService Service
}

func NewHandler(discordService Service) *Handler {
	return &Handler{
		discordService: discordService,
	}
}

func (discordHandler *Handler) OnMessageCreate(discSession *discordgo.Session, messagePayload *discordgo.MessageCreate) {
	if messagePayload.Author.ID == discSession.State.User.ID {
		return
	}

	if messagePayload.Content == "ping" {
		discSession.ChannelMessageSend(messagePayload.ChannelID, "Pong!")
	}

	fmt.Println("ONLINe!")
}
