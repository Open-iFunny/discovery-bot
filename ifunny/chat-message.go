package ifunny

const (
	UNK_0 messageType = iota
	MESSAGE
	UNK_2
	JOIN_CHANNEL
	EXIT_CHANNEL
)

type messageType int

type ChatMessage struct {
	ID   string `mapstructure:"id"`
	Text string `mapstructure:"text"`

	Type   messageType `mapstructure:"type"` // 1 = message, ???
	Status int         `mapstructure:"status"`
	PubAt  float64     `mapstructure:"pub_at"`

	User struct {
		ID         string `mapstructure:"user"`
		Nick       string `mapstructure:"nick"`
		IsVerified bool   `mapstructure:"is_verified"`
		LastSeenAt int64  `mapstructure:"last_seen_at"`
	} `mapstructure:"user"`
}
