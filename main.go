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

	vk.UserToken = "e96350e10f9dbf5f81aff0919bf69c879769a716c3e0d4cdce190cf5d10b3907c1f69132092786250e5d1"

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
		MessageEvent:       MessageEventt,
	}

	clbck.Settings = append(clbck.Settings, callbackApi.MessageEdit, callbackApi.MessageNew,
		callbackApi.MessageReply, callbackApi.MessageAllow, callbackApi.MessageDeny, callbackApi.MessageTypingState, callbackApi.MessageEvent)

	go func() {
		_, _ = clbck.AutoConnect() //автоматическое подключение Callback API (без вмешательства пользователя)
	}()

	k := api.GetKeyboard()

	k.OneTime = true
	k.Inline = false

	/* k.AddButton(api.KeyboardButtons{
		Action: api.KeyboardActionTypeText{
			Type:    api.ActionText,
			Payload: api.ToPayload("1"),
			Label:   "Дратути",
		},
		Color: api.Primary,
	})

	k.AddButton(api.KeyboardButtons{
		Action: api.KeyboardActionTypeText{
			Type:    api.ActionText,
			Payload: api.ToPayload("2"),
			Label:   "До свидания",
		},
		Color: api.Secondary,
	})

	k.AddLine()

	k.AddButton(api.KeyboardButtons{
		Action: api.KeyboardActionTypeOpenLink{
			Type:    api.ActionOpenLink,
			Payload: api.ToPayload("3"),
			Label:   "click me",
			Link:    "https://www.google.com/",
		},
	}) */

	k.AddButton(api.KeyboardButtons{
		Action: api.KeyboardActionTypeCallback{
			Type:    api.ActionCallback,
			Payload: api.ToPayload("btn1"),
			Label:   "Привет",
		},
		Color: api.Primary,
	})

	log.Println(k.Buttons[0][0])

	at, err := vk.GetAttachments(api.GetAttachmentsParams{
		//FilePaths: []string{"C:/Users/a1exCross/Desktop/VKElmaLib/_iOe4_DihIE.jpg", "C:/Users/a1exCross/Desktop/Безымянный.png"},
		//FilePaths: []string{"C:/Users/a1exCross/Desktop/VKElmaLib/Король и Шут - Дагон.mp3"},
		//FilePaths: []string{"C:/Users/a1exCross/Desktop/VKElmaLib/Описание алгоритмов Заболотских, Иванов.docx", "C:/Users/a1exCross/Desktop/VKElmaLib/Ответы Караваева.docx"},
		FilePaths: []string{"C:/Users/a1exCross/Desktop/VKElmaLib/Ви део.mp4"},
		//FilePaths:     []string{"C:/Users/a1exCross/Desktop/VKElmaLib/_iOe4_DihIE.jpg", "C:/Users/a1exCross/Desktop/VKElmaLib/Описание алгоритмов Заболотских, Иванов.docx"},
		PeerID: 106988557,
		//OriginalPhoto: true,
		//OriginalVideo: true,
	})

	if err != nil {
		log.Println(err)
	}

	id, err := vk.MessagesSend(api.MessagesSendParams{
		PeerIDs:  []int{106988557},
		Message:  "message",
		RandomID: 0,
		//Keyboard: k,
	})

	gr, _ := vk.GetCurrentGroup()

	res_edit, err := vk.MessagesEdit(api.MessagesEditParams{
		PeerID:                id.Response[0].PeerID,
		ConversationMessageID: id.Response[0].ConversationMessageID,
		Message:               "message111",
		Attachment:            at,
		GroupID:               gr.Response[0].ID,
	})

	log.Println(res_edit)

	if err != nil {
		log.Println(err)
	}

	/* res_del, err := vk.MessagesDelete(api.MessagesDeleteParams{
		//MessageIDs: []int{id.Response[0].ConversationMessageID},
		//GroupID:      gr.Response[0].ID,
		DeleteForAll: false,
		PeerID:       id.Response[0].PeerID,
		//Cmids: ,
	})

	if err != nil {
		log.Println(err)
	}

	log.Println(res_del) */

	/* url, err := oauth.ImplictFlow(oauth.AuthParams{
		Client_ID: 8117272,
		Display:   oauth.Page,
		Scope:     []oauth.Scope{oauth.Video},
	})

	if err != nil {
		log.Println(err)
	}

	fmt.Println(url) */

	/* 	at, err := vk.GetAttachments(api.GetAttachmentsParams{
	   		//FilePaths: []string{"C:/Users/a1exCross/Desktop/VKElmaLib/_iOe4_DihIE.jpg", "C:/Users/a1exCross/Desktop/Безымянный.png"},
	   		//FilePaths: []string{"C:/Users/a1exCross/Desktop/VKElmaLib/Король и Шут - Дагон.mp3"},
	   		//FilePaths: []string{"C:/Users/a1exCross/Desktop/VKElmaLib/Описание алгоритмов Заболотских, Иванов.docx", "C:/Users/a1exCross/Desktop/VKElmaLib/Ответы Караваева.docx"},
	   		//FilePaths: []string{"C:/Users/a1exCross/Desktop/VKElmaLib/Видео.mp4"},
	   		//FilePaths:     []string{"C:/Users/a1exCross/Desktop/VKElmaLib/_iOe4_DihIE.jpg", "C:/Users/a1exCross/Desktop/VKElmaLib/Описание алгоритмов Заболотских, Иванов.docx"},
	   		PeerID:        106988557,
	   		OriginalPhoto: true,
	   		OriginalVideo: true,
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
	   	} */

	http.HandleFunc("/callback", clbck.HandleFunc)
	http.ListenAndServe(":80", nil)
}

