package ifunny

type APIContentPost struct {
	ID          string   `json:"id"`
	Link        string   `json:"link"`
	DateCreated int64    `json:"date_created"`
	PublushAt   int64    `json:"publish_at"`
	Tags        []string `json:"tags"`
	State       string   `json:"state"`
	ShotStatus  string   `json:"shot_status"`

	FastStart  bool `json:"fast_start"`
	IsFeatured bool `json:"is_featured"`
	IsPinned   bool `json:"is_pinned"`
	IsAbused   bool `json:"is_abused"`
	IsUnsafe   bool `json:"is_unsafe"`

	IsRepublished bool `json:"is_republished"`
	IsSmiled      bool `json:"is_smiled"`
	IsUnsmiled    bool `json:"is_unsmiled"`

	Size struct {
		Height int `json:"h"`
		Width  int `json:"w"`
	} `json:"size"`

	Num struct {
		Comments    int `json:"comments"`
		Republished int `json:"republished"`
		Smiles      int `json:"smiles"`
		Unsmiles    int `json:"unsmiles"`
		Views       int `json:"views"`
	} `json:"num"`

	Creator struct {
		ID   string `json:"id"`
		Nick string `json:"nick"`
	} `json:"creator"`
}

type APIPaging struct {
	Cursors struct {
		Next string `json:"next,omitempty"`
		Prev string `json:"prev,omitempty"` // TODO is this right
	} `json:"cursors"`
	HasNext bool `json:"hasNext"`
	HasPrev bool `json:"hasPrev"`
}

type APIFeedPage struct {
	Items  []APIContentPost `json:"items"`
	Paging APIPaging        `json:"paging"`
}

type rFeed string

const (
	FeedCollective rFeed = "/feeds/collective"
	FeedFeatures   rFeed = "/feeds/collective"
	TimelineHome   rFeed = "/timelines/home"
)

var (
	TimelineUserID   = func(id string) rFeed { return rFeed("/timelines/users/" + id) }
	TimelineUserNick = func(nick string) rFeed { return rFeed("/timelines/users/by_nick/" + nick) }
)

/*
Get a feed page
A timeline is considered a feed - Timeline values may be used here too
*/
func (client *Client) Feed(path rFeed) (APIFeedPage, error) {
	response := new(struct {
		Data struct {
			Content APIFeedPage `json:"content"`
		} `json:"data"`
	})

	if path == FeedCollective {
		err := client.apiRequest(response, "POST", string(path), nil)
		return response.Data.Content, err
	}

	err := client.apiRequest(response, "GET", string(path), nil)
	return response.Data.Content, err
}
