package apiErrors

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type BodyError struct {
	Error struct {
		ErrorCode     int    `json:"error_code"`
		ErrorMsg      string `json:"error_msg"`
		RequestParams []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"request_params"`
	} `json:"error"`
}

func GetError(r *http.Response) string {
	data, err := ioutil.ReadAll(r.Body)

	r.Body.Close()

	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	if err != nil {
		return err.Error()
	}

	var b BodyError
	b.Error.ErrorCode = -1

	err = json.Unmarshal(data, &b)

	if err != nil {
		return err.Error()
	}

	//log.Println(b)

	if b.Error.ErrorCode != 200 && b.Error.ErrorCode != -1 /* && len(b.Error.RequestParams) > 0 */ {
		return "Error " + strconv.Itoa(b.Error.ErrorCode) + ": " + b.Error.ErrorMsg //+ ". Method: " + b.Error.RequestParams[2].Value
	}

	return "ok"
}
