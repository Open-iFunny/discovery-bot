package ifunny

import (
	"fmt"
	"io"
	"net/http"

	"github.com/jcelliott/turnpike"
	"github.com/mitchellh/mapstructure"
)

const (
	apiRoot   = "https://api.ifunny.mobi/v4"
	projectID = "iFunny"

	chatRoot      = "wss://chat.ifunny.co/chat"
	chatNamespace = "co.fun.chat"
)

type Client interface {
	request(method, path string, body io.Reader) (*http.Response, error)
	chat() (Chat, error)
}

func MakeClient(bearer, userAgent string) Client {
	return &staticClient{bearer, userAgent, http.DefaultClient}
}

type staticClient struct {
	bearer, userAgeng string
	http              *http.Client
}

func (client *staticClient) request(method, path string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, apiRoot+path, body)
	if err != nil {
		return nil, err
	}

	request.Header.Add("authorization", "bearer "+client.bearer)
	request.Header.Add("user-agent", client.userAgeng)
	request.Header.Add("ifunny-project-id", projectID)
	return client.http.Do(request)
}

func (client *staticClient) chat() (Chat, error) {
	ws, err := turnpike.NewWebsocketClient(turnpike.JSON, chatRoot, nil, nil, nil)
	if err != nil {
		panic(err)
	}

	ws.Auth = map[string]turnpike.AuthFunc{
		"ticket": turnpike.NewTicketAuthenticator(client.bearer),
	}
	hello, err := ws.JoinRealm(topic("ifunny"), nil)
	if err != nil {
		panic(err)
	}

	return &chat{ws, client, hello}, nil
}

type Chat interface {
	Chats(userID string) <-chan *WSChat
	Invites(userID string) <-chan *WSInvite
}

type chat struct {
	ws     *turnpike.Client
	client Client
	hello  map[string]interface{}
}

func topic(name string) string { return chatNamespace + "." + name }

func (chat *chat) Chats(userID string) <-chan *WSChat {
	result := make(chan *WSChat)
	chat.ws.Subscribe(topic("user."+userID+".chats"), nil, func(_ []interface{}, kwargs map[string]interface{}) {
		if kwargs["chats"] == nil {
			return
		}

		for _, chatRaw := range kwargs["chats"].([]interface{}) {
			wsChat := new(WSChat)
			mapstructure.Decode(chatRaw, wsChat)
			result <- wsChat
		}
	})

	return result
}

func (chat *chat) Invites(userID string) <-chan *WSInvite {
	result := make(chan *WSInvite)
	chat.ws.Subscribe(topic("user."+userID+".invites"), nil, func(_ []interface{}, kwargs map[string]interface{}) {
		if kwargs["invites"] == nil {
			return
		}

		for _, invRaw := range kwargs["invites"].([]interface{}) {
			fmt.Printf("invite: %+v\n", invRaw)

			wsInvite := new(WSInvite)
			mapstructure.Decode(invRaw, wsInvite)
			result <- wsInvite
		}
	})

	return result
}
