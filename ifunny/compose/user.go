package compose

import "github.com/jcelliott/turnpike"

/*
call out to get chat contacts
*/
func Contacts(limit int) turnpike.Call {
	return turnpike.Call{
		Procedure:   URI("list_contacts"),
		ArgumentsKw: map[string]interface{}{"limit": limit},
	}
}

func UserByID(id string) Request {
	return get("/users/"+id, nil)
}

func UserByNick(nick string) Request {
	return get("/users/by_nick/"+nick, nil)
}

func UserAccount() Request {
	return get("/account", nil)
}
