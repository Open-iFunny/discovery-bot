package ifunny

import (
	"net/http"

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
}

func connectChat(bearer, cookie string) (Chat, error) {
	header := http.Header{"Cookie": []string{cookie}}
	ws, err := turnpike.NewWebsocketClient(turnpike.JSON, chatRoot, header, nil, nil)
	if err != nil {
		panic(err)
	}

	return &chat{ws, bearer}, nil
}

func (chat *chat) Subscribe(topic string) (<-chan interface{}, func()) {
	panic("unimplemented")
}

func (chat *chat) Publish(topic string, event interface{}) error {
	panic("unimplemented")
}
