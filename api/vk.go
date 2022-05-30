package api

import (
	"net/http"
)

type VK struct {
	Token     string
	UserToken string
	Client    *http.Client
	Version   string
}

func Session(token string) VK {
	return VK{
		Token:   token,
		Client:  http.DefaultClient,
		Version: "5.131",
	}
}
