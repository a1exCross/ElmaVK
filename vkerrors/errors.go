//Проверка на наличие ошибок от API ВКонтакте
package vkerrors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetError(r *http.Response) string {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err.Error()
	}

	r.Body.Close()

	r.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	var b BodyError
	b.Error.ErrorCode = -1

	err = json.Unmarshal(data, &b)
	if err != nil {
		return err.Error()
	}

	if b.Error.ErrorCode != 200 && b.Error.ErrorCode != -1 {
		return "Error " + fmt.Sprint(b.Error.ErrorCode) + ": " + b.Error.ErrorMsg
	}

	return "ok"
}

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
