package discord

import (
	"encoding/json"
	"github.com/yamamushi/durouter/config"
	"github.com/yamamushi/durouter/models"
	"io/ioutil"
	"net/http"
)
type DiscordAPI struct {

	botToken              string
	clientID              string
	clientSecret          string
	alphaAuthorizedRoleID string
	guildID               string

}

func NewDiscordAPI(config config.Config) *DiscordAPI {

	result := &DiscordAPI{}
	result.botToken = config.DiscordConfig.Token
	result.guildID = config.DiscordConfig.ServerID
	result.alphaAuthorizedRoleID = config.DiscordConfig.AuthID
	result.clientID = config.DiscordConfig.OauthClientID
	result.clientSecret = config.DiscordConfig.OauthClientSecret

	return result
}

func (h *DiscordAPI) GetUserInfo(UserID string) (result models.AccountInfo) {

	result.Error = true
	result.Status = "500"
	result.AlphaAuthorized = false

	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://discordapp.com/api/v6/guilds/"+h.guildID+"/members/"+UserID, nil)
	if err != nil {
		return result
	}
	request.Header.Set("Authorization", "Bot "+h.botToken)
	resp, err := client.Do(request)
	if err != nil {
		return result
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result
	}

	discordUser, err := h.getDiscordUser(body)
	if err != nil {
		return result
	}

	result.DiscordID = discordUser.User.ID

	for _, role := range discordUser.Roles {
	 	if role == h.alphaAuthorizedRoleID {
	 		result.AlphaAuthorized = true
	    }
	}

	result.Error = false
	result.Status = "200"

	return result
}

func (h *DiscordAPI) getDiscordUser(body []byte) (*models.DiscordUser, error) {
	var s = new(models.DiscordUser)
	err := json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}
	return s, err
}
