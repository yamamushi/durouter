package controllers

import (
	"encoding/json"
	"github.com/yamamushi/durouter/config"
	"github.com/yamamushi/durouter/discord"
	"net/http"
	"strings"
)

type VerifyAuth struct {
	config config.Config
}

func NewVerifyAuth(config config.Config) (*VerifyAuth){
	verifyauth := &VerifyAuth{}
	verifyauth.config = config
	return verifyauth
}

func (h *VerifyAuth) Verify(w http.ResponseWriter, r *http.Request) {
	// Tell our client to expect a json response
	w.Header().Set("Content-Type", "application/json")

	userID := r.URL.Path
	userID = strings.TrimPrefix(userID, "/verifyauth/")

	discordapi := discord.NewDiscordAPI(h.config)
	account := discordapi.GetUserInfo(userID)
	account.DiscordID = userID
	if account.Error {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	_ = json.NewEncoder(w).Encode(account)
}
