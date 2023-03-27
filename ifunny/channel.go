package ifunny

type APIChannel struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Type        int    `json:"type"`
	Description string `json:"description"`
	JoinState   int    `json:"join_state"`
}

type rChannel string

const (
	ChannelsTrending rChannel = "/chats/trending"
)

func (client *Client) Channels(path rChannel) ([]APIChannel, error) {
	response := new(struct {
		Data struct {
			Channels []APIChannel `json:"channels"`
		} `json:"data"`
	})

	err := client.apiRequest(response, "GET", string(path), nil)
	return response.Data.Channels, err
}