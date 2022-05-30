package callbackApi

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	vk "github.com/a1exCross/ElmaVK/api"
)

func (c *Callback) AutoConnect() (int, error) {
	group, err := c.Vk.GetCurrentGroup()

	if err != nil {
		return 0, err
	}

	servers, err := c.Vk.GetCallbackServers(group.Response[0].ID)

	if err != nil {
		return 0, err
	}

	server_id := 0

	for _, s := range servers {
		if s.URL == c.URL {
			if s.Status != "ok" {
				res, err := c.Vk.DeleteCallbackServer(s.CreatorID, s.ID)
				if err != nil {
					return 0, err
				}

				if res == 1 {
					log.Println("Server with ID =", s.ID, "deleted")
				}

			} else {
				server_id = s.ID
				c.Title = s.Title
				c.Secret_key = s.SecretKey
				log.Println("Work server is finded with ID =", s.ID)
				break
			}
		}
	}

	if server_id == 0 {
		u, err := c.BuildRequestAddCallbackServer()

		if err != nil {
			return 0, err
		}

		c.ConfirmationKey, err = c.Vk.GetConfirmaionKey(group.Response[0].ID)

		if err != nil {
			return 0, err
		}

		res, err := c.Vk.AddCallbackServer(u)

		server_id = res

		if err != nil {
			return 0, err
		}

		log.Println("New server is created with ID =", res)
	}

	settings := setSettings(c.Vk, c.Settings, group.Response[0].ID, server_id)

	resp, err := c.Vk.SetCallbackSettings(settings)

	if err != nil {
		return 0, errors.New(err.Error())
	}

	if resp == 1 {
		log.Println("Server", c.Title, "is configured")
	}

	return server_id, nil
}

func setSettings(v vk.VK, e []EventType, group_id, server_id int) vk.CallbackSettings {
	settings := v.GetCallbackSettingsParams()
	settings.GroupID = group_id
	settings.ServerID = server_id

	for _, v := range e {
		switch v {
		case MessageEdit:
			{
				settings.MessageEdit = true
			}
		case MessageNew:
			{
				settings.MessageNew = true
			}

		case MessageReply:
			{
				settings.MessageReply = true
			}
		case MessageAllow:
			{
				settings.MessageAllow = true
			}
		case MessageDeny:
			{
				settings.MessageDeny = true
			}
		case MessageTypingState:
			{
				settings.MessageTypingState = true
			}
		case MessageEvent:
			{
				settings.MessageEvent = true
			}
		}
	}
	return settings
}

func (c Callback) BuildRequestAddCallbackServer() (string, error) {
	var u string = ""

	group, err := c.Vk.GetCurrentGroup()

	if err != nil {
		return "", err
	}

	if group != nil {
		u += "group_id=" + fmt.Sprint(group.Response[0].ID)
	} else {
		return "", errors.New("Required field 'GroupID' is empty, MethodName - BuildRequestAddCallbackServer()")
	}

	if c.Title == "" {
		c.Title = "VKGroupELMA365"
	} else {
		u += "&title=" + c.Title
	}

	if c.URL != "" {
		u += "&url=" + c.URL
	} else {
		return "", errors.New("Required field 'URL' is empty, MethodName - BuildRequestAddCallbackServer()")
	}

	if c.Vk.Token != "" {
		u += "&access_token=" + c.Vk.Token
	} else {
		return "", errors.New("Auth token is empty")
	}

	if c.Secret_key == "" {
		key := random_string()
		u += "&secret_key=" + key
		c.Secret_key = key
	} else {
		u += "&secret_key=" + c.Secret_key
	}

	u += "&v=" + c.Vk.Version

	return u, nil
}

func random_string() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := rand.Intn(50)
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
