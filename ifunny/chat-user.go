package ifunny

type ChatUser struct {
	ID         string `mapstructure:"id"`
	Nick       string `mapstructure:"nick"`
	IsVerified bool   `mapstructure:"is_verified"`
	LastSeenAt int64  `mapstructure:"last_seen_at"`

	Photo string `mapstructure:"photo"`
}
