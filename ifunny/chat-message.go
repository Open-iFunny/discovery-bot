package ifunny

import (
	"github.com/mitchellh/mapstructure"
)

type ChatMessage struct {
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

type ChatMessages struct {
	Chats []ChatMessage `mapstructure:"chats"`
}

type sMessages subscribe

func (client *Client) MessageUnread() sMessages {
	return sMessages{
		topic:   uri("user." + client.self.ID + ".chats"),
		options: map[string]interface{}{},
	}
}

func (chat *Chat) IterMessage(desc sMessages) <-chan *ChatMessage {
	result := make(chan *ChatMessage)
	chat.ws.Subscribe(desc.topic, desc.options, func(opts []interface{}, kwargs map[string]interface{}) {
		if kwargs["chats"] == nil {
			return
		}

		for _, messageRaw := range kwargs["chats"].([]interface{}) {
			message := new(ChatMessage)
			if err := mapstructure.Decode(messageRaw, message); err != nil {
				panic(err)
			}

			result <- message
		}
	})

	return result
}
