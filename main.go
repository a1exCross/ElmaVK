package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/a1exCross/ElmaVK/callbackApi"

	//longpool "github.com/a1exCross/VkElmaLib/longpool-user"
	"github.com/a1exCross/ElmaVK/api"
)

var vk = api.Session("9e63c5f226af4bb3bec5ce0d83e5a093e6ae2d5b034be0d67b6c0f87a46f0b7fcad7d8c7a2c76a283d787")

func main() {
	/* lp, err := longpool.New(vk, longpool.GetLongPoolServerParams{
		NeedPTS:   true,
		LpVersion: 3,
		Mode:      2,
	})

	if err != nil {
		log.Println(err)
	}

	lp.Run() */

	////// Авторизация
	/* auth, u := oauth.AuthCodeFlow(oauth.AuthParams{
		Client_ID:    v.Client_ID,                                  // идентификатор приложения ВК
		Display:      oauth.Page,                                   // формат отображения страницы при авторизации
		Group_IDs:    v.Group_IDs,                                  // идентификаторы групп
		Scope:        []oauth.Scope{oauth.Messages, oauth.Manage},  // рарещение прав при авторизации
		Redirect_URI: "https://e5c7-94-241-222-102.ngrok.io/token", // на данный адрес будет отправлен код доступа, необходимый для получения токена авторизации
		CleintSecret: "ClientSecret",
	}) // защищенный ключ приложения ВК

	// необходимо перенаправить пользователя по адресу в переменной u для авторизации

	group_tokens := auth.GetToken(r.URL) //r.URL - URL, с кодом доступа, полученный после авторизации */
	//////

	//group_tokens - массив с полученными токенами для указанных групп

	clbck := callbackApi.New()

	vk.UserToken = "3a9d86957773b960b7b2c34998124103c9ddfaaf8654ad6841c76ccf5e4671e9b06cac198b75dfbe5e7ec"

	clbck.Vk = vk

	clbck.Title = "GroupTitle"                                  //название группы
	clbck.URL = "https://0d1d-89-254-254-209.ngrok.io/callback" //адрес, на который будут приходить уведомления

	clbck.Functions = callbackApi.FuncList{ //список функций для отслеживания определенных событиый
		NewMessage:         MessageFromUser,
		MessageReply:       MessageFromGroup,
		MessageEdit:        MessageEditt,
		MessageTypingState: MessageTypingStatee,
		MessageAllow:       MessageAlloww,
		MessageDeny:        MessageDenyy,
	}

	clbck.Settings = append(clbck.Settings, callbackApi.MessageEdit, callbackApi.MessageNew,
		callbackApi.MessageReply, callbackApi.MessageAllow, callbackApi.MessageDeny, callbackApi.MessageTypingState)

	_, _ = clbck.AutoConnect() //автоматическое подключение Callback API (без вмешательства пользователя)

	//log.Println(s.Response[0])

	//

	/* url, err := oauth.ImplictFlow(oauth.AuthParams{
		Client_ID: 8117272,
		Display:   oauth.Page,
		Scope:     []oauth.Scope{oauth.Video},
	})

	fmt.Println(url) */

	at, err := vk.GetAttachments(api.GetAttachmentsParams{
		//FilePaths: []string{"C:/Users/a1exCross/Desktop/VKElmaLib/_iOe4_DihIE.jpg", "C:/Users/a1exCross/Desktop/Безымянный.png"},
		//FilePaths: []string{"C:/Users/a1exCross/Desktop/VKElmaLib/Король и Шут - Дагон.mp3"},
		//FilePaths: []string{"C:/Users/a1exCross/Desktop/VKElmaLib/Описание алгоритмов Заболотских, Иванов.docx", "C:/Users/a1exCross/Desktop/VKElmaLib/Ответы Караваева.docx"},
		//FilePaths: []string{"C:/Users/a1exCross/Desktop/VKElmaLib/Видео.mp4"},
		FilePaths:     []string{"C:/Users/a1exCross/Desktop/VKElmaLib/_iOe4_DihIE.jpg", "C:/Users/a1exCross/Desktop/VKElmaLib/Описание алгоритмов Заболотских, Иванов.docx"},
		PeerID:        106988557,
		OriginalPhoto: true,
	})

	if err != nil {
		log.Println(err)
	} else {
		vk.SendMessage(api.SendMessage{
			PeerIDs:    []int{106988557},
			Message:    "message",
			RandomID:   0,
			Attachment: at,
		})
	}

	http.HandleFunc("/callback", clbck.HandleFunc)
	http.ListenAndServe(":80", nil)
}

func MessageFromUser(e callbackApi.Events, obj callbackApi.MessageObject) { //callback функция для отслеживания новых сообщений
	log.Println("Пользователь с идентификатором", obj.Message.FromID, "отправил сообщение", obj.Message.Text)

	if obj.Message.Geo != nil {
		log.Println(obj.Message.Geo.Coordinates)
	}

	if obj.Message.Attachments != nil {
		for _, v := range obj.Message.Attachments {
			if v.Type == string(callbackApi.Photo) {
				log.Println(v.Photo.GetMaxSizePhotoUrl().Url)
			}

			if v.Type == string(callbackApi.Video) {
				r, _ := vk.GetVideo(api.GetVideoParams{
					OwnerID: v.Video.OwnerID,
					Videos:  strconv.Itoa(v.Video.OwnerID) + "_" + strconv.Itoa(v.Video.ID),
				})

				log.Println(r.Response.Items[0].Player)
			}

			if v.Type == string(callbackApi.Doc) {
				log.Println(v.Doc.Title)
			}

			if v.Type == string(callbackApi.Audio) {
				log.Println(v.Audio.Title)
			}

			if v.Type == string(callbackApi.AudioMessage) {
				log.Println(v.AudioMessage.OwnerID)
			}

			if v.Type == string(callbackApi.Graffiti) {
				log.Println(v.Graffiti.URL)
			}

			if v.Type == string(callbackApi.Sticker) {
				log.Println(v.Sticker.AnimationURL)
			}
		}
	}
}

func MessageFromGroup(e callbackApi.Events, obj callbackApi.MessageObjectMessage) {
	log.Println("Пользователь с идентификатором", e.GroupID, "отправил сообщение", obj.Text)
}

func MessageEditt(e callbackApi.Events, obj callbackApi.MessageObjectMessage) {
	log.Println(fmt.Sprintf("%d отредактировал сообщение", obj.FromID))
}

func MessageTypingStatee(e callbackApi.Events, obj callbackApi.MessageTypingStateObject) {
	log.Println(fmt.Sprintf("Пользователь с ID = %d набирает сообщение для ID %d", obj.FromID, obj.ToID))
}

func MessageDenyy(e callbackApi.Events, obj callbackApi.MessageDenyObject) {
	log.Println(fmt.Sprintf("Пользователь %d запретил сообщения от сообщества", obj.UserID))
}

func MessageAlloww(e callbackApi.Events, obj callbackApi.MessageAllowObject) {
	log.Println(fmt.Sprintf("Пользователь %d разрешил сообщения от сообщества", obj.UserID))
}
