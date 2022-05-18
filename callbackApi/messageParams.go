package callbackApi

type PhotoSize struct {
	Type   string `json:"type"`
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type AttachmentType string

const (
	Photo        AttachmentType = "photo"
	Video        AttachmentType = "video"
	Audio        AttachmentType = "audio"
	Doc          AttachmentType = "doc"
	Sticker      AttachmentType = "sticker"
	Graffiti     AttachmentType = "graffiti"
	AudioMessage AttachmentType = "audio_message"
)

type PhotoAttachment struct {
	ID      int         `json:"id"`
	AlbumID int         `json:"album_id"`
	OwnerID int         `json:"owner_id"`
	UserID  int         `json:"user_id"`
	Text    string      `json:"text"`
	Date    int         `json:"date"`
	Sizes   []PhotoSize `json:"sizes"`
	Width   int         `json:"width"`
	Height  int         `json:"height"`
}

type VideoAttachment struct {
	AccessKey     string `json:"access_key"`
	CanEdit       int    `json:"can_edit"`
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
	Title      string `json:"title"`
	IsFavorite bool   `json:"is_favorite"`
	TrackCode  string `json:"track_code"`
	Type       string `json:"type"`
	Views      int    `json:"views"`
}

type AudioAttachment struct {
	Artist       string `json:"artist"`
	ID           int    `json:"id"`
	OwnerID      int    `json:"owner_id"`
	Title        string `json:"title"`
	Duration     int    `json:"duration"`
	IsExplicit   bool   `json:"is_explicit"`
	IsFocusTrack bool   `json:"is_focus_track"`
	TrackCode    string `json:"track_code"`
	URL          string `json:"url"`
	Date         int    `json:"date"`
	MainArtists  []struct {
		Name   string `json:"name"`
		Domain string `json:"domain"`
		ID     string `json:"id"`
	} `json:"main_artists"`
	FeaturedArtists []struct {
		Name   string `json:"name"`
		Domain string `json:"domain"`
		ID     string `json:"id"`
	} `json:"featured_artists"`
	ShortVideosAllowed  bool `json:"short_videos_allowed"`
	StoriesAllowed      bool `json:"stories_allowed"`
	StoriesCoverAllowed bool `json:"stories_cover_allowed"`
}

type DocAttachment struct {
	ID        int    `json:"id"`
	OwnerID   int    `json:"owner_id"`
	Title     string `json:"title"`
	Size      int    `json:"size"`
	Ext       string `json:"ext"`
	Date      int    `json:"date"`
	Type      int    `json:"type"`
	URL       string `json:"url"`
	AccessKey string `json:"access_key"`
}

type GraffitiAttachment struct {
	ID        int    `json:"id"`
	OwnerID   int    `json:"owner_id"`
	URL       string `json:"url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	AccessKey string `json:"access_key"`
}

type StickerAttachment struct {
	StickerID int `json:"sticker_id"`
	ProductID int `json:"product_id"`
	Images    []struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"images"`
	ImagesWithBackground []struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"images_with_background"`
	AnimationURL string `json:"animation_url"`
	IsAllowed    bool   `json:"is_allowed"`
}

type AudioMessageAttachment struct {
	Duration  int    `json:"duration"`
	ID        int    `json:"id"`
	LinkMp3   string `json:"link_mp3"`
	LinkOgg   string `json:"link_ogg"`
	OwnerID   int    `json:"owner_id"`
	AccessKey string `json:"access_key"`
	Waveform  []int  `json:"waveform"`
}

type Attachments struct {
	Type         string                 `json:"type"`
	Photo        PhotoAttachment        `json:"photo"`
	Video        VideoAttachment        `json:"video"`
	Audio        AudioAttachment        `json:"audio"`
	Doc          DocAttachment          `json:"doc"`
	Graffiti     GraffitiAttachment     `json:"graffiti"`
	Sticker      StickerAttachment      `json:"sticker"`
	AudioMessage AudioMessageAttachment `json:"audio_message"`
}

type MessageObjectMessage struct {
	Date                  int                     `json:"date"`
	FromID                int                     `json:"from_id"`
	ID                    int                     `json:"id"`
	Out                   int                     `json:"out"`
	Attachments           []*Attachments          `json:"attachments"`
	ConversationMessageID int                     `json:"conversation_message_id"`
	FwdMessages           []*MessageObjectMessage `json:"fwd_messages"`
	Geo                   *struct {
		Coordinates struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"coordinates"`
		Place struct {
			City    string `json:"city"`
			Country string `json:"country"`
			Title   string `json:"title"`
		} `json:"place"`
		Type string `json:"type"`
	} `json:"geo,omitempty"`
	ReplyMessage *MessageObjectMessage `json:"reply_message"`
	Important    bool                  `json:"important"`
	IsHidden     bool                  `json:"is_hidden"`
	PeerID       int                   `json:"peer_id"`
	RandomID     int                   `json:"random_id"`
	Text         string                `json:"text"`
}

type MessageAllowObject struct {
	UserID int    `json:"user_id"`
	Key    string `json:"key"`
}

type MessageDenyObject struct {
	UserID int `json:"user_id"`
}

type MessageTypingStateObject struct {
	State  string `json:"state"`
	FromID int    `json:"from_id"`
	ToID   int    `json:"to_id"`
}

type MessageObjectClientInfo struct {
	ButtonsActions []string `json:"button_actions"`
	Keyboard       bool     `json:"keyboard"`
	InlineKeyboard bool     `json:"inline_keyboard"`
	Carousel       bool     `json:"carousel"`
	LangID         int      `json:"lang_id"`
}

type MessageObject struct {
	Message    MessageObjectMessage    `json:"message"`
	ClientInfo MessageObjectClientInfo `json:"client_info"`
}

type Message struct {
	Type    string        `json:"type"`
	Object  MessageObject `json:"object"`
	GroupID int           `json:"group_id"`
	Secret  string        `json:"secret"`
}

func (p PhotoAttachment) GetMaxSizePhotoUrl() PhotoSize {
	var max = p.Sizes[0]
	for _, v := range p.Sizes {
		if v.Height > max.Height && v.Width > max.Width {
			max = v
		}
	}
	return max
}
