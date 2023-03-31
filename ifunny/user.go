package ifunny

import "github.com/gastrodon/popplio/ifunny/compose"

type APIUser struct {
	Email            string `json:"email"`
	SafeMode         bool   `json:"safe_mode"`
	OriginalNick     string `json:"original_nick"`
	MessagingPrivacy string `json:"messaging_privacy_status"`

	ID    string `json:"id"`
	Nick  string `json:"nick"`
	About string `json:"about"`

	IsAvailableForChat bool `json:"is_available_for_chat"`
	IsBanned           bool `json:"is_banned"`
	IsDeleted          bool `json:"is_deleted"`
	IsModerator        bool `json:"is_moderator"`
	IsVerified         bool `json:"is_verified"`
}

func (client *Client) GetUser(desc compose.Request) (*APIUser, error) {
	user := new(struct {
		Data APIUser `json:"data"`
	})

	err := client.RequestJSON(desc, user)
	return &user.Data, err
}
