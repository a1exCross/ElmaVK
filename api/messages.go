package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/a1exCross/ElmaVK/vkerrors"
)

type ShowSnackbarAnswerType struct {
	Text string
}

type OpenLinkAnswerType struct {
	Link string
}

type OpenAppAnswerType struct {
	AppID   int
	OwnerID int
	Hash    string
}

type EventAnswerType struct {
	Type         string
	ShowSnackbar *ShowSnackbarAnswerType
	OpenLink     *OpenLinkAnswerType
	OpenApp      *OpenAppAnswerType
}

type SendMessageEventAnswerParams struct {
	EventID   string          `json:"event_id"`
	UserID    int             `json:"user_id"`
	PeerID    int             `json:"peer_id"`
	EventData EventAnswerType `json:"event_data"`
}

func toJsonParam(p1, p2 string) string {
	return fmt.Sprintf("\"%s\":\"%s\"", p1, p2)
}

//https://dev.vk.com/method/messages.sendMessageEventAnswer
func (v VK) SendMessageEventAnswer(p SendMessageEventAnswerParams) (SendMessageEventAnswerResponse, error) {
	data := url.Values{}

	if p.EventID != "" {
		data.Set("event_id", p.EventID)
	} else {
		return SendMessageEventAnswerResponse{}, errors.New("Required field 'EventID' is empty, MethodName - SendMessageEventAnswer()")
	}

	if p.PeerID != 0 {
		data.Set("peer_id", fmt.Sprint(p.PeerID))
	} else {
		return SendMessageEventAnswerResponse{}, errors.New("Required field 'PeerID' is empty, MethodName - SendMessageEventAnswer()")
	}

	if p.UserID != 0 {
		data.Set("user_id", fmt.Sprint(p.UserID))
	} else {
		return SendMessageEventAnswerResponse{}, errors.New("Required field 'UserID' is empty, MethodName - SendMessageEventAnswer()")
	}

	var prm string

	if p.EventData.OpenLink != nil {
		p.EventData.Type = "open_link"
		prm = "{" + toJsonParam("link", p.EventData.OpenLink.Link) + "," + toJsonParam("type", p.EventData.Type) + "}"
	}

	if p.EventData.OpenApp != nil {
		p.EventData.Type = "open_app"
		if p.EventData.Type == "open_app" {
			prm = "{" + toJsonParam("app_id", fmt.Sprint(p.EventData.OpenApp.AppID)) + "," +
				toJsonParam("owner_id", fmt.Sprint(p.EventData.OpenApp.OwnerID)) + "," +
				toJsonParam("hash", p.EventData.OpenApp.Hash) + "}"
		}
	}

	if p.EventData.ShowSnackbar != nil {
		p.EventData.Type = "show_snackbar"
		if p.EventData.Type == "show_snackbar" {
			prm = "{" + toJsonParam("text", p.EventData.ShowSnackbar.Text) + "," + toJsonParam("type", p.EventData.Type) + "}"
		}
	}

	data.Set("event_data", prm)

	var u string = ""

	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.Token
	} else {
		return SendMessageEventAnswerResponse{}, errors.New("Auth token is empty")
	}

	res, err := v.Reqeust_api_post("messages.sendMessageEventAnswer?", u, data)
	if err != nil {
		return SendMessageEventAnswerResponse{}, err
	}

	check := vkerrors.GetError(res)
	if check != "ok" {
		return SendMessageEventAnswerResponse{}, errors.New(check)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return SendMessageEventAnswerResponse{}, err
	}

	var t SendMessageEventAnswerResponse

	err = json.Unmarshal(body, &t)
	if err != nil {
		return SendMessageEventAnswerResponse{}, err
	}

	return t, nil
}

type SendMessageEventAnswerResponse struct {
	Response int `json:"response"`
}

type MessagesGetByConversationMessageIDParams struct {
	PeerID                 int      `json:"peer_id"`
	ConversationMessageIDs []int    `json:"conversation_messages_id"`
	Extended               bool     `json:"extended"`
	Feilds                 []string `json:"fields,omitempty"`
	GroupID                int      `json:"group_id"`
}

