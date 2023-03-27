package ifunny

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

type WSInvite struct {
}

func (chat *Chat) Invites() <-chan *WSInvite {
	result := make(chan *WSInvite)
	uri := uri("user." + chat.client.self.ID + ".invites")
	chat.ws.Subscribe(uri, nil, func(_ []interface{}, kwargs map[string]interface{}) {
		if kwargs["invites"] == nil {
			return
		}

		for _, invRaw := range kwargs["invites"].([]interface{}) {
			fmt.Printf("invite: %+v\n", invRaw)

			wsInvite := new(WSInvite)
			mapstructure.Decode(invRaw, wsInvite)
			result <- wsInvite
		}
	})

	return result
}
