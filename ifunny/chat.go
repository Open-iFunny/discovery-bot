package ifunny

import (
	"fmt"

	"github.com/jcelliott/turnpike"
	"github.com/mitchellh/mapstructure"
)

const (
	chatRoot      = "wss://chat.ifunny.co/chat"
	chatNamespace = "co.fun.chat"
)

type Chat interface {
	Chats(userID string) <-chan *WSChat
	Invites(userID string) <-chan *WSInvite
}

type chat struct {
	ws     *turnpike.Client
	bearer string
	hello  map[string]interface{}
}

func topic(name string) string { return chatNamespace + "." + name }

func connectChat(bearer string) (Chat, error) {
	ws, err := turnpike.NewWebsocketClient(turnpike.JSON, chatRoot, nil, nil, nil)
	if err != nil {
		panic(err)
	}

	ws.Auth = map[string]turnpike.AuthFunc{
		"ticket": turnpike.NewTicketAuthenticator(bearer),
	}
	hello, err := ws.JoinRealm(topic("ifunny"), nil)
	if err != nil {
		panic(err)
	}

	return &chat{ws, bearer, hello}, nil
}

func (chat *chat) Chats(userID string) <-chan *WSChat {
	result := make(chan *WSChat)
	chat.ws.Subscribe(topic("user."+userID+".chats"), nil, func(_ []interface{}, kwargs map[string]interface{}) {
		if kwargs["chats"] == nil {
			return
		}

		for _, chatRaw := range kwargs["chats"].([]interface{}) {
			wsChat := new(WSChat)
			mapstructure.Decode(chatRaw, wsChat)
			result <- wsChat
		}
	})

	return result
}

func (chat *chat) Invites(userID string) <-chan *WSInvite {
	result := make(chan *WSInvite)
	chat.ws.Subscribe(topic("user."+userID+".invites"), nil, func(_ []interface{}, kwargs map[string]interface{}) {
		if kwargs["invites"] == nil {
			return
		}

		for _, invRaw := range kwargs["invites"].([]interface{}) {
			fmt.Printf("invite: %+v\n", invRaw)

			wsInvite := new(WSInvite)
			mapstructure.Decode(invRaw, wsInvite)
			result <- wsInvite
		}
	})

	return result
}
