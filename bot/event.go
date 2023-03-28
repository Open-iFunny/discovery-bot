package bot

import (
	"github.com/gastrodon/popplio/ifunny"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type filter func(event ifunny.Event) bool
type handler func(event ifunny.Event) error

func (bot *Bot) On(filter filter, handle handler) func() {
	eventID := uuid.New().String()
	bot.handleEvents[eventID] = filtHandler{filter, handle}

	return func() { delete(bot.handleEvents, eventID) }
}

func (bot *Bot) Listen() {
	for event := range bot.recvEvents {
		go func(handlers map[string]filtHandler, event ifunny.Event) {
			log := bot.log.WithFields(logrus.Fields{
				"event_type": event.Type(),
			})

			log.Trace("start handling")
			for id, filtHandle := range handlers {
				if filtHandle.filter(event) {
					if err := filtHandle.handle(event); err != nil {
						log.WithField("handle_id", id).Error(err)
					}
				}
			}
		}(bot.handleEvents, event)
	}
}
