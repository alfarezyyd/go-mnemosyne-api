package config

import (
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type IdentityProvider struct {
	viperConfig          *viper.Viper
	GoogleProviderConfig oauth2.Config
}

func NewIdentityProvider(viperConfig *viper.Viper) *IdentityProvider {
	return &IdentityProvider{
		viperConfig: viperConfig,
	}
}

func (identityProvider *IdentityProvider) InitializeGoogleProviderConfig() oauth2.Config {

	identityProvider.GoogleProviderConfig = oauth2.Config{
		RedirectURL:  identityProvider.viperConfig.GetString("GOOGLE_REDIRECT_CALLBACK"),
		ClientID:     identityProvider.viperConfig.GetString("GOOGLE_CLIENT_ID"),
		ClientSecret: identityProvider.viperConfig.GetString("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return identityProvider.GoogleProviderConfig
}
