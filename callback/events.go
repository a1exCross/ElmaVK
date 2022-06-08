package callback

import (
	"encoding/json"
	"log"
)

type EventType string

const (
	MessageNew         EventType = "message_new"
	MessageReply       EventType = "message_reply"
	MessageEdit        EventType = "message_edit"
	MessageAllow       EventType = "message_allow"
	MessageDeny        EventType = "message_deny"
	MessageTypingState EventType = "message_typing_state"
	MessageEvent       EventType = "message_event"
)

type Events struct {
	Type    EventType       `json:"type"`
	Object  json.RawMessage `json:"object"`
	GroupID int             `json:"group_id"`
	EventID string          `json:"event_id"`
	Secret  string          `json:"secret"`
}

type FuncList struct {
	NewMessage         func(e Events, obj MessageObject)
	MessageReply       func(e Events, obj MessageObjectMessage)
	MessageEdit        func(e Events, obj MessageObjectMessage)
	MessageTypingState func(e Events, obj MessageTypingStateObject)
	MessageAllow       func(e Events, obj MessageAllowObject)
	MessageDeny        func(e Events, obj MessageDenyObject)
	MessageEvent       func(e Events, obj MessageEventObject)
}

func (c Callback) CallFuncList(data []byte, e Events) {
	switch e.Type {
	case MessageNew:
		{
			var obj MessageObject

			err := json.Unmarshal(e.Object, &obj)
			if err != nil {
				log.Println(err)
			}

			var p *Payload

			err = json.Unmarshal(e.Object, &p)
			if err != nil {
				log.Println(err)
			}

			obj.Message.Payload = p

			go c.Functions.NewMessage(e, obj)
		}
	case MessageReply:
		{
			var obj MessageObjectMessage

			err := json.Unmarshal(e.Object, &obj)
			if err != nil {
				log.Println(err)
			}

			var p *Payload

			err = json.Unmarshal(e.Object, &p)
			if err != nil {
				log.Println(err)
			}

			obj.Payload = p

			go c.Functions.MessageReply(e, obj)
		}
	case MessageEdit:
		{
			var obj MessageObjectMessage

			err := json.Unmarshal(e.Object, &obj)
			if err != nil {
				log.Println(err)
			}

			var p *Payload

			err = json.Unmarshal(e.Object, &p)
			if err != nil {
				log.Println(err)
			}

			obj.Payload = p

			go c.Functions.MessageEdit(e, obj)
		}
	case MessageAllow:
		{
			var obj MessageAllowObject

			err := json.Unmarshal(e.Object, &obj)
			if err != nil {
				log.Println(err)
			}

			go c.Functions.MessageAllow(e, obj)
		}
	case MessageDeny:
		{
			var obj MessageDenyObject

			err := json.Unmarshal(e.Object, &obj)
			if err != nil {
				log.Println(err)
			}

			go c.Functions.MessageDeny(e, obj)
		}
	case MessageTypingState:
		{
			var obj MessageTypingStateObject

			err := json.Unmarshal(e.Object, &obj)
			if err != nil {
				log.Println(err)
			}

			go c.Functions.MessageTypingState(e, obj)
		}
	case MessageEvent:
		{
			var obj MessageEventObject

			err := json.Unmarshal(e.Object, &obj)
			if err != nil {
				log.Println(err)
			}

			var p Payload

			err = json.Unmarshal(e.Object, &p)
			if err != nil {
				log.Println(err)
			}

			obj.Payload = &p

			go c.Functions.MessageEvent(e, obj)
		}
	}
}
