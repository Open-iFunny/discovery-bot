package ifunny

import (
	"github.com/gastrodon/popplio/ifunny/compose"
)

type ChatChannel struct {
	Name          string `json:"name"` // I think this is the unique id
	Title         string `json:"title"`
	MembersOnline int    `json:"members_online"`
	MembersTotal  int    `json:"members_total"`

	Type      int   `json:"type"`       // 1 = dm, ??
	JoinState int   `json:"join_state"` // 2 = true, ???
	Role      int   `json:"role"`
	TouchDT   int64 `json:"touch_dt"` // maybe when we last were online ??

	User struct {
		ID         string `json:"id"`
		Nick       string `json:"nick"`
		LastSeenAt int64  `json:"last_seen_at"`

		IsVerified bool `json:"is_verified"`
	} `json:"user"`
}

func (chat *Chat) handleChannelsRaw(handle func(channel *ChatChannel) error) EventHandler {
	return func(eventType int, kwargs map[string]interface{}) error {
		for _, channelRaw := range kwargs["chats"].([]interface{}) {
			if kwargs["chats"] == nil {
				chat.client.log.Warn("chats chunk is nil, skipping handler")
				return nil
			}

			channel := new(ChatChannel)
			if err := jsonDecode(channelRaw, channel); err != nil {
				return err
			}

			if err := handle(channel); err != nil {
				return err
			}
		}

		return nil
	}
}

func (chat *Chat) OnChannelJoin(handle func(channel *ChatChannel) error) (func(), error) {
	return chat.Subscribe(compose.JoinedChannels(chat.client.Self.ID), chat.handleChannelsRaw(handle))
}

func (chat *Chat) OnChannelInvite(handle func(channel *ChatChannel) error) (func(), error) {
	return chat.Subscribe(compose.PendingInvites(chat.client.Self.ID), chat.handleChannelsRaw(handle))
}
