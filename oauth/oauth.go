//Пакет для осуществления авторизации
package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/a1exCross/ElmaVK/vkerrors"
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
	ClientSecret string
	Redirect_URI string
	Group_IDs    []int
	GroupID      int
	Display      Display
	Scope        []Scope
	State        string
}

/* func GetGroupParams() AuthParams {
	return AuthParams{}
} */

//Генерация ссылки для получения токена пользователя
func ImplictFlow(p AuthParams) (*url.URL, error) {
	var u string

	if p.Client_ID != 0 {
		u += "client_id=" + fmt.Sprint(p.Client_ID)
	} else {
		return &url.URL{}, errors.New("Required field 'ClientID' is empty, MethodName - ImplictFlow()")
	}

	if p.Display != "" {
		u += fmt.Sprint("&display=") + string(p.Display)
	} else {
		u += "&display=page"
	}

	if p.Redirect_URI != "" {
		u += "&redirect_uri" + p.Redirect_URI
	}

	if p.Scope != nil {
		u += "&scope="
		for i, v := range p.Scope {
			if i > 0 {
				u += ","
			}
			u += string(v)
		}
	}

	if p.State != "" {
		u += "&state=" + p.State
	}

	u += "&response_type=token&v=" + version

	ur, err := url.Parse(auth_url + u)
	if err != nil {
		return &url.URL{}, err
	}

	return ur, nil
}

type Auth struct {
	AuthParams
	Client *http.Client
	URL    *url.URL
}

//Получение токена сообщества
func AuthCodeFlow(p AuthParams) Auth {
	var u string

	if p.Client_ID != 0 {
		u += "client_id=" + fmt.Sprint(p.Client_ID)
	}

	if p.GroupID != 0 {
		u += "&group_ids=" + fmt.Sprint(p.GroupID)
	}

	if p.Group_IDs != nil {
		if p.GroupID == 0 {
			u += "&group_ids="
		} else {
			u += ","
		}
		for i, v := range p.Group_IDs {
			if i > 0 {
				u += ","
			}
			u += fmt.Sprint(v)
		}
	}

	if p.Display != "" {
		u += "&display=" + string(p.Display)
	}

	if p.Redirect_URI != "" {
		u += "&redirect_uri=" + p.Redirect_URI
	}

	if p.Scope != nil {
		u += "&scope="

		for i, v := range p.Scope {
			if i > 0 {
				u += ","
			}
			u += string(v)
		}
	}

	if p.State != "" {
		u += "&state=" + p.State
	}

	u += "&response_type=code&v=" + version

	ur, _ := url.Parse(auth_url + u)

	return Auth{
		AuthParams: p,
		Client:     http.DefaultClient,
		URL:        ur,
	}
}

type GroupToken struct {
	GroupID     int    `json:"group_id"`
	AccessToken string `json:"access_token"`
}

type GroupTokens struct {
	Groups    []GroupToken `json:"groups"`
	ExpiresIn int          `json:"expires_in"`
}

func (a *Auth) get_token_request(code string) (GroupTokens, error) {
	var u string = ""

	if a.Client_ID != 0 {
		u += "client_id=" + fmt.Sprint(a.Client_ID)
	} else {
		return GroupTokens{}, errors.New("Required field 'ClientID' is empty, MethodName - get_token_request()")
	}

	if a.ClientSecret != "" {
		u += "&client_secret=" + a.ClientSecret
	} else {
		return GroupTokens{}, errors.New("Required field 'ClientSecret' is empty, MethodName - get_token_request()")
	}

	if code != "" {
		u += "&code=" + code
	} else {
		return GroupTokens{}, errors.New("Required field 'code' is empty, MethodName - get_token_request()")
	}

	u += "&redirect_uri=" + a.Redirect_URI

	req, err := http.NewRequest("GET", get_token_url+u, nil)
	if err != nil {
		return GroupTokens{}, err
	}

	res, err := a.Client.Do(req)
	if err != nil {
		return GroupTokens{}, err
	}

	check := vkerrors.GetError(res)
	if check != "ok" {
		return GroupTokens{}, errors.New(check)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return GroupTokens{}, err
	}

	var t GroupTokens

	err = json.Unmarshal(body, &t)
	if err != nil {
		return GroupTokens{}, err
	}

	return t, nil
}

func (a Auth) GetToken(u *url.URL) (GroupTokens, error) {
	code := u.Query().Get("code")

	tokens, err := a.get_token_request(code)
	if err != nil {
		return GroupTokens{}, err
	}

	return tokens, nil
}
