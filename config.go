package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func readConfig() error {

	// Env variables
	viper.SetEnvPrefix("BF_DISCORD_")
	viper.AutomaticEnv()

	// Setup config file lookup
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Setup defaults
	viper.SetDefault("bot_prefix", "!bf")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("config file not found: %v", err)
		}

		return fmt.Errorf("config file found, but another error ocurred: %v", err)

	}

	token := viper.GetString("bot_token")
	if token == "" {
		return fmt.Errorf("bot token is required: create a config.yml in the program's " +
			"directory with a key 'bot_token' with the token or define an env variable BF_DISCORD_BOT_TOKEN containing the token")
	}

	return nil
}
