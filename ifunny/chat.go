package ifunny

import "github.com/gastrodon/turnpike"

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
	ws := turnpike.NewClient()
	if err := ws.Connect(chatRoot, "http://localghost", cookie); err != nil {
		panic("connect: " + err.Error())
	}

	return &chat{ws, bearer}, nil
}

func (chat *chat) Subscribe(topic string) (<-chan interface{}, func()) {
	panic("unimplemented")
}

func (chat *chat) Publish(topic string, event interface{}) error {
	panic("unimplemented")
}
