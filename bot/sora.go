package bot

import (
	"context"
	"fmt"
	"reflect"
	"sora/commands"
	"sora/config"
	"sora/lib"
	"sort"
	"strings"
	"sync"
	"time"

	// "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type Bot struct {
	Client *whatsmeow.Client
	Config *config.Configuration
	MapID  sync.Map
}

func NewBot(conf *config.Configuration) *Bot {
	return &Bot{
		Config: conf,
	}
}

func (b *Bot) Start() error {
	go b.MapIDCleaner()
	wa, err := Konek()
	if err != nil {
		return fmt.Errorf("gagal connect ke client WhatsApp : %w", err)
	}
	b.Client = wa
	b.Client.AddEventHandler(b.eventHandler)
	b.Client.SendPresence(context.Background(), types.PresenceAvailable)
	return nil
}

func (b *Bot) Disconnect() {
	if b.Client != nil {
		b.Client.SendPresence(context.Background(), types.PresenceUnavailable)
		b.Client.Log.Infof("Disconnecting bot connection...")
		b.Client.Disconnect()
	}
}

// function untuk bersihin map ID
func (b *Bot) MapIDCleaner() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		b.MapID.Range(func(key, value any) bool {
			timestamp, ok := value.(time.Time)
			if !ok {
				return true
			}
			if now.Sub(timestamp) > 5*time.Minute {
				b.MapID.Delete(key)
			}
			return true
		})
	}
}

func (b *Bot) isOwner(user string) bool {
	for _, ownerNum := range b.Config.Owner {
		if user == ownerNum {
			return true
		}
	}
	return false
}

func (b *Bot) eventHandler(rawEvt any) {
	switch evt := rawEvt.(type) {
	case *events.Message:

		// isReply := false
		// if evt.Message.ExtendedTextMessage.ContextInfo.QuotedMessage != nil {
		// 	isReply = true
		// }
		// b.Client.Log.Warnf("isReply: %t\n", isReply)

		lib.UpdatePushname(evt.Info.Sender.ToNonAD().String(), evt.Info.PushName)
		// jangan proses pesan basi
		if time.Since(evt.Info.Timestamp) > 10*time.Second {
			return
		}
		// b.Client.Log.Infof("LOG « ID: %s", evt.Info.ID)
		msgText, isText := lib.GetText(evt.Message)
		if !isText || msgText == "" { // mastiin isinya itu pesan/text, supaya SenderKeyDistributionMessage, etc. gak masuk ke duplicate check
			go b.logMessageHandler(evt)
			return
		}
		// catat ID msg dan jangan proses pesan duplicate
		_, loaded := b.MapID.LoadOrStore(evt.Info.ID, time.Now())
		if loaded {
			// b.Client.Log.Warnf("CMD « Duplicate message from ID : %s", evt.Info.ID)
			return
		}
		b.commandHandler(evt)
		go b.logMessageHandler(evt)
	}
}

