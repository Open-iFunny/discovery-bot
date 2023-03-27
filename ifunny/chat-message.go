package ifunny

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
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
	traceID := uuid.New().String()
	chat.client.log.WithFields(logrus.Fields{
		"trace_id": traceID,
		"topic":    desc.topic,
		"options":  desc.options,
	}).Trace("iter message")

	result := make(chan *ChatMessage)
	chat.ws.Subscribe(desc.topic, desc.options, func(opts []interface{}, kwargs map[string]interface{}) {
		if kwargs["message"] == nil {
			fmt.Printf("got non-message data %+v\n", kwargs)
			mType := 0
			if err := mapstructure.Decode(kwargs["type"], &mType); err != nil {
				chat.client.log.WithField("trace_id", traceID).Error(err)
				return
			}

			chat.client.log.WithFields(logrus.Fields{
				"trace_id": traceID,
				"type":     mType,
			}).Warn("unknown message payload ", kwargs)
			return
		}

		message := new(ChatMessage)
		if err := mapstructure.Decode(kwargs["message"], message); err != nil {
			chat.client.log.WithField("trace_id", traceID).Error(err)
			return
		}

		chat.client.log.WithFields(logrus.Fields{
			"trace_id":     traceID,
			"message_id":   message.ID,
			"message_from": message.User.Nick,
			"message_text": message.Text,
		})
		result <- message
	})

	return result
}
