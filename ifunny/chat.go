package ifunny

import (
	"fmt"

	"github.com/gastrodon/popplio/ifunny/compose"
	"github.com/google/uuid"
	"github.com/jcelliott/turnpike"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

const (
	chatRoot = "wss://chat.ifunny.co/chat"
)

func (client *Client) Chat() (*Chat, error) {
	log := client.log.WithField("trace_id", uuid.New().String())

	log.Trace("start connect chat")
	ws, err := turnpike.NewWebsocketClient(turnpike.JSON, chatRoot, nil, nil, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.Trace("join realm ifunny")
	ws.Auth = map[string]turnpike.AuthFunc{"ticket": turnpike.NewTicketAuthenticator(client.bearer)}
	hello, err := ws.JoinRealm(string(compose.URI("ifunny")), nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &Chat{ws, client, hello}, nil
}

type Chat struct {
	ws     *turnpike.Client
	client *Client
	hello  map[string]interface{}
}

func (chat *Chat) call(desc turnpike.Call, output interface{}) error {
	log := chat.client.log.WithFields(logrus.Fields{
		"trace_id": uuid.New().String(),
		"type":     "CALL",
		"uri":      desc.Procedure,
		"kwargs":   desc.ArgumentsKw,
	})

	log.Trace("exec call")
	result, err := chat.ws.Call(string(desc.Procedure), desc.Options, desc.Arguments, desc.ArgumentsKw)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Trace(fmt.Sprintf("call OK recv: %+v\n", result.ArgumentsKw))
	if output != nil {
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			Result:  output,
			TagName: "json",
		})

		if err != nil {
			log.Error(err)
			return err
		}

		err = decoder.Decode(result.ArgumentsKw)
		if err != nil {
			log.Error(err)
		}

		return err
	}

	return nil
}

func (chat *Chat) publish(desc turnpike.Publish) error {
	return chat.ws.Publish(string(desc.Topic), desc.Options, desc.Arguments, desc.ArgumentsKw)
}

func (chat *Chat) subscribe(desc turnpike.Subscribe, handle EventHandler) (func(), error) {
	log := chat.client.log.WithFields(logrus.Fields{
		"trace_id": uuid.New().String(),
		"type":     "SUBSCRIBE",
		"uri":      desc.Topic,
		"options":  desc.Topic,
	})

	log.Trace("exec subscribe")
	err := chat.ws.Subscribe(string(desc.Topic), desc.Options, func(args []interface{}, kwargs map[string]interface{}) {
		eType := 0
		if kwargs["type"] == nil {
			log.WithField("kwargs", kwargs).Warn("event kwargs missing type")
			eType = EVENT_UNKNOWN
		} else if eFloat, ok := kwargs["type"].(float64); ok {
			log.WithField("event_type", eType).Warn(fmt.Sprintf("event type was float %.4f", kwargs["type"]))
			eType = int(eFloat)
		} else {
			eType = kwargs["type"].(int)
		}

		log.WithField("event_type", eType).Trace("exec handle")
		if err := handle(eType, kwargs); err != nil {
			log.Error(err)
		}
	})

	return func() { chat.ws.Unsubscribe(string(desc.Topic)) }, err
}
