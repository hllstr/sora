package lib

import (
	"context"
	//"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)


func GetEphemeralDuration(msg *events.Message) (duration uint32, isEphe bool) {
	if msg == nil || msg.Message == nil {
		return 0, false
	}
	// Disini bisa ditambahin lagi type messsage nya, misalnya GetImageMessage() dll.
	var ctxInfo *waE2E.ContextInfo
	if m := msg.Message.GetExtendedTextMessage(); m != nil {
		ctxInfo = m.GetContextInfo()
	} else if m := msg.Message.GetImageMessage(); m != nil {
		ctxInfo = m.GetContextInfo()
	} else if m := msg.Message.GetVideoMessage(); m != nil {
		ctxInfo = m.GetContextInfo()
	} else if msg.Message.GetDocumentMessage() != nil {
		ctxInfo = m.GetContextInfo()
	}

	if ctxInfo != nil {
		if exp := ctxInfo.GetExpiration(); exp > 0 {
			return exp, true
		}
	}

	if msgCtx := msg.Message.GetMessageContextInfo(); msgCtx != nil {
		if exp := msgCtx.GetMessageAddOnDurationInSecs(); exp > 0 {
			return exp, true
		}
	}
	return 0, false
}

// bypass participant function
func Bypass(client *whatsmeow.Client, chatJID types.JID) whatsmeow.SendRequestExtra {
	extra := whatsmeow.SendRequestExtra{}
	if chatJID.Server == types.GroupServer {
		ownID := client.Store.ID
		if ownID != nil {
			extra.TargetJID = []types.JID{*ownID}
		}
	}
	return extra
}

func Reply(client *whatsmeow.Client, msg *events.Message, text string) (whatsmeow.SendResponse, error) {
	ctxInfo := &waE2E.ContextInfo{
		StanzaID:      &msg.Info.ID,
		Participant:   proto.String(msg.Info.Sender.String()),
		QuotedMessage: msg.Message,
	}

	if duration, omkeh := GetEphemeralDuration(msg); omkeh {
		ctxInfo.Expiration = &duration
	}

	bypass := Bypass(client, msg.Info.Chat)
	return client.SendMessage(context.Background(), msg.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:        &text,
			ContextInfo: ctxInfo,
		},
	}, bypass) 
}

func GetText(msg *events.Message) (text string, ok bool) {
	if msg == nil || msg.Message == nil {
		return "", false
	}
	if m := msg.Message.GetConversation(); m != "" {
		return m, true
	} else if m := msg.Message.GetExtendedTextMessage(); m != nil {
		return m.GetText(), true
	} else if m := msg.Message.GetImageMessage(); m != nil {
		return m.GetCaption(), true
	} else if m := msg.Message.GetVideoMessage(); m != nil {
		return m.GetCaption(), true
	} else if m := msg.Message.GetDocumentMessage(); m != nil {
		return m.GetCaption(), true
	}

	return "", false
}