//https://dev.vk.com/method/messages.getByConversationMessageId
func (v VK) MessagesGetByConversationMessageID(p MessagesGetByConversationMessageIDParams) (MessagesGetByConversationMessageIDResponse, error) {
	data := url.Values{}

	if p.PeerID != 0 {
		data.Set("peer_id", fmt.Sprint(p.PeerID))
	} else {
		return MessagesGetByConversationMessageIDResponse{}, errors.New("Required field 'PeerID' is empty, MethodName - MessagesGetByConversationMessageID()")
	}

	if p.ConversationMessageIDs != nil {
		conv_ids := ""
		for i, v := range p.ConversationMessageIDs {
			if i > 0 {
				conv_ids += ","
			}
			conv_ids += fmt.Sprint(v)
		}
		data.Set("conversation_message_ids", conv_ids)
	} else {
		return MessagesGetByConversationMessageIDResponse{}, errors.New("Required field 'ConversationMessageIDs' is empty, MethodName - MessagesGetByConversationMessageID()")
	}

	if p.Feilds != nil {
		fields := ""
		for i, v := range p.Feilds {
			if i > 0 {
				fields += ","
			}
			fields += fmt.Sprint(v)
		}
		data.Set("fields", fields)
	}

	data.Set("extended", fmt.Sprint(p.Extended))

	if p.GroupID != 0 {
		data.Set("group_id", fmt.Sprint(p.GroupID))
	}

	var u string = ""

	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.Token
	} else {
		return MessagesGetByConversationMessageIDResponse{}, errors.New("Auth token is empty")
	}

	res, err := v.Reqeust_api_post("messages.getByConversationMessageId?", u, data)
	if err != nil {
		return MessagesGetByConversationMessageIDResponse{}, err
	}

	check := vkerrors.GetError(res)
	if check != "ok" {
		return MessagesGetByConversationMessageIDResponse{}, errors.New(check)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return MessagesGetByConversationMessageIDResponse{}, err
	}

	var t MessagesGetByConversationMessageIDResponse

	err = json.Unmarshal(body, &t)
	if err != nil {
		return MessagesGetByConversationMessageIDResponse{}, err
	}

	return t, nil
}

type MessagesGetByConversationMessageIDResponse struct {
	Response struct {
		Count int `json:"count"`
		Items []struct {
			Date                  int           `json:"date"`
			FromID                int           `json:"from_id"`
			ID                    int           `json:"id"`
			Out                   int           `json:"out"`
			Attachments           []interface{} `json:"attachments"`
			ConversationMessageID int           `json:"conversation_message_id"`
			FwdMessages           []interface{} `json:"fwd_messages"`
			Important             bool          `json:"important"`
			IsHidden              bool          `json:"is_hidden"`
			PeerID                int           `json:"peer_id"`
			RandomID              int           `json:"random_id"`
			Text                  string        `json:"text"`
		} `json:"items"`
		Profiles []struct {
			ID         int    `json:"id"`
			Sex        int    `json:"sex"`
			ScreenName string `json:"screen_name"`
			Photo50    string `json:"photo_50"`
			Photo100   string `json:"photo_100"`
			OnlineInfo struct {
				Visible  bool `json:"visible"`
				LastSeen int  `json:"last_seen"`
				IsOnline bool `json:"is_online"`
				AppID    int  `json:"app_id"`
				IsMobile bool `json:"is_mobile"`
			} `json:"online_info"`
			Online          int    `json:"online"`
			OnlineMobile    int    `json:"online_mobile"`
			OnlineApp       int    `json:"online_app"`
			FirstName       string `json:"first_name"`
			LastName        string `json:"last_name"`
			CanAccessClosed bool   `json:"can_access_closed"`
			IsClosed        bool   `json:"is_closed"`
		} `json:"profiles"`
		Groups []struct {
			ID         int    `json:"id"`
			Name       string `json:"name"`
			ScreenName string `json:"screen_name"`
			IsClosed   int    `json:"is_closed"`
			Type       string `json:"type"`
			Photo50    string `json:"photo_50"`
			Photo100   string `json:"photo_100"`
			Photo200   string `json:"photo_200"`
		} `json:"groups"`
	} `json:"response"`
}

type MessagesEditParams struct {
	PeerID                int
	ConversationMessageID int
	Message               string
	Attachment            []string
	GroupID               int
	Keyboard              Keyboard `json:"keyboard"`
	Template              string
	DontParseLinks        bool
	DisableMentions       bool
	MessageID             int
	Lat                   string
	Long                  string
	KeepForwardMessages   bool
	KeepSnipets           bool
}

