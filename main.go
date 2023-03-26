package main

import (
	"os"
	"time"

	"github.com/gastrodon/popplio/ifunny"
)

func main() {
	bearer := os.Getenv("IFUNNY_BEARER")
	if bearer == "" {
		panic("IFUNNY_BEARER must be set")
	}

	userAgent := os.Getenv("IFUNNY_USER_AGENT")
	if userAgent == "" {
		panic("IFUNNY_USER_AGENT must be set")
	}

	cookie := os.Getenv("IFUNNY_CHAT_COOKIE")
	if userAgent == "" {
		panic("IFUNNY_CHAT_COOKIE must be set")
	}

	client := ifunny.MakeClient("bearer "+bearer, userAgent)

	client.Request("GET", "/v4", nil)

	_, err := client.Connect(bearer, cookie)
	if err != nil {
		panic(err)
	}

	<-time.After(4 * time.Second)
}
