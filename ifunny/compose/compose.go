package compose

import (
	"io"
	"net/url"

	"github.com/jcelliott/turnpike"
)

const chatNamespace = "co.fun.chat"

func URI(name string) turnpike.URI { return turnpike.URI(chatNamespace + "." + name) }

type Request struct {
	Method, Path string
	Body         io.Reader
	Query        url.Values
}

func get(path string, query url.Values) Request {
	return Request{Method: "GET", Path: path, Query: query}
}
