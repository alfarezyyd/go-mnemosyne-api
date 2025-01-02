package config

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

type DiscordClient struct {
	viperConfig *viper.Viper
}

func NewDiscordClient(viperConfig *viper.Viper) *DiscordClient {
	return &DiscordClient{viperConfig: viperConfig}
}
func (discordClient *DiscordClient) InitializeDiscordConnection() (*discordgo.Session, error) {
	sess, err := discordgo.New(fmt.Sprintf("Bot %s", discordClient.viperConfig.GetString("DISCORD_OAUTH2_TOKEN")))
	if err != nil {
		return nil, err
	}
	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	err = sess.Open()
	if err != nil {
		panic(err)
	}

	return sess, nil
}
