package ifunny

type EventHandler func(eventType int, kwargs map[string]interface{}) error

const (
	EVENT_UNKNOWN = -1
	EVENT_JOIN    = 100
	EVENT_INVITED = 300
)