//https://dev.vk.com/method/messages.edit
func (v VK) MessagesEdit(p MessagesEditParams) (int, error) {
	data := url.Values{}

	if p.PeerID != 0 {
		data.Set("peer_id", fmt.Sprint(p.PeerID))
	} else {
		return 0, errors.New("Required field 'PeerID' is empty, MethodName - MessagesEdit()")
	}

	if p.ConversationMessageID != 0 {
		data.Set("conversation_message_id", fmt.Sprint(p.ConversationMessageID))
	}

	if p.GroupID != 0 {
		data.Set("group_id", fmt.Sprint(p.GroupID))
	}

	if p.Attachment != nil {
		var att = ""
		for i := 0; i < len(p.Attachment); i++ {
			if i > 0 {
				att += ","
			}
			att += p.Attachment[i]
		}
		data.Set("attachment", att)
	}

	if p.Message != "" {
		data.Set("message", p.Message)
	} else if p.Attachment == nil {
		return 0, errors.New("Required field 'Message' is empty, MethodName - MessagesEdit()")
	}

	if p.Template != "" {
		data.Set("template", p.Template)
	}

	if p.MessageID != 0 {
		data.Set("message_id", fmt.Sprint(p.MessageID))
	}

	data.Set("dont_parse_links", fmt.Sprint(p.DontParseLinks))

	data.Set("disable_mentions", fmt.Sprint(p.DisableMentions))

	if p.Lat != "" {
		data.Set("lat", p.Lat)
	}

	if p.Long != "" {
		data.Set("long", p.Long)
	}

	data.Set("keep_forward_messages", fmt.Sprint(p.KeepForwardMessages))

	data.Set("keep_snippets", fmt.Sprint(p.KeepSnipets))

	if p.Keyboard.Buttons == nil {
		p.Keyboard.Buttons = [][]KeyboardButtons{}
	}

	ps, err := json.Marshal(p.Keyboard)
	if err != nil {
		return 0, err
	}

	data.Set("keyboard", string(ps))

	var u string = ""

	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.Token
	} else {
		return 0, errors.New("Auth token is empty")
	}

	res, err := v.Reqeust_api_post("messages.edit?", u, data)
	if err != nil {
		return 0, err
	}

	check := vkerrors.GetError(res)
	if check != "ok" {
		return 0, errors.New(check)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	var t MessagesEditResponse

	err = json.Unmarshal(body, &t)
	if err != nil {
		return 0, err
	}

	return t.Response, nil
}

type MessagesEditResponse struct {
	Response int `json:"response"`
}
type MessagesDeleteParams struct {
	MessageIDs   []int
	GroupID      int
	DeleteForAll bool
	PeerID       int
	//Cmids        []int
}

//https://dev.vk.com/method/messages.delete
func (v VK) MessagesDelete(p MessagesDeleteParams) (map[string]int, error) {
	data := url.Values{}

	var msg_ids = ""

	if p.MessageIDs != nil {
		for i, v := range p.MessageIDs {
			if i > 0 {
				msg_ids += ","
			}
			msg_ids += fmt.Sprint(v)
		}
	} else {
		return nil, errors.New("Required field 'MessagesIDs' is empty, MethodName - MessagesDelete()")
	}

	if p.PeerID != 0 {
		data.Set("peer_id", fmt.Sprint(p.PeerID))
	} else {
		return nil, errors.New("Required field 'PeerID' is empty, MethodName - MessagesDelete()")
	}

	if p.GroupID != 0 {
		data.Set("group_id", fmt.Sprint(p.GroupID))
	}

	data.Set("cmids", msg_ids)

	data.Set("delete_for_all", fmt.Sprint(p.DeleteForAll))

	var u string = ""

	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.Token
	} else {
		return nil, errors.New("Auth token is empty")
	}

	res, err := v.Reqeust_api_post("messages.delete?", u, data)
	if err != nil {
		return nil, err
	}

	check := vkerrors.GetError(res)
	if check != "ok" {
		return nil, errors.New(check)
	}

	dat, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var t MessagesDeleteResponse

	err = json.Unmarshal(dat, &t)
	if err != nil {
		return nil, err
	}

	var arr map[string]int
	err = json.Unmarshal(t.Response, &arr)

	if err != nil {
		log.Println(err)
	}

	return arr, nil
}

