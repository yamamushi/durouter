package models

import "time"

type DiscordUser struct {
	Nick interface{} `json:"nick"`
	User struct {
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		ID            string `json:"id"`
		Avatar        string `json:"avatar"`
	} `json:"user"`
	Roles        []string  `json:"roles"`
	PremiumSince time.Time `json:"premium_since"`
	Deaf         bool      `json:"deaf"`
	Mute         bool      `json:"mute"`
	JoinedAt     time.Time `json:"joined_at"`
}