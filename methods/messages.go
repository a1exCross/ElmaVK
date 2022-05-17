package Methods

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
	"strconv"
	"strings"

	ApiErrors "github.com/a1exCross/ElmaVK/errors"
)

type SendMessage struct {
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
	GroupID   int
	Token     string
	//Keyboard text
	//Template text
	//Payload text
	//ContentSource text json
	DontParseLink   bool
	DisableMentions bool
	Intent          string
	SubscribeID     int
}

type ResponseSendMessage struct {
	Response int `json:"response"`
}

func (v VK) GetSendMessageParam() *SendMessage {
	return &SendMessage{}
}

type GetAttachmentsParams struct {
	FilePaths []string
	PeerID    int
}

func (v VK) GetAttachments(p GetAttachmentsParams) ([]string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	var at []string

	for _, path := range p.FilePaths {
		file, err := os.Open(path)

		if err != nil {
			log.Println(err)
		}

		defer file.Close()

		filename := filepath.Base(path)
		log.Println(filename)

		ext := filepath.Ext(filename)
		fmt.Println(ext)

		title := strings.Replace(filename, ext, "", len(ext))

		log.Println(title)

		ext = strings.ToLower(ext)

		up := ""
		field := ""

		if ext == ".jpg" || ext == ".png" || ext == ".gif" {
			up, err = v.GetMessagesUploadServerPhoto(p.PeerID)

			if err != nil {
				return nil, err
			}
			field = "photo"
		}

		if ext == ".docx" {
			up, err = v.GetMessagesUploadServerDoc("doc", p.PeerID)

			if err != nil {
				return nil, err
			}

			field = "doc"
			/* up, err = v.GetMessagesUploadServerAudio()

			if err != nil {
				return nil, err
			} */
		}

		var r *SaveVideoResponse

		if ext == ".mp4" {
			r, err = v.SaveVideo(SaveVideoParam{
				Name:      title, //filename,
				IsPrivate: true,
				Wallpost:  false,
				//GroupID:   203374987,
			})

			up = r.Response.UploadURL

			if err != nil {
				return nil, err
			}

			field = "video"
		}

		part, _ := writer.CreateFormFile("file", filepath.Base(file.Name())) // или file

		io.Copy(part, file)

		//log.Println(writer.FormDataContentType())

		err = writer.Close()

		if err != nil {
			return nil, err
		}

		res, err := http.DefaultClient.Post(up, writer.FormDataContentType(), body)

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

		if field == "photo" {

			var r ResponseAfterUploadPhoto

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

		if field == "doc" {
			var r ResponseAfterUploadDoc

			err = json.Unmarshal(data, &r)

			if err != nil {
				return nil, err
			}

			s, err := v.SaveDoc(r.File)

			if err != nil {
				return nil, err
			}

			at = append(at, MediaToAttachment("doc", s.Response.Doc.OwnerID, s.Response.Doc.ID, ""))
		}

		if field == "video" {
			at = append(at, MediaToAttachment("video", r.Response.OwnerID, r.Response.VideoID, r.Response.AccessKey))
			//at = append(at, "video"+strconv.Itoa(r.Response.OwnerID)+"_"+strconv.Itoa(r.Response.VideoID)+"_"+r.Response.AccessKey)
		}

		//log.Println(s.Response)

		/* at := PhotoToAttachment(s.Response[0].OwnerID, s.Response[0].ID) */

		//for _, r := range s.Response {

		//}
	}

	return at, nil
}

//https://dev.vk.com/method/messages.send
func (v *VK) SendMessage(s SendMessage) (*ResponseSendMessage, error) {
	var u string = ""

	data := url.Values{}

	if s.UserID != 0 {
		//u += "user_id=" + strconv.Itoa(s.UserID)
		data.Set("user_id", strconv.Itoa(s.UserID))
	} /* else {
		return nil, errors.New("User ID is null")
	} */

	if len(s.UserIDs) > 0 {
		var ids = ""
		for i := 0; i < len(s.UserIDs); i++ {
			if i > 0 {
				ids += ","
			}
			ids += strconv.Itoa(s.UserIDs[i])
		}
		data.Set("user_ids", ids)
	}

	if s.PeerID != 0 {
		//u += "user_id=" + strconv.Itoa(s.UserID)
		data.Set("peer_id", strconv.Itoa(s.PeerID))
	} /* else {
		return nil, errors.New("User ID is null")
	} */

	if len(s.PeerIDs) > 0 {
		var ids = ""
		for i := 0; i < len(s.PeerIDs); i++ {
			if i > 0 {
				ids += ","
			}
			ids += strconv.Itoa(s.PeerIDs[i])
		}
		data.Set("peer_ids", ids)
	}

	if s.Domain != "" {
		//u += "&message=" + url.QueryEscape(s.Message)
		data.Set("domain", s.Domain)
	}

	if s.ChatID != 0 {
		data.Set("chat_id", strconv.Itoa(s.ChatID))
	}

	if s.Message != "" {
		//u += "&message=" + url.QueryEscape(s.Message)
		data.Set("message", url.QueryEscape(s.Message))
	} /* else {
		return nil, errors.New("Message text is null")
	} */

	if s.Guid != 0 {
		data.Set("guid", strconv.Itoa(s.Guid))
	}

	if len(s.Attachment) > 0 {
		var att = ""
		for i := 0; i < len(s.Attachment); i++ {
			if i > 0 {
				att += ","
			}
			att += s.Attachment[i]
		}
		data.Set("attachment", att)
	}

	data.Set("random_id", strconv.Itoa(s.RandomID))

	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.Token
	} /* else {
		return nil, errors.New("Token group is null")
	} */

	res, err := v.Reqeust_api_post("messages.send?", u, data)

	if err != nil {
		return nil, err
	}

	check := ApiErrors.GetError(res)

	if check != "ok" {
		return nil, errors.New(check)
	}

	dat, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var t ResponseSendMessage

	err = json.Unmarshal(dat, &t)

	if err != nil {
		return nil, err
	}

	log.Println(t)

	return &t, nil
}