type MessagesDeleteResponse struct {
	Response json.RawMessage `json:"response"`
}

type GetAttachmentsParams struct {
	FilePaths     []string
	PeerID        int
	OriginalPhoto bool
	OriginalVideo bool
}

//https://dev.vk.com/reference/objects/attachments-message
func (v VK) GetAttachments(p GetAttachmentsParams) ([]string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	var at []string

	group_id, err := v.GetCurrentGroup()
	if err != nil {
		return nil, err
	}

	for _, path := range p.FilePaths {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		defer file.Close()

		filename := filepath.Base(path)

		ext := filepath.Ext(filename)

		ext = strings.ToLower(ext)

		up := ""
		attachment_type := ""
		var r SaveVideoResponse

		if (ext == ".jpg" && !p.OriginalPhoto) || (ext == ".png" && !p.OriginalPhoto) || (ext == ".gif" && !p.OriginalPhoto) {
			up, err = v.GetMessagesUploadServerPhoto(p.PeerID)
			if err != nil {
				return nil, err
			}

			attachment_type = "photo"
		} else if (ext == ".mp4" && !p.OriginalVideo) || (ext == ".avi" && !p.OriginalVideo) ||
			(ext == ".3gp" && !p.OriginalVideo) || (ext == ".mpeg" && !p.OriginalVideo) ||
			(ext == ".mov" && !p.OriginalVideo) || (ext == ".flw" && !p.OriginalVideo) ||
			(ext == ".wmv" && !p.OriginalVideo) /* || ext == ".mp3"  */ {

			title := strings.Replace(filename, ext, "", len(ext))

			r, err = v.SaveVideo(SaveVideoParams{
				Name:      title,
				IsPrivate: true,
				Wallpost:  false,
				GroupID:   group_id.Response[0].ID,
			})
			if err != nil {
				return nil, err
			}

			up = r.Response.UploadURL

			attachment_type = "video"
		} else if ext != ".mp3" && ext != ".exe" {
			up, err = v.GetMessagesUploadServerDoc("doc", p.PeerID)
			if err != nil {
				return nil, err
			}

			attachment_type = "doc"
		}

		if attachment_type == "" {
			return nil, errors.New("This file type is not supported")
		}

		part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return nil, err
		}

		err = writer.Close()
		if err != nil {
			return nil, err
		}

		res, err := http.DefaultClient.Post(up, writer.FormDataContentType(), body)
		if err != nil {
			return nil, err
		}

		check := vkerrors.GetError(res)
		if check != "ok" {
			return nil, errors.New(check)
		}

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		if attachment_type == "photo" {

			var r AfterUploadPhotoResponse

			err = json.Unmarshal(data, &r)
			if err != nil {
				return nil, err
			}

			s, err := v.SaveMessagesPhoto(r.Hash, r.Photo, r.Server)
			if err != nil {
				return nil, err
			}

			at = append(at, MediaToAttachment("photo", s.Response[0].OwnerID, s.Response[0].ID, ""))
		}

		if attachment_type == "doc" {
			var r AfterUploadDocResponse

			err = json.Unmarshal(data, &r)
			if err != nil {
				return nil, err
			}

			s, err := v.SaveDoc(SaveDocParams{
				File: r.File,
			})
			if err != nil {
				return nil, err
			}

			at = append(at, MediaToAttachment("doc", s.Response.Doc.OwnerID, s.Response.Doc.ID, ""))
		}

		if attachment_type == "video" {
			at = append(at, MediaToAttachment("video", r.Response.OwnerID, r.Response.VideoID, r.Response.AccessKey))
		}
	}

	return at, nil
}

type AfterUploadPhotoResponse struct {
	Server int    `json:"server"`
	Photo  string `json:"photo"`
	Hash   string `json:"hash"`
}

type AfterUploadDocResponse struct {
	File string `json:"file"`
}

type MessagesSendParams struct {
	UserID          int
	UserIDs         []int
	PeerID          int
	PeerIDs         []int
	Domain          string
	ChatID          int
	Message         string
	RandomID        int
	Guid            int
	Attachment      []string
	ReplyTo         int
	ForwardMessages []int
	//Forward json object
	StickerID int
	//GroupID   int
	Keyboard Keyboard `json:"keyboard"`
	Template string
	Payload  string
	//ContentSource text json
	DontParseLinks  bool
	DisableMentions bool
	//Intent          string
	//SubscribeID int
	Lat  string
	Long string
}

