package gemini

import (
	"context"
	"fmt"
	"log"
	cmd "sora/commands"
	"sora/lib"
	"time"

	"go.mau.fi/whatsmeow/types"
	"google.golang.org/genai"
)

func init() {
	cmd.Plugin(cmd.Cmd{
		Name:     "gemini",
		Alias:    []string{"gm"},
		Desc:     "Gemini E Ay",
		Exec:     gemini,
		Category: "ai",
	})
}

func InitGemini() {
	// api key udah otomatis ambil dari env
	client, err := genai.NewClient(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	GenaiClient = client
	log.Println("Gemini Initialized")
}

func gemini(ctx *cmd.CommandContext) {
	var nol int32
	jakarta, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(jakarta)
	timestamp := now.Format("02/01/2006, 15:04")
	senderName := ctx.Message.Info.PushName
	text := ctx.RawArgs
	chatID := ctx.Message.Info.Chat.String()
	config := &genai.GenerateContentConfig{
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingBudget: &nol,
		},
		SystemInstruction: genai.NewContentFromText("Kamu adalah seorang E Ay yang santai, namamu gemini, jawab teks pengguna dengan sesingkat dan sesantai mungkin seperti sedang chattingan di whatsapp. nanti pengguna bawa metadata kaya nama, timestamp, dll. Lu jawabnya nanti gausah ikut ikutan pake format/metadata ya", genai.RoleUser),
	}
	session, ok := ActiveSessions[chatID]
	if !ok {
		var err error
		history := LoadSession(chatID)
		session, err = GenaiClient.Chats.Create(context.Background(), "gemini-2.5-flash", config, history)
		if err != nil {
			ctx.Client.Log.Errorf("Error creating chat: %v", err)
			return
		}
		ActiveSessions[chatID] = session
	}
	if text == "" {
		text = "halo gemini"
	}
	finalMsg := fmt.Sprintf("[%s | %s] Pesan: %s", senderName, timestamp, text)

	if isReply(ctx.Message.Message) == true {
		ctxInfo := ctx.Message.Message.ExtendedTextMessage.ContextInfo
		repliedJID := *ctxInfo.Participant
		repliedText, _ := lib.GetText(ctxInfo.QuotedMessage)
		repliedName := "Gemini"
		ownJID := ctx.Client.Store.ID.ToNonAD().String()
		if repliedJID != ownJID {
			repliedName = lib.GetCachedName(repliedJID)
		}
		finalMsg = fmt.Sprintf(`[%s | %s] (Membalas ke "%s: %s") Dengan pesan: %s`, senderName, timestamp, repliedName, repliedText, text)
	}
	// ctx.Client.Log.Warnf("Final Message: %s", finalMsg)
	ctx.Client.SendChatPresence(ctx.Ctx, ctx.Message.Info.Chat, types.ChatPresenceComposing, types.ChatPresenceMediaText)
	result, err := session.SendMessage(context.Background(), genai.Part{Text: finalMsg})
	if err != nil {
		ctx.Client.Log.Errorf("Error generating content: %v", err)
		return
	}
	SaveSession(chatID, session.History(true))
	ctx.Reply(result.Text())
	ctx.Client.SendChatPresence(ctx.Ctx, ctx.Message.Info.Chat, types.ChatPresencePaused, types.ChatPresenceMediaText)
	// ctx.Client.Log.Warnf("Respon Gemini: %s", result.Text())
}
