package ifunny

import (
	"github.com/mitchellh/mapstructure"
)

type WSChat struct {
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

type WSChats struct {
	Chats string `mapstructure:"chats"`
}

func (chat *Chat) Chats() <-chan *WSChat {
	result := make(chan *WSChat)
	uri := topic("user." + chat.client.self.ID + ".chats")
	chat.ws.Subscribe(uri, nil, func(_ []interface{}, kwargs map[string]interface{}) {
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

type cChannel call

func ChannelName(channel string) cChannel {
	return cChannel{
		procedure: topic("get_chat"),
		options:   map[string]interface{}{},
		args:      []interface{}{},
		kwargs:    map[string]interface{}{"chat_name": channel},
	}
}

func (client *Client) ChannelDM(them string) cChannel {
	return ChannelName(them + "_" + client.self.ID)
}

/*
Get a ws chat, and whether or not it exists
*/
func (chat *Chat) Channel(desc cChannel) (*WSChat, bool, error) {
	result, err := chat.ws.Call(desc.procedure, desc.options, desc.args, desc.kwargs)
	if err != nil {
		return nil, false, err
	}

	if result.ArgumentsKw["chat"] == nil {
		return nil, false, nil
	}

	wsChat := new(WSChat)
	err = mapstructure.Decode(result.ArgumentsKw["chat"], wsChat)
	return wsChat, true, err
}
