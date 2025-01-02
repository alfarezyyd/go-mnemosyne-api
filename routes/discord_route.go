package routes

import (
	"github.com/bwmarrin/discordgo"
	"go-mnemosyne-api/discord"
)

func DiscordRoutes(discordSession *discordgo.Session) {
	discordSession.AddHandler(discord.OnMessageCreate)
}
