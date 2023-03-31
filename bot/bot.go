package bot

import (
	"github.com/gastrodon/popplio/ifunny"
	"github.com/gastrodon/popplio/ifunny/compose"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type filtHandler struct {
	filter filter
	handle handler
}

type Bot struct {
	Client *ifunny.Client
	Chat   *ifunny.Chat
	log    *logrus.Logger

	recvEvents   chan map[string]interface{}
	unsubEvents  map[string]func()
	handleEvents map[string]filtHandler
}

func MakeBot(bearer, userAgent string) (*Bot, error) {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(ifunny.LogLevel)
	client, err := ifunny.MakeClientLog(bearer, userAgent, log)
	if err != nil {
		return nil, err
	}

	chat, err := client.Chat()
	if err != nil {
		return nil, err
	}

	return &Bot{
		Client:       client,
		Chat:         chat,
		log:          log,
		recvEvents:   make(chan map[string]interface{}),
		unsubEvents:  make(map[string]func()),
		handleEvents: make(map[string]filtHandler, 0),
	}, nil
}

func (bot *Bot) Subscribe(channel string) {
	traceID := uuid.New().String()
	if unsub, ok := bot.unsubEvents[channel]; ok {
		bot.log.WithFields(logrus.Fields{
			"trace_id":     traceID,
			"channel_name": channel,
		}).Warn("SubscribeChat on subscribed channel")

		unsub()
	}

	bot.Chat.Subscribe(compose.EventsIn(channel), func(eventType int, event map[string]interface{}) error {
		bot.recvEvents <- event
		return nil
	})
}

func (bot *Bot) Unsubscribe(channel string) {
	traceID := uuid.New().String()
	if unsub, ok := bot.unsubEvents[channel]; !ok {
		bot.log.WithFields(logrus.Fields{
			"trace_id":     traceID,
			"channel_name": channel,
		}).Warn("UnsubscribeChat on not subscribed channel")
	} else {
		unsub()
	}
}