func MessageFromUser(e callbackApi.Events, obj callbackApi.MessageObject) { //callback функция для отслеживания новых сообщений
	log.Println("Пользователь с идентификатором", obj.Message.FromID, "отправил сообщение", obj.Message.Text)

	if obj.Message.Payload != nil {
		log.Println(obj.Message.Payload.Payload)
	}

	if obj.Message.Geo != nil {
		log.Println(obj.Message.Geo.Coordinates)
	}

	if obj.Message.Attachments != nil {
		for _, v := range obj.Message.Attachments {
			if v.Type == callbackApi.Photo {
				log.Println(v.Photo.GetMaxSizePhotoUrl().Url)
			}

			if v.Type == callbackApi.Video {
				r, _ := vk.GetVideo(api.GetVideoParams{
					OwnerID: v.Video.OwnerID,
					Videos:  strconv.Itoa(v.Video.OwnerID) + "_" + strconv.Itoa(v.Video.ID),
				})

				log.Println(r.Response.Items[0].Player)
			}

			if v.Type == callbackApi.Doc {
				log.Println(v.Doc.Title)
			}

			if v.Type == callbackApi.Audio {
				log.Println(v.Audio.Title)
			}

			if v.Type == callbackApi.AudioMessage {
				log.Println(v.AudioMessage.OwnerID)
			}

			if v.Type == callbackApi.Graffiti {
				log.Println(v.Graffiti.URL)
			}

			if v.Type == callbackApi.Sticker {
				log.Println(v.Sticker.AnimationURL)
			}
		}
	}
}

func MessageFromGroup(e callbackApi.Events, obj callbackApi.MessageObjectMessage) {
	log.Println("Пользователь с идентификатором", e.GroupID, "отправил сообщение", obj.Text)

	_, err := vk.MessagesGetByConversationMessageID(api.MessagesGetByConversationMessageIDParams{
		PeerID:                 obj.PeerID,
		ConversationMessageIDs: []int{obj.ConversationMessageID},
		Extended:               false,
		GroupID:                e.GroupID,
	})

	if err != nil {
		log.Println("conv:", err)
	}
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

func MessageEventt(e callbackApi.Events, obj callbackApi.MessageEventObject) {
	log.Println(fmt.Sprintf("Нажата callback кнопка в чат-боте %s", obj.Payload.Payload.Button))

	_, err := vk.SendMessageEventAnswer(api.SendMessageEventAnswerParams{
		EventID: obj.EventID,
		UserID:  obj.UserID,
		PeerID:  obj.PeerID,
		EventData: api.EventAnswerType{
			ShowSnackbar: &api.ShowSnackbarAnswerType{
				Text: "Hello!!",
			},
		},
		/* 	EventData: api.EventAnswerType{
			OpenLink: &api.OpenLinkAnswerType{
				Link: "https://www.google.com",
			},
		}, */
	})

	if err != nil {
		log.Println(err)
	}
}
