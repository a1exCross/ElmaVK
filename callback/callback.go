//Пакет для работы с Callback API
package callback

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

func New(v vk.VK) Callback {
	return Callback{
		Vk: v,
	}
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

	if e.Secret == c.Secret_key {

		c.callFuncList(data, e)

		if e.Type == "confirmation" {
			fmt.Fprintf(w, c.ConfirmationKey)
			return
		}
	}

	_, err = w.Write([]byte("ok"))
	if err != nil {
		log.Println(err)
	}
}