func (b *Bot) commandHandler(evt *events.Message) {
	messageText, ok := lib.GetText(evt.Message)
	if !ok {
		return
	}

	var foundPrefix string
	// sort prefix biar no prefix di check paling akir
	sort.Slice(b.Config.Prefix, func(i, j int) bool {
		return len(b.Config.Prefix[i]) > len(b.Config.Prefix[j])
	})
	for _, prefix := range b.Config.Prefix {
		if strings.HasPrefix(messageText, prefix) {
			foundPrefix = prefix
			break
		}
	}

	trimmedText := strings.TrimPrefix(messageText, foundPrefix)
	parts := strings.Fields(trimmedText)
	if len(parts) == 0 {
		return
	}

	commandName := strings.ToLower(parts[0])
	args := parts[1:]
	rawArgs := strings.Join(args, " ")

	cmd, found := commands.Commands[commandName]
	if !found {
		return
	}
	sender := evt.Info.Sender.ToNonAD()
	// b.Client.Log.Warnf("AddressingMode : %s", evt.Info.Sender.Server)
	if evt.Info.Sender.Server == "lid" {
		var err error
		sender, err = b.Client.Store.LIDs.GetPNForLID(context.Background(), evt.Info.Sender)
		if err != nil {
			b.Client.Log.Errorf("CMD « Sum error happen bruh : ", err)
			return
		}
	}
	senderNum := sender.User
	if b.Config.Mode == "self" && !evt.Info.IsFromMe && !b.isOwner(senderNum) {
		return
	}
	if b.Config.Mode == "public" && !evt.Info.IsFromMe && !b.isOwner(senderNum) {
		if foundPrefix == "" {
			return
		}
	}

	if cmd.Permission > commands.Public {
		if cmd.Permission == commands.Owner {
			if !b.isOwner(senderNum) {
				b.Client.Log.Warnf("CMD « Permission denied for user '%s'", senderNum)
				return
			}
			b.Client.Log.Infof("CMD « Permission granted for user '%s'", senderNum)
		}
	}

	ctx := &commands.CommandContext{
		Ctx:     context.Background(),
		Client:  b.Client,
		Message: evt,
		Args:    args,
		RawArgs: rawArgs,
		Conf:    b.Config,
	}
	go func() {
		b.Client.SendChatPresence(context.Background(), evt.Info.Chat, types.ChatPresenceComposing, types.ChatPresenceMediaText)
		defer b.Client.SendChatPresence(context.Background(), evt.Info.Chat, types.ChatPresencePaused, types.ChatPresenceMediaText)
		cmd.Exec(ctx)
	}()

	var sourceName string
	if evt.Info.IsGroup {
		groupInfo, _ := b.Client.GetGroupInfo(ctx.Ctx, evt.Info.Chat)
		if groupInfo != nil {
			sourceName = fmt.Sprintf("in Group '%s'", groupInfo.Name)
		} else {
			sourceName = fmt.Sprintf("in Group '%s'", evt.Info.Chat.String())
		}
	} else {
		sourceName = "in Private"
	}
	b.Client.Log.Infof("CMD « Executing '%s' for %s %s", commandName, evt.Info.PushName, sourceName)
}

func (b *Bot) logMessageHandler(evt *events.Message) {
	var source string
	if evt.Info.Chat.String() == "status@broadcast" {
		source = fmt.Sprintf("Status from '%s'", evt.Info.PushName)
	} else if evt.Info.IsGroup {
		groupInfo, err := b.Client.GetGroupInfo(context.Background(), evt.Info.Chat)
		if err != nil {
			return
		}
		if groupInfo != nil {
			source = fmt.Sprintf("Group '%s'", groupInfo.Name)
		} else {
			source = fmt.Sprintf("Group '%s'", evt.Info.Chat.String())
		}
	} else if evt.Info.Chat.Server == "newsletter" {
		newsletterIngfo, err := b.Client.GetNewsletterInfo(context.Background(), evt.Info.Chat)
		if err != nil {
			return
		}
		if newsletterIngfo != nil && newsletterIngfo.ThreadMeta.Name.Text != "" {
			source = fmt.Sprintf("Channel '%s'", newsletterIngfo.ThreadMeta.Name.Text)
		} else {
			// fallback
			source = fmt.Sprintf("Channel '%s'", evt.Info.Chat.String())
		}
	} else {
		source = "Private"
	}

	msgType := "(Unknown)"
	if evt.Message.GetConversation() != "" {
		msgType = "Conversation"
	} else {
		v := reflect.ValueOf(evt.Message)
		if v.IsValid() && !v.IsNil() {
			v = v.Elem()
			for i := 0; i < v.NumField(); i++ {
				field := v.Field(i)
				if field.Kind() == reflect.Ptr && !field.IsNil() {
					msgType = reflect.TypeOf(field.Interface()).Elem().Name()
					break
				}
			}
		}
	}
	content := "(Nothing)"
	if text, ok := lib.GetText(evt.Message); ok && text != "" {
		content = fmt.Sprintf(`"%s"`, text)
	} else if react := evt.Message.GetReactionMessage(); react != nil {
		content = react.GetText()
	}

	msgTime := evt.Info.Timestamp.Format("15:04:05")

	b.Client.Log.Infof("MSG « From: %s (%s) « In: %s", evt.Info.PushName, evt.Info.Sender.String(), source)
	b.Client.Log.Infof("    └── %s | Type: %s | Text: %s", msgTime, msgType, content)
}
