package gemini

import "go.mau.fi/whatsmeow/proto/waE2E"

func isReply(msg *waE2E.Message) bool {
	if extMsg := msg.ExtendedTextMessage; extMsg != nil && extMsg.ContextInfo != nil {
		if ctxInfo := extMsg.ContextInfo; ctxInfo.Participant != nil && ctxInfo.QuotedMessage != nil {
			return true
		}
	}
	return false
}
