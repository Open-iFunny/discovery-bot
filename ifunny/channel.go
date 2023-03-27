package ifunny

const (
	CHATS_TRENDING = "/chats/trending"
)

type APIChannel struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Type        int    `json:"type"`
	Description string `json:"description"`
	JoinState   int    `json:"join_state"`
}

func (client *Client) Channels(path string) ([]APIChannel, error) {
	response := new(struct {
		Data struct {
			Channels []APIChannel `json:"channels"`
		} `json:"data"`
	})

	err := client.apiRequest(response, "GET", path, nil)
	return response.Data.Channels, err
}
