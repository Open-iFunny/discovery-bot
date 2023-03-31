package ifunny

type ChatChannel struct {
	Name          string `mapstructure:"name"` // I think this is the unique id
	Title         string `mapstructure:"title"`
	MembersOnline int    `mapstructure:"members_online"`
	MembersTotal  int    `mapstructure:"members_total"`

	Type      int   `mapstructure:"type"`       // 1 = dm, ??
	JoinState int   `mapstructure:"join_state"` // 2 = true, ???
	Role      int   `mapstructure:"role"`
	TouchDT   int64 `mapstructure:"touch_dt"` // maybe when we last were online ??

	User struct {
		ID         string `mapstructure:"id"`
		Nick       string `mapstructure:"nick"`
		LastSeenAt int64  `mapstructure:"last_seen_at"`

		IsVerified bool `mapstructure:"is_verified"`
	} `mapstructure:"user"`
}
