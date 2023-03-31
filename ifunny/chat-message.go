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
	ID   string `json:"id"`
	Text string `json:"text"`

	Type   messageType `json:"type"` // 1 = message, ???
	Status int         `json:"status"`
	PubAt  float64     `json:"pub_at"`

	User struct {
		ID         string `json:"user"`
		Nick       string `json:"nick"`
		IsVerified bool   `json:"is_verified"`
		LastSeenAt int64  `json:"last_seen_at"`
	} `json:"user"`
}
