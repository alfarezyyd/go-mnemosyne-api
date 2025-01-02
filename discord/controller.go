package discord

import "github.com/bwmarrin/discordgo"

type Controller interface {
	OnMessageCreate(discSession *discordgo.Session, messagePayload *discordgo.MessageCreate)
}
