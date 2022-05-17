package Methods

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"strconv"

	ApiErrors "github.com/a1exCross/ElmaVK/errors"
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

func (v *VK) GetCurrentGroup() (*GroupInfo, error) {
	u := "access_token=" + v.Token +
		"&v=5.131"
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

func (v *VK) GetConfirmaionKey(group_id int) (string, error) {
	u := "access_token=" + v.Token +
		"&group_id=" + strconv.Itoa(group_id) +
		"&v=" + v.Version

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

func (v *VK) GetCallbackServers(group_id int) ([]ServerItem, error) { //groups.getCallbackServers
	u := "access_token=" + v.Token +
		"&group_id=" + strconv.Itoa(group_id) +
		"&v=" + v.Version

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

type AutoSetResponse struct {
	Response struct {
		ServerID int `json:"server_id"`
	} `json:"response"`
}

func (v *VK) AddCallbackServer(u string) (int, error) {
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

	resp := AutoSetResponse{}

	err = json.Unmarshal(data, &resp)

	if err != nil {
		return 0, err
	}

	return resp.Response.ServerID, nil
}

type ResponseServerDeleteOrSet struct {
	Response int `json:"response"`
}

func (v *VK) DeleteCallbackServer(group_id, serv_id int) (int, error) {
	u := "access_token=" + v.Token +
		"&group_id=" + strconv.Itoa(-group_id) +
		"&server_id=" + strconv.Itoa(serv_id) +
		"&v=" + v.Version

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

type CallbackEvents string

const (
	MessageNew   CallbackEvents = "message_new"
	MessageReply CallbackEvents = "message_reply"
	MessageEdit  CallbackEvents = "message_edit"
)

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

func (v *VK) SetCallbackSettings(m *CallbackSettings) (int, error) {
	u := "access_token=" + v.Token +
		"&group_id=" + strconv.Itoa(int(math.Abs(float64(m.GroupID)))) +
		"&server_id=" + strconv.Itoa(m.ServerID) +
		"&message_new=" + strconv.FormatBool(m.MessageNew) +
		"&v=" + v.Version

	if m.MessageReply != false {
		u += "&message_reply=" + strconv.FormatBool(m.MessageReply)
	}

	if m.MessageEdit != false {
		u += "&message_edit=" + strconv.FormatBool(m.MessageEdit)
	}

	if m.MessageAllow != false {
		u += "&message_allow=" + strconv.FormatBool(m.MessageAllow)
	}

	if m.MessageDeny != false {
		u += "&message_deny=" + strconv.FormatBool(m.MessageDeny)
	}

	if m.MessageEvent != false {
		u += "&message_event=" + strconv.FormatBool(m.MessageEvent)
	}

	if m.MessageTypingState != false {
		u += "&message_typing_state=" + strconv.FormatBool(m.MessageTypingState)
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
