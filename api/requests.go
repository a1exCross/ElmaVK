package api

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/a1exCross/ElmaVK/vkerrors"
)

const api_host = "https://api.vk.com/method/"

func (v VK) reqeustApiGet(method, param string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, api_host+method+param, nil)
	if err != nil {
		return nil, err
	}

	res, err := v.Client.Do(req)

	check := vkerrors.GetError(res)

	if check != "ok" {
		return nil, errors.New(check)
	}

	return res, nil
}

func (v VK) reqeustApiPost(method, param string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, api_host+method+param, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := v.Client.Do(req)

	check := vkerrors.GetError(res)

	if check != "ok" {
		return nil, errors.New(check)
	}

	return res, nil
}
