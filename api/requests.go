package api

import (
	"net/http"
	"net/url"
	"strings"
)

const api_host = "https://api.vk.com/method/"

func (v VK) Reqeust_api_get(method, param string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, api_host+method+param, nil)

	if err != nil {
		return nil, err
	}

	return v.Client.Do(req)
}

func (v VK) Reqeust_api_post(method, param string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, api_host+method+param, strings.NewReader(data.Encode()))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return v.Client.Do(req)
}