type SaveVideoParam struct {
	Name          string
	Description   string
	IsPrivate     bool
	Wallpost      bool
	Link          string
	GroupID       int
	AlbumID       int
	PrivacyView   string
	PrivacComment string
	NoComments    bool
	Repeat        bool
	Compression   bool
}

func (v VK) SaveVideo(p SaveVideoParam) (*SaveVideoResponse, error) {
	var u string = ""

	if p.Name != "" {
		u += "name=" + p.Name //пробелы убрать
	}

	u += "&is_private=1" //+ strconv.FormatBool(p.IsPrivate)

	//u += "&description=ab"

	u += "&wallpost=0" //+ strconv.FormatBool(p.Wallpost)

	if p.GroupID != 0 {
		u += "&group_id=" + strconv.Itoa(p.GroupID)
	}

	/* if f != "" {
		u += "file=" + url.QueryEscape(f)
	} */

	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.UserToken
	}

	/* else {
		return nil, errors.New("Token group is null")
	} */

	res, err := v.Reqeust_api_get("video.save?", u)

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

	/* 	log.Println(string(data))

	   	os.Exit(1) */

	var r SaveVideoResponse

	err = json.Unmarshal(data, &r)

	if err != nil {
		return nil, err
	}

	return &r, nil
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

func (v VK) GetMessagesUploadServerDoc(tp string, peer_id int) (string, error) {
	var u string = ""

	if tp != "" {
		u += "type=" + tp
	}

	if peer_id != 0 {
		u += "&peer_id=" + strconv.Itoa(peer_id)
	}

	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.Token
	} /*  else {
		return nil, errors.New("Token group is null")
	} */

	res, err := v.Reqeust_api_get("docs.getMessagesUploadServer?", u)

	if err != nil {
		return "", err
	}

	check := ApiErrors.GetError(res)

	if check != "ok" {
		return "", errors.New(check)
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	var g GetMessagesUploadServerDocResponse

	err = json.Unmarshal(data, &g)

	if err != nil {
		return "", err
	}

	return g.Response.UploadURL, nil
}

type ResponseAfterUploadDoc struct {
	File string `json:"file"`
}

func (v VK) GetMessagesUploadServerAudio() (string, error) {
	var u string = ""

	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.Token
	} /*  else {
		return nil, errors.New("Token group is null")
	} */

	res, err := v.Reqeust_api_get("audio.getUploadServer?", u)

	if err != nil {
		return "", err
	}

	check := ApiErrors.GetError(res)

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

type GetMessagesUploadServerAudioResponse struct {
	UploadURL string `json:"upload_url"`
}

type GetMessagesUploadServerDocResponse struct {
	Response struct {
		UploadURL string `json:"upload_url"`
	} `json:"response"`
}

type GetMessagesUploadServerPhotoResponse struct {
	AlbumID   int    `json:"album_id"`
	UploadURL string `json:"upload_url"`
	UserID    int    `json:"user_id"`
	GroupID   int    `json:"group_id"`
}

