package longpool

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	ApiErrors "github.com/a1exCross/ElmaVK/errors"
	vk "github.com/a1exCross/ElmaVK/methods"
)

type Longpool struct {
	Server  string
	Key     string
	Ts      int
	Pts     int
	Wait    int
	VK      vk.VK
	Client  *http.Client
	Version int
	Mode    int
}

func New(v vk.VK, p GetLongPoolServerParams) (*Longpool, error) {
	lp := Longpool{
		Wait:    25,
		Version: 3,
		VK:      v,
		Client:  http.DefaultClient,
		Mode:    p.Mode,
	}

	err := lp.getLongPoolServer(p)

	if err != nil {
		return nil, err
	}

	return &lp, nil
}

type LongpoolParams struct {
	Server string `json:"server"`
	Key    string `json:"key"`
	Ts     int    `json:"ts"`
	Pts    int    `json:"pts"`
}

type GetLongPoolServerResponse struct {
	Response LongpoolParams `json:"response"`
}

type GetLongPoolServerParams struct {
	NeedPTS   bool `json:"need_pts"`
	GroupID   int  `json:"group_id,omitempty"`
	LpVersion int  `json:"lp_version"`
	Mode      int
}

func (lp *Longpool) getLongPoolServer(p GetLongPoolServerParams) error {
	var u string = ""

	if p.GroupID != 0 {
		u += "group_id=" + strconv.Itoa(p.GroupID)
	}

	u += "&lp_version=" + strconv.Itoa(lp.Version)

	u += "&need_pts=" + strconv.FormatBool(p.NeedPTS)

	if lp.VK.Token != "" {
		u += "&v=" + lp.VK.Version + "&access_token=" + lp.VK.Token
	} /* else {
		return nil, errors.New("Token group is null")
	} */

	res, err := lp.VK.Reqeust_api_get("messages.getLongPollServer?", u)

	if err != nil {
		return err
	}

	check := ApiErrors.GetError(res)

	if check != "ok" {
		return errors.New(check)
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	log.Println(string(data))

	var r GetLongPoolServerResponse

	err = json.Unmarshal(data, &r)

	if err != nil {
		return err
	}

	lp.Key = r.Response.Key
	lp.Pts = r.Response.Pts
	lp.Server = r.Response.Server
	lp.Ts = r.Response.Ts

	return nil
}

func (lp *Longpool) Run() error {

	for true {
		r, err := lp.check_events()

		if err != nil {
			return err
		}

		for _, event := range r.Updates {
			log.Println("Length:", len(event))
			log.Println("Event=", event[0])
			fmt.Printf("%T\n", event[0])
			if event[0] == 4.0 {
				log.Println(event[0])
				log.Println(event[1])
				log.Println(event[2])
				log.Println(event[3])
				log.Println(event[4])
				log.Println(event[5])
				log.Println(event[6])
				log.Println(event[7])
			}
		}
	}

	return nil
}

type LongpoolResponse struct {
	Ts      int             `json:"ts"`
	Updates [][]interface{} `json:"updates"`
}

func (lp *Longpool) check_events() (*LongpoolResponse, error) {
	var u string = ""

	/* if p.GroupID != 0 {
		u += "group_id=" + strconv.Itoa(p.GroupID)
	}

	if p.LpVersion != 0 {
		p.LpVersion = 3
		u += "&lp_version=" + strconv.Itoa(p.LpVersion)
	}

	u += "&need_pts=" + strconv.FormatBool(p.NeedPTS)

	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.Token
	}  */

	u += "act=a_check&key=" + lp.Key + "&ts=" + strconv.Itoa(lp.Ts) + "&wait=" + strconv.Itoa(lp.Wait) + "&mode=" + strconv.Itoa(lp.Mode) + "&version=" + strconv.Itoa(lp.Version)

	res, err := lp.Reqeust_api_get(lp.Server+"?", u)
	//{$server}?act=a_check&key={$key}&ts={$ts}&wait=25&mode=2&version=2

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

	log.Println(string(data))

	var r LongpoolResponse

	err = json.Unmarshal(data, &r)

	if err != nil {
		return nil, err
	}

	lp.Ts = r.Ts

	return &r, nil
}
