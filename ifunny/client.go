package ifunny

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	apiRoot   = "https://api.ifunny.mobi/v4"
	projectID = "iFunny"

	logLevel = logrus.InfoLevel
)

func MakeClient(bearer, userAgent string) (*Client, error) {
	client := &Client{bearer, userAgent, http.DefaultClient, nil, logrus.New()}
	client.log.SetFormatter(&logrus.JSONFormatter{})
	client.log.SetLevel(logLevel)

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
	log               *logrus.Logger
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
	traceID := uuid.New().String()
	client.log.WithFields(logrus.Fields{
		"trace_id": traceID,
		"path":     path,
		"method":   method,
		"has_body": body != nil},
	).Trace("make request")

	response, err := request(method, path, body, client.header(), client.http)
	if err != nil {
		client.log.WithField("trace_id", traceID).Error(err)
		return err
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		client.log.WithField("trace_id", traceID).Error(err)
		return err
	}

	err = json.Unmarshal(bodyBytes, apiBody)
	if err != nil {
		client.log.WithField("trace_id", traceID).Error(err)
	}

	return err
}
