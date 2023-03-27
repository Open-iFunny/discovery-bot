package ifunny

import (
	"fmt"

	"github.com/jcelliott/turnpike"
	"github.com/mitchellh/mapstructure"
)

type WSChat struct {
	Name  string `mapstructure:"name"`
	Title string `mapstructure:"title"`
}

type WSChats struct {
	Chats string `mapstructure:"chats"`
}

type WSInvite struct {
}

func (client *Client) Chat() (*Chat, error) {
	ws, err := turnpike.NewWebsocketClient(turnpike.JSON, chatRoot, nil, nil, nil)
	if err != nil {
		panic(err)
	}

	ws.Auth = map[string]turnpike.AuthFunc{
		"ticket": turnpike.NewTicketAuthenticator(client.bearer),
	}
	hello, err := ws.JoinRealm(topic("ifunny"), nil)
	if err != nil {
		panic(err)
	}

	return &Chat{ws, client, hello}, nil
}

type Chat struct {
	ws     *turnpike.Client
	client *Client
	hello  map[string]interface{}
}

func topic(name string) string { return chatNamespace + "." + name }

func (chat *Chat) Chats() <-chan *WSChat {
	result := make(chan *WSChat)
	uri := topic("user." + chat.client.self.ID + ".chats")
	chat.ws.Subscribe(uri, nil, func(_ []interface{}, kwargs map[string]interface{}) {
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

func (chat *Chat) Invites() <-chan *WSInvite {
	result := make(chan *WSInvite)
	uri := topic("user." + chat.client.self.ID + ".invites")
	chat.ws.Subscribe(uri, nil, func(_ []interface{}, kwargs map[string]interface{}) {
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