//https://dev.vk.com/method/messages.send
func (v VK) MessagesSend(p MessagesSendParams) (MessagesSendResponseIDs, MessagesSendResponseID, error) {
	// первый возвращаемый параметр функции возвращает структуру, если были переданы UserIDs или PeerIDs
	// второй возвращаемый параметр функции возвращает структуру, если были переданы UserID или PeerID
	data := url.Values{}

	if p.UserID != 0 {
		data.Set("user_id", fmt.Sprint(p.UserID))
	}

	if p.UserIDs != nil {
		var ids = ""
		for i := 0; i < len(p.UserIDs); i++ {
			if i > 0 {
				ids += ","
			}
			ids += fmt.Sprint(p.UserIDs[i])
		}
		data.Set("user_ids", ids)
	}

	if p.PeerID != 0 {
		data.Set("peer_id", fmt.Sprint(p.PeerID))
	}

	if p.PeerIDs != nil {
		var ids = ""
		for i := 0; i < len(p.PeerIDs); i++ {
			if i > 0 {
				ids += ","
			}
			ids += fmt.Sprint(p.PeerIDs[i])
		}
		data.Set("peer_ids", ids)
	}

	if p.Domain != "" {
		data.Set("domain", p.Domain)
	}

	if p.ChatID != 0 {
		data.Set("chat_id", fmt.Sprint(p.ChatID))
	}

	if p.Message != "" {
		data.Set("message", url.QueryEscape(p.Message))
	} else if p.Attachment == nil {
		return MessagesSendResponseIDs{}, MessagesSendResponseID{}, errors.New("Required field 'Message' is empty, MethodName - MessagesSend()")
	}

	if p.ForwardMessages != nil {
		var frw_msg = ""
		for i := 0; i < len(p.ForwardMessages); i++ {
			if i > 0 {
				frw_msg += ","
			}
			frw_msg += fmt.Sprint(p.ForwardMessages[i])
		}
		data.Set("forward_messages", frw_msg)
	}

	if p.Guid != 0 {
		data.Set("guid", fmt.Sprint(p.Guid))
	}

	if p.Attachment != nil {
		var att = ""
		for i := 0; i < len(p.Attachment); i++ {
			if i > 0 {
				att += ","
			}
			att += p.Attachment[i]
		}
		data.Set("attachment", att)
	}

	if p.Keyboard.Buttons == nil {
		p.Keyboard.Buttons = [][]KeyboardButtons{}
	}

	ps, err := json.Marshal(p.Keyboard)
	if err != nil {
		return MessagesSendResponseIDs{}, MessagesSendResponseID{}, err
	}

	data.Set("keyboard", string(ps))

	data.Set("random_id", fmt.Sprint(p.RandomID))

	if p.Template != "" {
		data.Set("template", p.Template)
	}

	data.Set("dont_parse_links", fmt.Sprint(p.DontParseLinks))

	data.Set("disable_mentions", fmt.Sprint(p.DisableMentions))

	if p.Lat != "" {
		data.Set("lat", p.Lat)
	}

	if p.Long != "" {
		data.Set("long", p.Long)
	}

	if p.ReplyTo != 0 {
		data.Set("reply_to", fmt.Sprint(p.ReplyTo))
	}

	if p.Payload != "" {
		data.Set("payload", p.Payload)
	}

	var u string = ""

	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.Token
	} else {
		return MessagesSendResponseIDs{}, MessagesSendResponseID{}, errors.New("Auth token is empty")
	}

	res, err := v.Reqeust_api_post("messages.send?", u, data)
	if err != nil {
		return MessagesSendResponseIDs{}, MessagesSendResponseID{}, err
	}

	check := vkerrors.GetError(res)
	if check != "ok" {
		return MessagesSendResponseIDs{}, MessagesSendResponseID{}, errors.New(check)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return MessagesSendResponseIDs{}, MessagesSendResponseID{}, err
	}

	var id MessagesSendResponseID
	var ids MessagesSendResponseIDs

	_ = json.Unmarshal(body, &id)
	_ = json.Unmarshal(body, &ids)

	return ids, id, nil
}

