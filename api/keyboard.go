package api

import "fmt"

type Keyboard struct {
	OneTime bool                `json:"one_time"`
	Buttons [][]KeyboardButtons `json:"buttons"`
	Inline  bool                `json:"inline"`
}

type ButtonColor string

const (
	Primary   ButtonColor = "primary"
	Secondary ButtonColor = "secondary"
	Negative  ButtonColor = "negative"
	Positive  ButtonColor = "positive"
)

type ActionType string

const (
	ActionText     ActionType = "text"
	ActionOpenLink ActionType = "open_link"
	ActionLocation ActionType = "location"
	ActionVkPay    ActionType = "vkpay"
	ActionVkApps   ActionType = "open_app"
	ActionCallback ActionType = "callback"
)

type KeyboardActionTypeText struct {
	Type    ActionType `json:"type"`
	Payload string     `json:"payload,omitempty"`
	Label   string     `json:"label"`
}

type KeyboardActionTypeOpenLink struct {
	Type    ActionType `json:"type"`
	Payload string     `json:"payload,omitempty"`
	Label   string     `json:"label"`
	Link    string     `json:"link"`
}

type KeyboardActionTypeVkApps struct {
	Type    ActionType `json:"type"`
	Payload string     `json:"payload,omitempty"`
	Label   string     `json:"label"`
	AppID   int        `json:"app_id"`
	OwnerID int        `json:"owner_id"`
	Hash    string     `json:"hash"`
}

type KeyboardActionTypeVkPay struct {
	Type    ActionType `json:"type"`
	Payload string     `json:"payload,omitempty"`
	Hash    string     `json:"hash"`
}

type KeyboardActionTypeLocation struct {
	Type    ActionType `json:"type"`
	Payload string     `json:"payload,omitempty"`
}

type KeyboardActionTypeCallback struct {
	Type    ActionType `json:"type"`
	Payload string     `json:"payload,omitempty"`
	Label   string     `json:"label"`
}

type Action struct {
	Text     KeyboardActionTypeText
	OpenLink KeyboardActionTypeOpenLink
	Location KeyboardActionTypeLocation
	VkPay    KeyboardActionTypeVkPay
	VkApps   KeyboardActionTypeVkApps
	Callback KeyboardActionTypeCallback
}

type KeyboardButtons struct {
	Action interface{} `json:"action"`
	Color  ButtonColor `json:"color,omitempty"`
}

func (b *Keyboard) AddLine() {
	var p []KeyboardButtons
	b.Buttons = append(b.Buttons, p)
}

func ToPayload(s string) string {
	return fmt.Sprintf("{\"button\": \"%s\"}", s)
}

func (b *Keyboard) AddButton(p KeyboardButtons) {
	str_count := 0

	if len(b.Buttons) > 0 {
		str_count = len(b.Buttons)
		b.Buttons[str_count-1] = append(b.Buttons[str_count-1], p)
	}
}

func GetKeyboard() Keyboard {
	return Keyboard{
		Buttons: make([][]KeyboardButtons, 1),
	}
}
