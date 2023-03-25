package main

import (
	"fmt"
	"io"
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

	client := ifunny.MakeClient("bearer "+bearer, userAgent)
	response, err := client.Request("GET", "/account", nil)
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", data)
}
