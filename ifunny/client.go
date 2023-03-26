package ifunny

import (
	"io"
	"net/http"
)

const (
	apiRoot   = "https://api.ifunny.mobi/v4"
	projectID = "iFunny"
)

type Client interface {
	Request(method, path string, body io.Reader) (*http.Response, error)
	Connect(bearer, cookie string) (Chat, error)
}

func MakeClient(authorization, userAgent string) Client {
	return &staticClient{authorization, userAgent, http.DefaultClient}
}

type staticClient struct {
	authorization, userAgeng string
	http                     *http.Client
}

func (client *staticClient) Request(method, path string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, apiRoot+path, body)
	if err != nil {
		return nil, err
	}

	request.Header.Add("authorization", client.authorization)
	request.Header.Add("user-agent", client.userAgeng)
	request.Header.Add("ifunny-project-id", projectID)
	return client.http.Do(request)
}

func (client *staticClient) Connect(bearer, cookie string) (Chat, error) {
	return connectChat(bearer, cookie)
}
