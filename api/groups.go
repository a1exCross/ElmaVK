package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"

	"github.com/a1exCross/ElmaVK/ApiErrors"
)

type group_info struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
	IsClosed   int    `json:"is_closed"`
	Type       string `json:"type"`
	Photo50    string `json:"photo_50"`
	Photo100   string `json:"photo_100"`
	Photo200   string `json:"photo_200"`
}

type GroupInfo struct {
	Response []group_info `json:"response"`
}

//https://dev.vk.com/method/groups.get
func (v VK) GetCurrentGroup() (*GroupInfo, error) {
	var u string = ""

	if v.Token != "" {
		u = "access_token=" + v.Token +
			"&v=5.131"
	} else {
		return nil, errors.New("Auth token is empty")
	}

	res, err := v.Reqeust_api_get("groups.getById?", u)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	check := ApiErrors.GetError(res)

	if check != "ok" {
		return nil, errors.New(check)
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	group := GroupInfo{}

	err = json.Unmarshal(data, &group)

	if err != nil {
		return nil, err
	}

	return &group, nil
}

type ConfirmationKey struct {
	Response struct {
		Code string `json:"code"`
	} `json:"response"`
}

//https://dev.vk.com/method/groups.getCallbackConfirmationCode
func (v VK) GetConfirmaionKey(group_id int) (string, error) {
	var u string = ""

	if group_id != 0 {
		u += "group_id=" + fmt.Sprint(group_id)
	} else {
		return "", errors.New("Required field 'GroupID' is empty, MethodName - GetConfirmationKey()")
	}

	if v.Token != "" {
		u = "&access_token=" + v.Token + "&v=" + v.Version
	} else {
		return "", errors.New("Auth token is empty")
	}

	res, err := v.Reqeust_api_get("groups.getCallbackConfirmationCode?", u)

	if err != nil {
		return "", err
	}

	check := ApiErrors.GetError(res)

	if check != "ok" {
		return "", errors.New(check)
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	key := ConfirmationKey{}

	err = json.Unmarshal(data, &key)

	if err != nil {
		return "", err
	}

	return key.Response.Code, nil
}

type ServerItem struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	CreatorID int    `json:"creator_id"`
	URL       string `json:"url"`
	SecretKey string `json:"secret_key"`
	Status    string `json:"status"`
}

type CallbackServersResponse struct {
	Response struct {
		Count int          `json:"count"`
		Items []ServerItem `json:"items"`
	} `json:"response"`
}

//https://dev.vk.com/method/groups.getCallbackServers
func (v VK) GetCallbackServers(group_id int) ([]ServerItem, error) { //groups.getCallbackServers
	var u string = ""

	if group_id != 0 {
		u += "group_id=" + fmt.Sprint(group_id)
	} else {
		return nil, errors.New("Required field 'GroupID' is empty, MethodName - GetCallbackServers()")
	}

	if v.Token != "" {
		u = "&access_token=" + v.Token + "&v=" + v.Version
	} else {
		return nil, errors.New("Auth token is empty")
	}

	res, err := v.Reqeust_api_get("groups.getCallbackServers?", u)

	if err != nil {
		return nil, err
	}

	check := ApiErrors.GetError(res)

	if check != "ok" {
		return nil, errors.New(check)
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	resp := CallbackServersResponse{}

	err = json.Unmarshal(data, &resp)

	if err != nil {
		return nil, err
	}

	return resp.Response.Items, nil
}

type AddCallbackServerResponse struct {
	Response struct {
		ServerID int `json:"server_id"`
	} `json:"response"`
}

//https://dev.vk.com/method/groups.addCallbackServer
func (v VK) AddCallbackServer(u string) (int, error) {
	res, err := v.Reqeust_api_get("groups.addCallbackServer?", u)

	if err != nil {
		return 0, err
	}

	check := ApiErrors.GetError(res)

	if check != "ok" {
		return 0, errors.New(check)
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return 0, err
	}

	resp := AddCallbackServerResponse{}

	err = json.Unmarshal(data, &resp)

	if err != nil {
		return 0, err
	}

	return resp.Response.ServerID, nil
}

type ResponseServerDeleteOrSet struct {
	Response int `json:"response"`
}

//https://dev.vk.com/method/groups.deleteCallbackServer
func (v VK) DeleteCallbackServer(group_id, serv_id int) (int, error) {
	var u string = ""

	if group_id != 0 {
		u += "group_id=" + fmt.Sprint(-group_id)
	} else {
		return 0, errors.New("Required field 'GroupID' is empty, MethodName - DeleteCallbackServer()")
	}

	if serv_id != 0 {
		u += "&server_id=" + fmt.Sprint(serv_id)
	} else {
		return 0, errors.New("Required field 'ServerID' is empty, MethodName - DeleteCallbackServer()")
	}

	if v.Token != "" {
		u = "&access_token=" + v.Token + "&v=" + v.Version
	}

	res, err := v.Reqeust_api_get("groups.deleteCallbackServer?", u)

	if err != nil {
		return -1, err
	}

	check := ApiErrors.GetError(res)

	if check != "ok" {
		return 0, errors.New(check)
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return -1, err
	}

	resp := ResponseServerDeleteOrSet{}

	err = json.Unmarshal(data, &resp)

	if err != nil {
		return -1, err
	}

	return resp.Response, nil
}

type CallbackSettings struct {
	GroupID  int `json:"group_id"`
	ServerID int `json:"server_id"`
	CallbackEventSettings
}

type CallbackEventSettings struct {
	MessageNew         bool `json:"message_new"`
	MessageReply       bool `json:"message_reply"`
	MessageEdit        bool `json:"message_edit"`
	MessageAllow       bool `json:"message_allow"`
	MessageDeny        bool `json:"message_deny"`
	MessageTypingState bool `json:"message_typing_state"`
	MessageEvent       bool `json:"message_event"`
}

func (v VK) GetCallbackSettingsParams() CallbackSettings {
	return CallbackSettings{}
}

//https://dev.vk.com/method/groups.setCallbackSettings
func (v VK) SetCallbackSettings(m CallbackSettings) (int, error) {
	var u string = ""

	if m.GroupID != 0 {
		u += "group_id=" + fmt.Sprint(int(math.Abs(float64(m.GroupID))))
	} else {
		return 0, errors.New("Required field 'GroupID' is empty, MethodName - SetCallbackSettings()")
	}

	if m.ServerID != 0 {
		u += "&server_id=" + fmt.Sprint(m.ServerID)
	}

	u += "&message_new=" + fmt.Sprint(m.MessageNew)

	u += "&message_reply=" + fmt.Sprint(m.MessageReply)

	u += "&message_edit=" + fmt.Sprint(m.MessageEdit)

	u += "&message_allow=" + fmt.Sprint(m.MessageAllow)

	u += "&message_deny=" + fmt.Sprint(m.MessageDeny)

	u += "&message_event=" + fmt.Sprint(m.MessageEvent)

	u += "&message_typing_state=" + fmt.Sprint(m.MessageTypingState)

	if v.Token != "" {
		u += "&access_token=" + v.Token + "&v=" + v.Version
	} else {
		return 0, errors.New("Auth token is empty")
	}

	res, err := v.Reqeust_api_get("groups.setCallbackSettings?", u)

	if err != nil {
		return -1, err
	}

	check := ApiErrors.GetError(res)

	if check != "ok" {
		return 0, errors.New(check)
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return -1, err
	}

	resp := ResponseServerDeleteOrSet{}

	err = json.Unmarshal(data, &resp)

	if err != nil {
		return -1, err
	}

	return resp.Response, nil
}
