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
		//UserToken: "71455ddd108b9cbc4e8b2d57b6f391a121de2f208ce82aa9c3b6d6d336c5e3cc9e83b773949ecaa4d8fa0",
	}
}
