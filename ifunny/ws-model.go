package ifunny

type WSChat struct {
	Name  string `mapstructure:"name"`
	Title string `mapstructure:"title"`
}

type WSChats struct {
	Chats string `mapstructure:"chats"`
}
