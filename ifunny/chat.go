package ifunny

import (
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

func uri(name string) string { return chatNamespace + "." + name }

func (client *Client) Chat() (*Chat, error) {
	ws, err := turnpike.NewWebsocketClient(turnpike.JSON, chatRoot, nil, nil, nil)
	if err != nil {
		panic(err)
	}

	ws.Auth = map[string]turnpike.AuthFunc{"ticket": turnpike.NewTicketAuthenticator(client.bearer)}
	hello, err := ws.JoinRealm(uri("ifunny"), nil)
	if err != nil {
		panic(err)
	}

	return &Chat{ws, client, hello}, nil
}

type Chat struct {
	ws     *turnpike.Client
	client *Client
	hello  map[string]interface{}
}
