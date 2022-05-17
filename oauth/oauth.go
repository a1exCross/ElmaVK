package oauth

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const auth_url = "https://oauth.vk.com/authorize?"
const token_geturl = "https://oauth.vk.com/access_token?"
const version = "5.131"

type Display string

const (
	Page   Display = "page"
	Popup  Display = "popup"
	Mobile Display = "mobile"
)

type Scope string

const (
	Messages Scope = "messages"
	Manage   Scope = "manage"
	Photos   Scope = "photos"
	Docs     Scope = "docs"
)

type Group_params struct {
	Client_ID    int
	Redirect_URI string
	Group_IDs    []int
	Display      Display
	Scope        []Scope
	State        string
}

func GetGroupParams() *Group_params {
	return &Group_params{}
}

type Auth struct {
	Group_params
	client_secret string
	Client        *http.Client
}

func CodeAuthInUrl(g Group_params, client_secret string) (Auth, string) {
	var u string

	if g.Client_ID != 0 {
		u += "client_id=" + strconv.Itoa(g.Client_ID)
	}

	if g.Group_IDs != nil {
		u += "&group_ids="
	}

	for i, v := range g.Group_IDs {
		if i > 0 {
			u += ","
		}
		u += strconv.Itoa(v)
	}

	if g.Display != "" {
		u += "&display=" + string(g.Display)
	}

	if g.Redirect_URI != "" {
		u += "&redirect_uri=" + g.Redirect_URI
	}

	if g.Scope != nil {
		u += "&scope="
	}

	for i, v := range g.Scope {
		if i > 0 {
			u += ","
		}
		u += string(v)
	}

	u += "&response_type=code&v=" + version

	return Auth{
		Group_params:  g,
		client_secret: client_secret,
		Client:        http.DefaultClient,
	}, auth_url + u
}

type GroupToken struct {
	GroupID     int    `json:"group_id"`
	AccessToken string `json:"access_token"`
}

type GroupTokens struct {
	Groups    []GroupToken `json:"groups"`
	ExpiresIn int          `json:"expires_in"`
}

func (a *Auth) req_token(code string) *GroupTokens {
	u := "client_id=" + strconv.Itoa(a.Client_ID) +
		"&client_secret=" + a.client_secret +
		"&redirect_uri=" + a.Redirect_URI +
		"&code=" + code

	req, _ := http.NewRequest("GET", token_geturl+u, nil)

	res, err := a.Client.Do(req)

	if err != nil {
		log.Println(err)
	}

	data, err := ioutil.ReadAll(res.Body)

	var t GroupTokens
	err = json.Unmarshal(data, &t)

	return &t
}

func (a Auth) GetToken(u *url.URL) *GroupTokens {
	code := u.Query().Get("code")

	return a.req_token(code)
}
