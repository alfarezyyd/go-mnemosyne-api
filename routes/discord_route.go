package routes

import (
	"github.com/bwmarrin/discordgo"
	"go-mnemosyne-api/discord"
)

func DiscordRoutes(discordSession *discordgo.Session, discordController discord.Controller) {
	discordSession.AddHandler(discordController.OnMessageCreate)
}
