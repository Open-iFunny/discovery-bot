package ifunny

import (
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type ChatUser struct {
	ID         string `mapstructure:"id"`
	Nick       string `mapstructure:"nick"`
	IsVerified bool   `mapstructure:"is_verified"`
	LastSeenAt int64  `mapstructure:"last_seen_at"`

	Photo string `mapstructure:"photo"`
}

type cUser call

func Contacts(limit int) cUser {
	return cUser{
		procedure: uri("list_contacts"),
		options:   nil,
		args:      nil,
		kwargs:    map[string]interface{}{"limit": limit},
	}
}

func (chat *Chat) GetUsers(desc cUser) ([]*ChatUser, error) {
	chats := new(struct {
		Users []*ChatUser `mapstructure:"users"`
	})

	err := chat.call(call(desc), &chats)
	return chats.Users, err
}

func (chat *Chat) _GetUsers(desc cUser) ([]*ChatUser, error) {
	traceID := uuid.New().String()
	chat.client.log.WithFields(logrus.Fields{
		"trace_id":  traceID,
		"procedure": desc.procedure,
		"kwargs":    desc.kwargs,
	}).Trace("begin get users")

	result, err := chat.ws.Call(desc.procedure, desc.options, desc.args, desc.kwargs)
	if err != nil {
		chat.client.log.WithField("trace_id", traceID).Trace(err)
	}

	if result.ArgumentsKw["users"] == nil {
		chat.client.log.WithField("trace_id", traceID).Trace("no users returned")
		return nil, nil
	}

	usersRaw := result.ArgumentsKw["users"].([]interface{})
	users := make([]*ChatUser, len(usersRaw))
	for index, eachRaw := range usersRaw {
		user := new(ChatUser)
		if err := mapstructure.Decode(eachRaw, user); err != nil {
			chat.client.log.WithField("trace_id", traceID).Trace(err)
			return nil, err
		}

		users[index] = user
	}

	return users, nil
}
