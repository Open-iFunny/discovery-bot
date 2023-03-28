package ifunny

import (
	"github.com/google/uuid"
	"github.com/jcelliott/turnpike"
)

type call struct {
	procedure string
	options   map[string]interface{}
	args      []interface{}
	kwargs    map[string]interface{}
}

type subscribe struct {
	topic   string
	options map[string]interface{}
}

type publish struct {
	topic   string
	options map[string]interface{}
	args    []interface{}
	kwargs  map[string]interface{}
}

const (
	chatRoot      = "wss://chat.ifunny.co/chat"
	chatNamespace = "co.fun.chat"
)

func uri(name string) string { return chatNamespace + "." + name }

func (client *Client) Chat() (*Chat, error) {
	traceID := uuid.New().String()
	client.log.WithField("trace_id", traceID).Trace("start websocket")

	ws, err := turnpike.NewWebsocketClient(turnpike.JSON, chatRoot, nil, nil, nil)
	if err != nil {
		client.log.WithField("trace_id", traceID).Error(err)
		return nil, err
	}

	client.log.WithField("trace_id", traceID).Trace("join realm ifunny")
	ws.Auth = map[string]turnpike.AuthFunc{"ticket": turnpike.NewTicketAuthenticator(client.bearer)}
	hello, err := ws.JoinRealm(uri("ifunny"), nil)
	if err != nil {
		client.log.WithField("trace_id", err).Error(err)
		return nil, err
	}

	return &Chat{ws, client, hello}, nil
}

type Chat struct {
	ws     *turnpike.Client
	client *Client
	hello  map[string]interface{}
}
