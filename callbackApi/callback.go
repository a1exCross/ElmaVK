package callbackApi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	//events "vksdk/callbackApi/events"
	vk "github.com/a1exCross/ElmaVK/api"
)

type Callback struct {
	Secret_key      string
	ConfirmationKey string
	Title           string
	Vk              vk.VK
	URL             string
	Settings        []EventType
	Functions       FuncList
}

/* func GetCallbackParams() *callback{
	return &callback{}
} */

func New() Callback {
	return Callback{}
}

func (c *Callback) HandleFunc(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
	}

	defer r.Body.Close()

	var e Events

	err = json.Unmarshal(data, &e)

	if err != nil {
		log.Println(err)
	}

	log.Println(e.Type)

	log.Println(string(data))

	if e.Secret == c.Secret_key {

		c.CallFuncList(data, e)

		if e.Type == "confirmation" {
			fmt.Fprintf(w, c.ConfirmationKey)
			return
		}
	}

	_, _ = w.Write([]byte("ok"))
}

/* type CallbackServersArray struct {
	Items []vk.ServerItem
} */
