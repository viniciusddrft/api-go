package entities

type Tweet struct {
	Id          int64  `json:"id"`
	Description string `json:"description"`
}

func NewTweet() *Tweet {
	tweet := Tweet{}

	return &tweet
}
