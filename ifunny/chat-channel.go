package ifunny

import (
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type ChatChannel struct {
	Name          string `mapstructure:"name"` // I think this is the unique id
	Title         string `mapstructure:"title"`
	MembersOnline int    `mapstructure:"members_online"`
	MembersTotal  int    `mapstructure:"members_total"`

	Type      int   `mapstructure:"type"`       // 1 = dm, ??
	JoinState int   `mapstructure:"join_state"` // 2 = true, ???
	Role      int   `mapstructure:"role"`
	TouchDT   int64 `mapstructure:"touch_dt"` // maybe when we last were online ??

	User struct {
		ID         string `mapstructure:"id"`
		Nick       string `mapstructure:"nick"`
		LastSeenAt int64  `mapstructure:"last_seen_at"`

		IsVerified bool `mapstructure:"is_verified"`
	} `mapstructure:"user"`
}

type cChannel call

func ChannelName(channel string) cChannel {
	return cChannel{
		procedure: uri("get_chat"),
		options:   map[string]interface{}{},
		args:      []interface{}{},
		kwargs:    map[string]interface{}{"chat_name": channel},
	}
}

func (client *Client) ChannelDM(them ...string) cChannel {
	us := append(them, client.self.ID)
	sort.Strings(us)
	size := len(us)
	backwards := make([]string, size)
	for index, each := range us {
		backwards[size-1-index] = each
	}

	return cChannel{
		procedure: uri("get_or_create_chat"),
		options:   nil,
		args:      nil,
		kwargs: map[string]interface{}{
			"type":  1,
			"users": them,
			"name":  strings.Join(backwards, "_"),
		},
	}
}

/*
Get a ws chat, and whether or not it exists
*/
func (chat *Chat) GetChannel(desc cChannel) (*ChatChannel, bool, error) {
	traceID := uuid.New().String()
	chat.client.log.WithFields(logrus.Fields{
		"trace_id":  traceID,
		"procedure": desc.procedure,
		"kwargs":    desc.kwargs,
	}).Trace("get channel")

	result, err := chat.ws.Call(desc.procedure, desc.options, desc.args, desc.kwargs)
	if err != nil {
		chat.client.log.WithField("trace_id", traceID).Error(err)
		return nil, false, err
	}

	if result.ArgumentsKw["chat"] == nil {
		chat.client.log.WithField("trace_id", traceID).Trace("no such channel")
		return nil, false, nil
	}

	wsChat := new(ChatChannel)
	err = mapstructure.Decode(result.ArgumentsKw["chat"], wsChat)
	return wsChat, true, err
}

type sChannel subscribe

func ChannelsIn(topic string) sChannel {
	return sChannel{
		topic:   uri(topic),
		options: nil,
	}
}

func (client *Client) ChannelsJoined() sChannel {
	return ChannelsIn("user." + client.self.ID + ".chats")
}

func (chat *Chat) IterChannel(desc sChannel) <-chan *ChatChannel {
	traceID := uuid.New().String()
	chat.client.log.WithFields(logrus.Fields{
		"trace_id": traceID,
		"topic":    desc.topic,
		"options":  desc.options,
	}).Trace("iter channel")

	result := make(chan *ChatChannel)
	chat.ws.Subscribe(desc.topic, desc.options, func(opts []interface{}, kwargs map[string]interface{}) {
		if kwargs["chats"] == nil {
			chat.client.log.WithField("trace_id", traceID).Trace("recv no channels")
			return
		}

		for _, channelRaw := range kwargs["chats"].([]interface{}) {
			channel := new(ChatChannel)
			if err := mapstructure.Decode(channelRaw, channel); err != nil {
				chat.client.log.WithField("trace_id", traceID).Error(err)
				continue
			}

			chat.client.log.WithFields(logrus.Fields{
				"trace_id":      traceID,
				"channel_name":  channel.Name,
				"channel_title": channel.Title,
			})
			result <- channel
		}
	})

	return result
}
