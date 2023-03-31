package ifunny

type ChatChannel struct {
	Name          string `json:"name"` // I think this is the unique id
	Title         string `json:"title"`
	MembersOnline int    `json:"members_online"`
	MembersTotal  int    `json:"members_total"`

	Type      int   `json:"type"`       // 1 = dm, ??
	JoinState int   `json:"join_state"` // 2 = true, ???
	Role      int   `json:"role"`
	TouchDT   int64 `json:"touch_dt"` // maybe when we last were online ??

	User struct {
		ID         string `json:"id"`
		Nick       string `json:"nick"`
		LastSeenAt int64  `json:"last_seen_at"`

		IsVerified bool `json:"is_verified"`
	} `json:"user"`
}
