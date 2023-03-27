package main

import (
	"fmt"
	"os"

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

	client, _ := ifunny.MakeClient(bearer, userAgent)

	chat, _ := client.Chat()

	dm, e, err := chat.GetDM("5d19bdf524aac73ffb2b2e81")
	if err != nil {
		panic(err)
	}

	fmt.Printf("exists: %t, %+v", e, dm)
}
