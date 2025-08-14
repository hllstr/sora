package bot

import (
	"context"
	"fmt"
	"sora/commands"
	"sora/config"
	"sora/lib"
	"strings"

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

func (b *Bot) eventHandler(rawEvt any) {
	switch evt := rawEvt.(type) {
	case *events.Message:
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

		b.Client.Log.Infof("Executing command '%s' for %s", commandName, evt.Info.PushName)
		go cmd.Exec(ctx)
	}
}
