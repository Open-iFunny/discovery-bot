package ifunny

import (
	"github.com/jcelliott/turnpike"
	"github.com/mitchellh/mapstructure"
)

const (
	chatRoot      = "wss://chat.ifunny.co/chat"
	chatNamespace = "co.fun.chat"
)

type Chat interface {
	Chats(userID string) <-chan *WSChat
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
	chat.ws.Subscribe(topic("user."+userID+".chats"), nil, func(args []interface{}, kwargs map[string]interface{}) {
		for _, chatRaw := range kwargs["chats"].([]interface{}) {
			wsChat := new(WSChat)
			mapstructure.Decode(chatRaw, wsChat)
			result <- wsChat
		}

		close(result)
	})

	return result
}

func (chat *chat) Invites() {

}
