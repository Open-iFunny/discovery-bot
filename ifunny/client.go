package ifunny

import (
	"encoding/json"
	"io"
	"net/http"
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
