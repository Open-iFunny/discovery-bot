package ifunny

import "github.com/mitchellh/mapstructure"

type WSChat struct {
	Name  string `mapstructure:"name"`
	Title string `mapstructure:"title"`
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
