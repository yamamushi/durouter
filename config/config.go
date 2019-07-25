package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

// ReadConfig function
func ReadConfig(path string) (config Config, err error) {

	if _, err := toml.DecodeFile(path, &config); err != nil {
		fmt.Println(err)
		return config, err
	}

	return config, nil
}

// Config struct
type Config struct {
	HostConfig    hostConfig     `toml:"host"`
	DiscordConfig discordConfig  `toml:"discord"`
	DBConfig      databaseConfig `toml:"database"`
}

// hostConfig struct
type hostConfig struct {
	Port    string `toml:"port"`
	Host    string `toml:"host"`
}

// discordConfig struct
type discordConfig struct {
	Token   string `toml:"bot_token"`
	GuildID string `toml:"guild_id"`
	AuthID  string `toml:"authorized_role_id"`
	OauthClientID string `toml:"oauth_client_id"`
	OauthClientSecret string `toml:"oauth_client_secret"`
	ServerID string `toml:"server_id"`
}

// databaseConfig struct
type databaseConfig struct {
	DBFile string `toml:"filename"`
	MongoHost string `toml:"mongohost"`
	MongoDB string  `toml:"mongodb"`
	MongoUser string    `toml:"mongouser"`
	MongoPass string    `toml:"mongopass"`
	AccountColumn string `toml:"accountcolumn"`
	CraftingRecordColumn string `toml:"craftingrecordcolumn"`
}