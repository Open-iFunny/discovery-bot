package ifunny

import (
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

const (
	UNK_0 messageType = iota
	MESSAGE
	UNK_2
	JOIN_CHANNEL
	EXIT_CHANNEL
)

type messageType int

type ChatMessage struct {
	ID   string `mapstructure:"id"`
	Text string `mapstructure:"text"`

	Type   messageType `mapstructure:"type"` // 1 = message, ???
	Status int         `mapstructure:"status"`
	PubAt  float64     `mapstructure:"pub_at"`

	User struct {
		ID         string `mapstructure:"user"`
		Nick       string `mapstructure:"nick"`
		IsVerified bool   `mapstructure:"is_verified"`
		LastSeenAt int64  `mapstructure:"last_seen_at"`
	} `mapstructure:"user"`
}

type sEvent subscribe

func MessageIn(channel string) sEvent {
	return sEvent{
		topic:   uri("chat." + channel),
		options: nil,
	}
}

func (chat *Chat) SubscribeEvent(desc sEvent, handle func(WSResource) error) func() {
	traceID := uuid.New().String()
	chat.client.log.WithFields(logrus.Fields{
		"trace_id": traceID,
		"topic":    desc.topic,
		"options":  desc.options,
	}).Trace("begin subscribe message")

	chat.ws.Subscribe(desc.topic, desc.options, func(opts []interface{}, kwargs map[string]interface{}) {
		eType := 0
		if err := mapstructure.Decode(kwargs["type"], &eType); err != nil {
			chat.client.log.WithField("trace_id", traceID).Error(err)
		}

		chat.client.log.WithFields(logrus.Fields{
			"trace_id":   traceID,
			"event_type": eType,
			"event_data": kwargs,
		}).Trace("handle event")

		if err := handle(makeEvent(eType, kwargs)); err != nil {
			chat.client.log.WithField("trace_id", traceID).Error("subscribe message handler: " + err.Error())
		}
	})

	return func() { chat.ws.Unsubscribe(desc.topic) }
}

func (chat *Chat) IterEvent(desc sEvent) (<-chan WSResource, func()) {
	traceID := uuid.New().String()
	chat.client.log.WithFields(logrus.Fields{
		"trace_id": traceID,
		"topic":    desc.topic,
		"options":  desc.options,
	}).Trace("begin iter message")

	result := make(chan WSResource)
	return result, chat.SubscribeEvent(desc, func(chat WSResource) error {
		result <- chat
		return nil
	})
}

type pMessage publish

func MessageTo(channel, text string) pMessage {
	return pMessage{
		topic:   uri("chat." + channel),
		options: map[string]interface{}{"acknowledge": 1, "exclude_me": 1},
		args:    nil,
		kwargs: map[string]interface{}{
			"message_type": 1,
			"type":         200,
			"text":         text,
		},
	}
}

func (chat *Chat) SendMessage(desc pMessage) error {
	return chat.ws.Publish(desc.topic, desc.options, desc.args, desc.kwargs)
}
