package ifunny

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type ChatMessage struct {
	ID   string `mapstructure:"id"`
	Text string `mapstructure:"text"`

	Type   int     `mapstructure:"type"` // 1 = message, ???
	Status int     `mapstructure:"status"`
	PubAt  float64 `mapstructure:"pub_at"`

	User struct {
		ID         string `mapstructure:"user"`
		Nick       string `mapstructure:"nick"`
		IsVerified bool   `mapstructure:"is_verified"`
		LastSeenAt int64  `mapstructure:"last_seen_at"`
	} `mapstructure:"user"`
}

type sMessage subscribe

func MessageIn(channel string) sMessage {
	return sMessage{
		topic:   uri("chat." + channel),
		options: nil,
	}
}

func (chat *Chat) IterMessage(desc sMessage) <-chan *ChatMessage {
	result := make(chan *ChatMessage)
	chat.ws.Subscribe(desc.topic, desc.options, func(opts []interface{}, kwargs map[string]interface{}) {
		if kwargs["message"] == nil {
			fmt.Printf("got non-message data %+v\n", kwargs)
			return
		}

		message := new(ChatMessage)
		if err := mapstructure.Decode(kwargs["message"], message); err != nil {
			panic(err)
		}
		result <- message
	})

	return result
}
