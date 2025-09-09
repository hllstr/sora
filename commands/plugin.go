package commands

import (
	"context"
	"log"
	"sora/config"
	"sora/lib"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

type CommandContext struct {
	Ctx     context.Context
	Client  *whatsmeow.Client
	Message *events.Message
	Args    []string
	RawArgs string
	Conf    *config.Configuration
}

// Disini helper function Reply dari package lib dibungkus lagi, biar pas dipake di package commands
// bakalan jauh lebih simpel...
// Selain Reply() nanti lu juga bisa tambahin function lain kaya Send() atau React() atau Edit()
func (c *CommandContext) Reply(text string) {
	if _, err := lib.Reply(c.Client, c.Message, text); err != nil {
		log.Printf("error when replying message : %v", err)
	}
}

type PermissionLevel int

const (
	Public PermissionLevel = iota
	Owner
)

// ini struct buat Cmd nanti disini bisa nambahin kaya permission dll.
// misal kalo mau nambahin OwnerOnly nanti tambahin field nya disini
// terus handle logic nya di eventHandler
type Cmd struct {
	Name       string
	Alias      []string
	Desc       string
	Category   string
	Permission PermissionLevel
	Exec       func(ctx *CommandContext)
}

var Commands = make(map[string]Cmd)

func Plugin(cmd Cmd) {
	Commands[cmd.Name] = cmd
	for _, alias := range cmd.Alias {
		Commands[alias] = cmd
	}
}
