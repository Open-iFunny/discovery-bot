package bot

import (
	"strings"

	"github.com/gastrodon/popplio/ifunny"
)

type prefix struct {
	prefix string
}

func (us filter) And(also filter) filter {
	return func(event *ifunny.ChatEvent) bool {
		return us(event) && also(event)
	}
}

func (us filter) Not(also filter) filter {
	return func(event *ifunny.ChatEvent) bool {
		return us(event) && !also(event)
	}
}

func Prefix(fix string) prefix { return prefix{fix} }

func (fix prefix) Cmd(name string) filter {
	return func(event *ifunny.ChatEvent) bool {
		if event.Text != "" {
			return event.Text == fix.prefix+name || strings.HasPrefix(event.Text, fix.prefix+name+" ")
		}

		return false
	}
}

func AuthoredBy(nick string) filter {
	return func(event *ifunny.ChatEvent) bool {
		return strings.ToLower(event.User.Nick) == nick
	}
}
