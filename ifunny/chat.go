package ifunny

import (
	"github.com/jcelliott/turnpike"
)

const chatRoot = "wss://chat.ifunny.co/chat"

type Chat interface {
	Subscribe(id string) (<-chan interface{}, func())
	Publish(topic string, event interface{}) error
}

type chat struct {
	ws     *turnpike.Client
	bearer string
	hello  map[string]interface{}
}

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

func (chat *chat) Subscribe(topic string) (<-chan interface{}, func()) {
	panic("unimplemented")
}

func (chat *chat) Publish(topic string, event interface{}) error {
	panic("unimplemented")
}
