package compose

import (
	"sort"
	"strings"

	"github.com/jcelliott/turnpike"
)

type channelType int

const (
	ChannelDM      channelType = 1
	ChannelPrivate channelType = 2
	ChannelPublic  channelType = 3
)

func PendingInvites(id string) turnpike.Subscribe {
	return turnpike.Subscribe{Topic: URI("user." + id + ".invites")}
}

func JoinedChannels(id string) turnpike.Subscribe {
	return turnpike.Subscribe{Topic: URI("user." + id + ".chats")}
}

func HideChannel(channel string) turnpike.Call {
	return turnpike.Call{Procedure: URI("hide_chat")}
}

func DMChannelName(self string, them []string) string {
	us := append(them, self)
	sort.Strings(us)
	size := len(us)
	backwards := make([]string, size)
	for index, each := range us {
		backwards[size-1-index] = each
	}

	return strings.Join(backwards, "_")
}

func GetDMChannel(id string, them ...string) turnpike.Call {
	return turnpike.Call{
		Procedure: URI("get_or_create_chat"),
		ArgumentsKw: map[string]interface{}{
			"type":  ChannelDM,
			"users": them,
			"name":  DMChannelName(id, them),
		},
	}
}

func NewChannel(title, name, description string, invite []string, channelType channelType) turnpike.Call {
	if description != "" && channelType == ChannelPrivate {
		panic("cannot add a description to a private channel")
	}

	return turnpike.Call{
		Procedure: URI("new_chat"),
		ArgumentsKw: map[string]interface{}{
			"users":       invite,
			"title":       title,
			"name":        name,
			"description": description,
			"type":        channelType,
		},
	}
}

func GetChannel(channel string) turnpike.Call {
	return turnpike.Call{
		Procedure:   URI("get_chat"),
		ArgumentsKw: map[string]interface{}{"chat_name": channel},
	}
}

var (
	ChatsTrending = Request{Method: "GET", Path: "/chats/trending"}
)
