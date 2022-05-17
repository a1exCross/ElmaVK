package Methods

import (
	"net/http"
)

type VK struct {
	Token     string
	UserToken string
	Client    *http.Client
	Version   string
}

func Session(token string) *VK {
	return &VK{
		Token:     token,
		Client:    http.DefaultClient,
		Version:   "5.131",
		UserToken: "b0624602220b43dbc2f151db2af83a75e8e740f18d1cf0d36b26f63027498420569182933181cb590b5b1",
	}
}