type GetMessagesUploadServerPhoto struct {
	GetMessagesUploadServerPhotoResponse `json:"response"`
}

func (v VK) GetMessagesUploadServerPhoto(peer_id int) (string, error) {
	var u string = ""

	if peer_id != 0 {
		u += "peer_id=" + strconv.Itoa(peer_id)
	}

	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.Token
	} /*  else {
		return nil, errors.New("Token group is null")
	} */

	res, err := v.Reqeust_api_get("photos.getMessagesUploadServer?", u)

	if err != nil {
		return "", err
	}

	check := ApiErrors.GetError(res)

	if check != "ok" {
		return "", errors.New(check)
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	var g GetMessagesUploadServerPhoto

	err = json.Unmarshal(data, &g)

	if err != nil {
		return "", err
	}

	return g.UploadURL, nil
}

type ResponseAfterUploadPhoto struct {
	Server int    `json:"server"`
	Photo  string `json:"photo"`
	Hash   string `json:"hash"`
}

type SizesPhoto struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Type   string `json:"type"`
	Width  int    `json:"width"`
}

type ResponseSaveMessagesPhoto struct {
	AlbumID   int          `json:"album_id"`
	Date      int          `json:"date"`
	ID        int          `json:"id"`
	OwnerID   int          `json:"owner_id"`
	AccessKey string       `json:"access_key"`
	Sizes     []SizesPhoto `json:"sizes"`
	Text      string       `json:"text"`
	HasTags   bool         `json:"has_tags"`
}

func (v VK) SaveDoc(f string) (*SaveDoc, error) {
	var u string = ""

	/* if Hash != "" {
		u += "hash=" + Hash
	}

	if Photo != "" {
		u += "&photo=" + Photo
	}

	if Server != 0 {
		u += "&server=" + strconv.Itoa(Server)
	} */

	if f != "" {
		u += "file=" + url.QueryEscape(f)
	}

	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.Token
	}

	/* else {
		return nil, errors.New("Token group is null")
	} */

	res, err := v.Reqeust_api_get("docs.save?", u)

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

	var p SaveDoc

	err = json.Unmarshal(data, &p)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

type SaveDoc struct {
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

type SaveMessagesPhoto struct {
	Response []ResponseSaveMessagesPhoto `json:"response"`
}

func (v VK) SaveMessagesPhoto(Hash, Photo string, Server int) (*SaveMessagesPhoto, error) {
	var u string = ""

	if Hash != "" {
		u += "hash=" + Hash
	}

	if Photo != "" {
		u += "&photo=" + Photo
	}

	if Server != 0 {
		u += "&server=" + strconv.Itoa(Server)
	}

	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.Token
	}

	/* else {
		return nil, errors.New("Token group is null")
	} */

	res, err := v.Reqeust_api_get("photos.saveMessagesPhoto?", u)

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

	var p SaveMessagesPhoto

	err = json.Unmarshal(data, &p)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func MediaToAttachment(typ string, OwnerID, MediaID int, access_key string) string {
	if typ != "video" {
		return fmt.Sprintf("%s%d_%d", typ, OwnerID, MediaID)
	} else {
		return fmt.Sprintf("%s%d_%d_%s", typ, OwnerID, MediaID, access_key)
	}
}

func (r ResponseSaveMessagesPhoto) GetMaxSizePhotoUrl() *SizesPhoto {
	var max = r.Sizes[0]
	for _, v := range r.Sizes {
		if v.Height > max.Height && v.Width > max.Width {
			max = v
		}
	}
	return &max
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

func (v VK) GetVideo(p GetVideoParams) (*GetVideoResponse, error) {
	var u string = ""

	/* if p. != "" {
		u += "hash=" + Hash
	} */

	/* 	if Photo != "" {
	   		u += "&photo=" + Photo
	   	}

	   	if Server != 0 {
	   		u += "&server=" + strconv.Itoa(Server)
	   	} */

	u += "owner_id=" + strconv.Itoa(p.OwnerID) + "&videos=" + p.Videos

	if v.Token != "" {
		u += "&v=" + v.Version + "&access_token=" + v.UserToken
	}

	res, err := v.Reqeust_api_get("video.get?", u)

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

	//log.Println(string(data))

	var r GetVideoResponse

	err = json.Unmarshal(data, &r)

	if err != nil {
		return nil, err
	}

	return &r, nil
}
