package lib

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

/*
	Disini isinya function-function untuk helper aja sih,
	Biar gak nulis ulang-ulang kode panjang, kalian bisa
	Tambahin function lain disini juga misal kaya React, Edit, dll.
*/

func GetEphemeralDuration(msg *events.Message) (duration uint32, isEphe bool) {
	if msg == nil || msg.Message == nil {
		return 0, false
	}
	// Disini bisa ditambahin lagi type messsage nya, misalnya GetImageMessage() dll.
	var ctxInfo *waE2E.ContextInfo
	if m := msg.Message.GetExtendedTextMessage(); m != nil {
		ctxInfo = m.GetContextInfo()
	} // Tambahin type message lain disini, buat ngambil ContextInfo nya

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

func Reply(client *whatsmeow.Client, msg *events.Message, text string) (whatsmeow.SendResponse, error) {
	ctxInfo := &waE2E.ContextInfo{
		StanzaID:      &msg.Info.ID,
		Participant:   proto.String(msg.Info.Sender.String()),
		QuotedMessage: msg.Message,
	}

	// Biar gak pentung :p
	if duration, omkeh := GetEphemeralDuration(msg); omkeh {
		ctxInfo.Expiration = &duration
	}
	// bypass participants
	extra := whatsmeow.SendRequestExtra{}
	if msg.Info.Chat.Server == types.GroupServer {
		ownID := client.Store.ID
		if ownID != nil {
			extra.TargetJID = []types.JID{*ownID}
		}
	}
	return client.SendMessage(context.Background(), msg.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:        &text,
			ContextInfo: ctxInfo,
		},
	}, extra)
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
