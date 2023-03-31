package ifunny

import (
	"github.com/gastrodon/popplio/ifunny/compose"
)

const (
	UNK_0 messageType = iota
	MESSAGE
	UNK_2
	JOIN_CHANNEL
	EXIT_CHANNEL
)

type messageType int

type ChatEvent struct {
	ID   string `json:"id"`
	Text string `json:"text"`

	Type   messageType `json:"type"`
	Status int         `json:"status"`
	PubAt  float64     `json:"pub_at"`

	User struct {
		ID         string `json:"user"`
		Nick       string `json:"nick"`
		IsVerified bool   `json:"is_verified"`
		LastSeenAt int64  `json:"last_seen_at"`
	} `json:"user"`
}

func (chat *Chat) OnChanneEvent(channel string, handle func(event *ChatEvent) error) (func(), error) {
	return chat.Subscribe(compose.EventsIn(channel), func(eventType int, kwargs map[string]interface{}) error {
		log := chat.client.log.WithField("event_type", eventType)

		if kwargs["message"] == nil {
			log.WithField("kwargs", kwargs).Warn("channel event message is nil")
		}

		message := new(ChatEvent)
		if err := jsonDecode(kwargs["message"], message); err != nil {
			return err
		}

		return handle(message)
	})
}
