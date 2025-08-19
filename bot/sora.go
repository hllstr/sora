package bot

import (
	"context"
	"fmt"
	"reflect"
	"sora/commands"
	"sora/config"
	"sora/lib"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

type Bot struct {
	Client *whatsmeow.Client
	Config *config.Configuration
}

func NewBot(conf *config.Configuration) *Bot {
	return &Bot{
		Config: conf,
	}
}

func (b *Bot) Start() error {
	wa, err := Konek()
	if err != nil {
		return fmt.Errorf("gagal connect ke client WhatsApp : %w", err)
	}
	b.Client = wa
	b.Client.AddEventHandler(b.eventHandler)
	return nil
}

func (b *Bot) Disconnect() {
	if b.Client != nil {
		b.Client.Log.Infof("Disconnecting bot connection...")
		b.Client.Disconnect()
	}
}

// refactor dikit biar rapih, + nambah logger mesejjj ama timestamp filter sekalian
func (b *Bot) eventHandler(rawEvt any) {
	switch evt := rawEvt.(type) {
	case *events.Message:
		b.commandHandler(evt)
		go b.logMessageHandler(evt)
	}
}

func (b *Bot) commandHandler(evt *events.Message) {
	if time.Since(evt.Info.Timestamp) > 10*time.Second {
		return
	}
	messageText, ok := lib.GetText(evt)
	if !ok {
		return
	}
	if b.Config.Mode == "self" && !evt.Info.IsFromMe {
		return
	}
	var foundPrefix string
	for _, prefix := range b.Config.Prefix {
		if strings.HasPrefix(messageText, prefix) {
			foundPrefix = prefix
			break
		}
	}

	if foundPrefix == "" {
		return
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

	ctx := &commands.CommandContext{
		Ctx:     context.Background(),
		Client:  b.Client,
		Message: evt,
		Args:    args,
		RawArgs: rawArgs,
	}

	go cmd.Exec(ctx)

	var sourceName string
	if evt.Info.IsGroup {
		groupInfo, _ := b.Client.GetGroupInfo(evt.Info.Chat)
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
		groupInfo, err := b.Client.GetGroupInfo(evt.Info.Chat)
		if err != nil {
			return
		}
		if groupInfo != nil {
			source = fmt.Sprintf("Group '%s'", groupInfo.Name)
		} else {
			source = fmt.Sprintf("Group '%s'", evt.Info.Chat.String())
		}
	} else if evt.Info.Chat.Server == "newsletter" {
		newsletterIngfo, err := b.Client.GetNewsletterInfo(evt.Info.Chat)
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
	if text, ok := lib.GetText(evt); ok && text != "" {
		content = fmt.Sprintf(`"%s"`, text)
	} else if react := evt.Message.GetReactionMessage(); react != nil {
		content = react.GetText()
	}

	msgTime := evt.Info.Timestamp.Format("15:04:05")

	b.Client.Log.Infof("MSG « From: %s (%s) « In: %s", evt.Info.PushName, evt.Info.Sender.String(), source)
	b.Client.Log.Infof("    └── %s | Type: %s | Text: %s", msgTime, msgType, content)
}