type MessagesSendResponseIDs struct {
	Response []struct {
		PeerID                int `json:"peer_id"`
		MessageID             int `json:"message_id"`
		ConversationMessageID int `json:"conversation_message_id"`
	} `json:"response"`
}

type MessagesSendResponseID struct {
	Response int `json:"response"`
}

type SaveVideoParams struct {
	Name           string
	Description    string
	IsPrivate      bool
	Wallpost       bool
	Link           string
	GroupID        int
	AlbumID        int
	PrivacyView    string
	PrivacyComment string
	NoComments     bool
	Repeat         bool
	//Compression    bool
}

//https://dev.vk.com/method/video.save
func (v VK) SaveVideo(p SaveVideoParams) (SaveVideoResponse, error) {
	data := url.Values{}

	if p.Name != "" {
		data.Set("name", url.QueryEscape(p.Name))
	}

	data.Set("is_private", fmt.Sprint(p.IsPrivate))

	if p.Description != "" {
		data.Set("description", p.Description)
	}

	data.Set("wallpost", fmt.Sprint(p.Wallpost))

	if p.GroupID != 0 {
		data.Set("group_id", fmt.Sprint(p.GroupID))
	}

	if p.Link != "" {
		data.Set("link", p.Link)
	}

	if p.AlbumID != 0 {
		data.Set("album_id", fmt.Sprint(p.AlbumID))
	}

	if p.PrivacyComment != "" {
		data.Set("privacy_comment", p.PrivacyComment)
	}

	if p.PrivacyView != "" {
		data.Set("privacy_view", p.PrivacyView)
	}

	data.Set("no_comments", fmt.Sprint(p.NoComments))

	data.Set("repeat", fmt.Sprint(p.Repeat))

	var u string = ""

	if v.UserToken != "" {
		u += "&v=" + v.Version + "&access_token=" + v.UserToken
	} else {
		return SaveVideoResponse{}, errors.New("Auth token is empty")
	}

	res, err := v.Reqeust_api_get("video.save?", u)
	if err != nil {
		return SaveVideoResponse{}, err
	}

	check := vkerrors.GetError(res)
	if check != "ok" {
		return SaveVideoResponse{}, errors.New(check)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return SaveVideoResponse{}, err
	}

	var r SaveVideoResponse

	err = json.Unmarshal(body, &r)
	if err != nil {
		return SaveVideoResponse{}, err
	}

	return r, nil
}

type SaveVideoResponse struct {
	Response struct {
		AccessKey   string `json:"access_key"`
		Description string `json:"description"`
		OwnerID     int    `json:"owner_id"`
		Title       string `json:"title"`
		UploadURL   string `json:"upload_url"`
		VideoID     int    `json:"video_id"`
	} `json:"response"`
}

//https://dev.vk.com/method/docs.getMessagesUploadServer
func (v VK) GetMessagesUploadServerDoc(doc_type string, peer_id int) (string, error) {
	var u string = ""

	if doc_type != "" {
		u += "type=" + doc_type
	}

	if peer_id != 0 {
		u += "&peer_id=" + fmt.Sprint(peer_id)
	}

	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.Token
	} else {
		return "", errors.New("Auth token is empty")
	}

	res, err := v.Reqeust_api_get("docs.getMessagesUploadServer?", u)
	if err != nil {
		return "", err
	}

	check := vkerrors.GetError(res)
	if check != "ok" {
		return "", errors.New(check)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var g GetMessagesUploadServerDocResponse

	err = json.Unmarshal(body, &g)
	if err != nil {
		return "", err
	}

	return g.Response.UploadURL, nil
}

type GetMessagesUploadServerDocResponse struct {
	Response struct {
		UploadURL string `json:"upload_url"`
	} `json:"response"`
}

/* func (v VK) GetMessagesUploadServerAudio() (string, error) {
	var u string = ""

	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.Token
	} else {
		return "", errors.New("Token group is null")
	}

	res, err := v.Reqeust_api_get("audio.getUploadServer?", u)

	if err != nil {
		return "", err
	}

	check := vkerrors.GetError(res)

	if check != "ok" {
		return "", errors.New(check)
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	var g GetMessagesUploadServerAudioResponse

	err = json.Unmarshal(data, &g)

	if err != nil {
		return "", err
	}

	return g.UploadURL, nil
}
*/

