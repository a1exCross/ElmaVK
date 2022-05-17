package longpool

import (
	"net/http"
)

func (lp *Longpool) Reqeust_api_get(method, param string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, "http://"+method+param, nil)

	if err != nil {
		return nil, err
	}

	return lp.Client.Do(req)
}
