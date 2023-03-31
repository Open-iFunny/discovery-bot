package main

import (
	"os"

	"github.com/gastrodon/popplio/bot"
)

var bearer = ""
var userAgent = ""

func init() {
	bearer = os.Getenv("IFUNNY_BEARER")
	if bearer == "" {
		panic("IFUNNY_BEARER must be set")
	}

	userAgent = os.Getenv("IFUNNY_USER_AGENT")
	if userAgent == "" {
		panic("IFUNNY_USER_AGENT must be set")
	}
}

func main() {
	bot, err := bot.MakeBot(bearer, userAgent)
	if err != nil {
		panic(err)
	}

	bot.Listen()
}