/* type GetMessagesUploadServerAudioResponse struct {
	UploadURL string `json:"upload_url"`
} */

//https://dev.vk.com/method/photos.getMessagesUploadServer
func (v VK) GetMessagesUploadServerPhoto(peer_id int) (string, error) {
	var u string = ""

	if peer_id != 0 {
		u += "peer_id=" + fmt.Sprint(peer_id)
	}

	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.Token
	} else {
		return "", errors.New("Auth token is empty")
	}

	res, err := v.Reqeust_api_get("photos.getMessagesUploadServer?", u)
	if err != nil {
		return "", err
	}

	check := vkerrors.GetError(res)
	if check != "ok" {
		return "", errors.New(check)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var g GetMessagesUploadServerPhotoResponse

	err = json.Unmarshal(data, &g)
	if err != nil {
		return "", err
	}

	return g.Response.UploadURL, nil
}

type GetMessagesUploadServerPhotoResponse struct {
	Response struct {
		AlbumID   int    `json:"album_id"`
		UploadURL string `json:"upload_url"`
		UserID    int    `json:"user_id"`
		GroupID   int    `json:"group_id"`
	} `json:"response"`
}

type SaveDocParams struct {
	File  string
	Title string
	/* 	Tags       string
	   	ReturnTags bool */
}

//https://dev.vk.com/method/docs.save
func (v VK) SaveDoc(p SaveDocParams) (SaveDocResponse, error) {
	data := url.Values{}

	if p.File != "" {
		data.Set("file", fmt.Sprint(p.File))
	} else {
		return SaveDocResponse{}, errors.New("Required field 'File' is empty, MethodName - SaveDoc()")
	}

	/* 	data.Set("return_tags", fmt.Sprint(p.ReturnTags))

	   	if p.Tags != "" {
	   		data.Set("tags", p.Tags)
	   	} */

	if p.Title != "" {
		data.Set("title", p.Title)
	}

	var u string = ""

	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.Token
	} else {
		return SaveDocResponse{}, errors.New("Auth token is empty")
	}

	res, err := v.Reqeust_api_post("docs.save?", u, data)
	if err != nil {
		return SaveDocResponse{}, err
	}

	check := vkerrors.GetError(res)
	if check != "ok" {
		return SaveDocResponse{}, errors.New(check)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return SaveDocResponse{}, err
	}

	var r SaveDocResponse

	err = json.Unmarshal(body, &r)
	if err != nil {
		return SaveDocResponse{}, err
	}

	return r, nil
}

type SaveDocResponse struct {
	Response struct {
		Type string `json:"type"`
		Doc  struct {
			ID      int    `json:"id"`
			OwnerID int    `json:"owner_id"`
			Title   string `json:"title"`
			Size    int    `json:"size"`
			Ext     string `json:"ext"`
			Date    int    `json:"date"`
			Type    int    `json:"type"`
			URL     string `json:"url"`
		} `json:"doc"`
	} `json:"response"`
}

//https://dev.vk.com/method/photos.saveMessagesPhoto
func (v VK) SaveMessagesPhoto(Hash, Photo string, Server int) (SaveMessagesPhotoResponse, error) {
	data := url.Values{}

	if Hash != "" {
		data.Set("hash", Hash)
	}

	if Photo != "" {
		data.Set("photo", Photo)
	}

	if Server != 0 {
		data.Set("server", fmt.Sprint(Server))
	}

	var u string = ""
	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.Token
	} else {
		return SaveMessagesPhotoResponse{}, errors.New("Auth token is empty")
	}

	res, err := v.Reqeust_api_post("photos.saveMessagesPhoto?", u, data)
	if err != nil {
		return SaveMessagesPhotoResponse{}, err
	}

	check := vkerrors.GetError(res)
	if check != "ok" {
		return SaveMessagesPhotoResponse{}, errors.New(check)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return SaveMessagesPhotoResponse{}, err
	}

	var r SaveMessagesPhotoResponse

	err = json.Unmarshal(body, &r)
	if err != nil {
		return SaveMessagesPhotoResponse{}, err
	}

	return r, nil
}

type SaveMessagesPhotoResponse struct {
	Response []SaveMessagesPhoto `json:"response"`
}

type SizesPhoto struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Type   string `json:"type"`
	Width  int    `json:"width"`
}

