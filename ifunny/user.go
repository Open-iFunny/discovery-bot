package ifunny

const (
	USER_SELF = "/account"
)

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

func (client *Client) User(path string) (APIUser, error) {
	response := new(struct {
		Data APIUser `json:"data"`
	})

	err := client.apiRequest(response, "GET", path, nil)
	return response.Data, err
}
