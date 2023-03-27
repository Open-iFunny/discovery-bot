package ifunny

type sMessages subscribe

func MessageIn(channel string) sMessages {
	return sMessages{
		topic:   uri(channel),
		options: nil,
	}
}