type SaveMessagesPhoto struct {
	AlbumID   int          `json:"album_id"`
	Date      int          `json:"date"`
	ID        int          `json:"id"`
	OwnerID   int          `json:"owner_id"`
	AccessKey string       `json:"access_key"`
	Sizes     []SizesPhoto `json:"sizes"`
	Text      string       `json:"text"`
	HasTags   bool         `json:"has_tags"`
}

func MediaToAttachment(typ string, OwnerID, MediaID int, access_key string) string {
	if typ != "video" {
		return fmt.Sprintf("%s%d_%d", typ, OwnerID, MediaID)
	} else {
		return fmt.Sprintf("%s%d_%d_%s", typ, OwnerID, MediaID, access_key)
	}
}

func (r SaveMessagesPhoto) GetMaxSizePhotoUrl() SizesPhoto {
	var max = r.Sizes[0]
	for _, v := range r.Sizes {
		if v.Height > max.Height && v.Width > max.Width {
			max = v
		}
	}
	return max
}

type GetVideoParams struct {
	OwnerID  int
	Videos   string
	AlbumID  int
	Count    int
	Offset   int
	Extended bool
	Fields   string
}

//https://dev.vk.com/method/video.get
func (v VK) GetVideo(p GetVideoParams) (GetVideoResponse, error) {
	data := url.Values{}

	if p.AlbumID != 0 {
		data.Set("album_id", fmt.Sprint(p.AlbumID))
	}

	if p.Count != 0 {
		data.Set("count", fmt.Sprint(p.Count))
	}

	data.Set("extended", fmt.Sprint(p.Extended))

	if p.Fields != "" {
		data.Set("fields", p.Fields)
	}

	if p.Offset != 0 {
		data.Set("offset", fmt.Sprint(p.Offset))
	}

	if p.OwnerID != 0 {
		data.Set("owner_id", fmt.Sprint(p.OwnerID))
	}

	if p.Videos != "" {
		data.Set("videos", p.Videos)
	}

	var u string = ""

	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.UserToken
	} else {
		return GetVideoResponse{}, errors.New("Auth token is empty")
	}

	res, err := v.Reqeust_api_post("video.get?", u, data)
	if err != nil {
		return GetVideoResponse{}, err
	}

	check := vkerrors.GetError(res)
	if check != "ok" {
		return GetVideoResponse{}, errors.New(check)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return GetVideoResponse{}, err
	}

	var r GetVideoResponse

	err = json.Unmarshal(body, &r)
	if err != nil {
		return GetVideoResponse{}, err
	}

	return r, nil
}

type GetVideoResponse struct {
	Response struct {
		Count int `json:"count"`
		Items []struct {
			CanComment    int    `json:"can_comment"`
			CanEdit       int    `json:"can_edit"`
			CanLike       int    `json:"can_like"`
			CanRepost     int    `json:"can_repost"`
			CanSubscribe  int    `json:"can_subscribe"`
			CanAddToFaves int    `json:"can_add_to_faves"`
			CanAdd        int    `json:"can_add"`
			CanAttachLink int    `json:"can_attach_link"`
			Comments      int    `json:"comments"`
			Date          int    `json:"date"`
			Description   string `json:"description"`
			Duration      int    `json:"duration"`
			Image         []struct {
				URL         string `json:"url"`
				Width       int    `json:"width"`
				Height      int    `json:"height"`
				WithPadding int    `json:"with_padding"`
			} `json:"image"`
			FirstFrame []struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"first_frame"`
			Width      int    `json:"width"`
			Height     int    `json:"height"`
			ID         int    `json:"id"`
			OwnerID    int    `json:"owner_id"`
			OvID       string `json:"ov_id"`
			Title      string `json:"title"`
			IsFavorite bool   `json:"is_favorite"`
			Player     string `json:"player"`
			Converting int    `json:"converting"`
			Added      int    `json:"added"`
			Type       string `json:"type"`
			Views      int    `json:"views"`
			Likes      struct {
				Count     int `json:"count"`
				UserLikes int `json:"user_likes"`
			} `json:"likes"`
			Reposts struct {
				Count        int `json:"count"`
				WallCount    int `json:"wall_count"`
				MailCount    int `json:"mail_count"`
				UserReposted int `json:"user_reposted"`
			} `json:"reposts"`
		} `json:"items"`
	} `json:"response"`
}
