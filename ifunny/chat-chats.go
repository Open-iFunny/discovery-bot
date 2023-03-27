package ifunny

import "github.com/mitchellh/mapstructure"

type WSChat struct {
	Name          string `mapstructure:"name"`
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

/*
Get a ws chat, and whether or not it exists
*/
func (chat *Chat) GetDM(id string) (*WSChat, bool, error) {
	kwargs := map[string]interface{}{"chat_name": id + "_" + chat.client.self.ID}
	result, err := chat.ws.Call(topic("get_chat"), nil, nil, kwargs)
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
