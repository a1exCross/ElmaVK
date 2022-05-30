package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/a1exCross/ElmaVK/ApiErrors"
)

const auth_url = "https://oauth.vk.com/authorize?"
const get_token_url = "https://oauth.vk.com/access_token?"
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
	Video    Scope = "video"
)

type AuthParams struct {
	Client_ID    int
	CleintSecret string
	Redirect_URI string
	//Group_IDs    []int
	GroupID int
	Display Display
	Scope   []Scope
	State   string
	Revoke  bool
}

func GetGroupParams() AuthParams {
	return AuthParams{}
}

type Auth struct {
	AuthParams
	client_secret string
	Client        *http.Client
}

//Генерация ссылки для получения токена пользователя
func ImplictFlow(p AuthParams) (*url.URL, error) {
	var u string

	if p.Client_ID != 0 {
		u += fmt.Sprintf("client_id=%d", p.Client_ID)
	} else {
		return &url.URL{}, errors.New("ClientID не указан")
	}

	if p.Display != "" {
		u += fmt.Sprintf("display=%s", p.Display)
	} else {
		u += "display=page"
	}

	if p.Redirect_URI != "" {
		u += fmt.Sprintf("redirect_uri=%s", p.Redirect_URI)
	}

	if p.Scope != nil {
		u += "&scope="
	}

	for i, v := range p.Scope {
		if i > 0 {
			u += ","
		}
		u += string(v)
	}

	u += "&response_type=token&v=" + version

	u = auth_url + u

	ur, err := url.Parse(u)

	if err != nil {
		return &url.URL{}, err
	}

	return ur, nil
}

//Получение токена сообщества
func AuthCodeFlow(g AuthParams) (Auth, string) {
	var u string

	if g.Client_ID != 0 {
		u += "client_id=" + strconv.Itoa(g.Client_ID)
	}

	if g.GroupID != 0 {
		u += "&group_ids="
	}

	/* 	if g.Group_IDs != nil {
	   		u += "&group_ids="
	   	}

	   	for i, v := range g.Group_IDs {
	   		if i > 0 {
	   			u += ","
	   		}
	   		u += strconv.Itoa(v)
	   	} */

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
		AuthParams: g,
		Client:     http.DefaultClient,
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

func (a *Auth) req_token(code string) (GroupTokens, error) {
	u := "client_id=" + strconv.Itoa(a.Client_ID) +
		"&client_secret=" + a.client_secret +
		"&redirect_uri=" + a.Redirect_URI +
		"&code=" + code

	req, err := http.NewRequest("GET", get_token_url+u, nil)

	if err != nil {
		return GroupTokens{}, err
	}

	res, err := a.Client.Do(req)

	if err != nil {
		return GroupTokens{}, err
	}

	check := ApiErrors.GetError(res)

	if check != "ok" {
		return GroupTokens{}, errors.New(check)
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return GroupTokens{}, err
	}

	var t GroupTokens
	err = json.Unmarshal(data, &t)

	return t, nil
}

func (a Auth) GetToken(u *url.URL) (GroupTokens, error) {
	code := u.Query().Get("code")

	tokens, err := a.req_token(code)

	if err != nil {
		return GroupTokens{}, err
	}

	return tokens, nil
}
