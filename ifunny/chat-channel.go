package ifunny

import (
	"github.com/mitchellh/mapstructure"
)

type cChannel call

func ChannelName(channel string) cChannel {
	return cChannel{
		procedure: uri("get_chat"),
		options:   map[string]interface{}{},
		args:      []interface{}{},
		kwargs:    map[string]interface{}{"chat_name": channel},
	}
}

/*
Get a ws chat, and whether or not it exists
*/
func (chat *Chat) Channel(desc cChannel) (*ChatMessage, bool, error) {
	result, err := chat.ws.Call(desc.procedure, desc.options, desc.args, desc.kwargs)
	if err != nil {
		return nil, false, err
	}

	if result.ArgumentsKw["chat"] == nil {
		return nil, false, nil
	}

	wsChat := new(ChatMessage)
	err = mapstructure.Decode(result.ArgumentsKw["chat"], wsChat)
	return wsChat, true, err
}
