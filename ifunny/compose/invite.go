package compose

import "github.com/jcelliott/turnpike"

func inviteCall(accept bool) string {
	if accept {
		return "invite.accept"
	}

	return "invite.decline"
}

func Invite(channel string, accept bool) turnpike.Call {
	return turnpike.Call{
		Procedure:   URI(inviteCall(accept)),
		ArgumentsKw: map[string]interface{}{"chat_name": channel},
	}
}
