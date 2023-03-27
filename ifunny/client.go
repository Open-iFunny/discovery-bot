package ifunny

import (
	"encoding/json"
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

func MakeClient(bearer, userAgent string) (*Client, error) {
	client := &Client{bearer, userAgent, http.DefaultClient, nil}
	self, err := client.User(UserAccount)
	if err != nil {
		return nil, err
	}

	client.self = &self
	return client, nil
}

type Client struct {
	bearer, userAgeng string
	http              *http.Client
	self              *APIUser
}

func request(method, path string, body io.Reader, header http.Header, client *http.Client) (*http.Response, error) {
	request, err := http.NewRequest(method, apiRoot+path, body)
	if err != nil {
		return nil, err
	}
	request.Header = header
	return client.Do(request)
}

func (client *Client) header() http.Header {
	return http.Header{
		"authorization":     []string{"bearer " + client.bearer},
		"user-agent":        []string{client.userAgeng},
		"ifunny-project-id": []string{projectID},
	}
}

func (client *Client) apiRequest(apiBody interface{}, method, path string, body io.Reader) error {
	response, err := request(method, path, body, client.header(), client.http)
	if err != nil {
		return err
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, apiBody)
	return err
}

func (client *Client) chat() (*Chat, error) {
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

	return &Chat{ws, client, hello}, nil
}

type Chat struct {
	ws     *turnpike.Client
	client *Client
	hello  map[string]interface{}
}

func topic(name string) string { return chatNamespace + "." + name }

func (chat *Chat) Chats(userID string) <-chan *WSChat {
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

func (chat *Chat) Invites(userID string) <-chan *WSInvite {
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
