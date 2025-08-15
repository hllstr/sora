package commands

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
	"runtime"
	"sora/lib"
	"strings"
	"time"
)

func init() {
	Plugin(Cmd{
		Name:  "ping",
		Alias: []string{"p"},
		Desc:  "Check bot response.",
		// inline code
		Exec: func(ctx *CommandContext) {
			now := time.Now()
			latency := now.Sub(ctx.Message.Info.Timestamp)
			ctx.Reply(fmt.Sprintf("*Pong!* `%s`", latency.Round(time.Millisecond)))
		},
	})

	Plugin(Cmd{
		Name:  "info",
		Alias: []string{"i"},
		Desc:  "Show Information.",
		Exec:  info,
	})
}

func getMemoryUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024)
}

func countCommands() int {
	count := 0
	printed := make(map[string]bool)
	for _, cmd := range Commands {
		if _, omkeh := printed[cmd.Name]; !omkeh {
			count++
			printed[cmd.Name] = true
		}
	}
	return count
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func info(ctx *CommandContext) {
	myrepo := "https://githuh.com/hllstr/sora"
	bullet := "✦"
	osInfo, _ := host.Info()
	cpuInfo, _ := cpu.Info()
	vmInfo, _ := mem.VirtualMemory()

	platform := fmt.Sprintf("%s/%s", capitalize(runtime.GOOS), runtime.GOARCH)
	memUsage := getMemoryUsage()

	registeredCmds := countCommands()
	//	cpuUseggs, _ := cpu.Percent(time.Second, false)
	totalCores, err := cpu.Counts(true)
	if err != nil {
		return
	}
	var sb strings.Builder
	sb.WriteString("*Sora (空) adalah sebuah Base Bot Wangsaf yang dibangun menggunakan Go (Golang) dan library Whatsmeow.*\n")
	sb.WriteString(fmt.Sprintf("%s *OS:* %s\n", bullet, capitalize(osInfo.Platform)))
	sb.WriteString(fmt.Sprintf("%s *CPU:* %s\n", bullet, cpuInfo[0].ModelName))
	sb.WriteString(fmt.Sprintf("%s *Cores:* %d\n", bullet, totalCores))
	sb.WriteString(fmt.Sprintf("%s *Freq:* %.2f Mhz\n", bullet, cpuInfo[0].Mhz))
	sb.WriteString(fmt.Sprintf("%s *RAM:* %.2f GB / %.2f GB\n", bullet, float64(vmInfo.Used)/1024/1024/1024, float64(vmInfo.Total)/1024/1024/1024))
	sb.WriteString(fmt.Sprintf("%s *Platform:* %s\n", bullet, platform))

	sb.WriteString(fmt.Sprintf("%s *Memory Usage:* %s (Sora Alloc)\n", bullet, memUsage))
	// CPU usex bikin delay 1 sex jir buat ngukurnya, gw comment aja, kalo lu mau nampilin cpu usage bisa di uncomment sajah.
	// sb.WriteString(fmt.Sprintf("%s *CPU Usage:* %.2f%%\n", bullet, cpuUseggs[0]))
	sb.WriteString(fmt.Sprintf("%s *Commands:* %d Registered\n", bullet, registeredCmds))
	sb.WriteString(fmt.Sprintf("> *Download/Clone script nya hanya di %s*\n> *Jangan lupa kasih star hehe :v*", myrepo))
	thumb := "https://i.pinimg.com/originals/7e/2b/fb/7e2bfb2629b8e72826b818a5e749839b.jpg"
	finalText := sb.String()
	ctxInfo := &waE2E.ContextInfo{
		Participant: proto.String(ctx.Message.Info.Sender.String()),
		ExternalAdReply: &waE2E.ContextInfo_ExternalAdReplyInfo{
			Title:                 proto.String("✨ System Information ✨"),
			Body:                  proto.String("A lightweight & high-performance base bot."),
			ThumbnailURL:          proto.String(thumb),
			RenderLargerThumbnail: proto.Bool(true),
			MediaURL:              proto.String(thumb),
			SourceURL:             proto.String(myrepo),
			WtwaAdFormat:          proto.Bool(true),
			MediaType:             waE2E.ContextInfo_ExternalAdReplyInfo_IMAGE.Enum(),
		},
	}
	if exp, omkeh := lib.GetEphemeralDuration(ctx.Message); omkeh {
		ctxInfo.Expiration = proto.Uint32(exp)
	}
	extra := whatsmeow.SendRequestExtra{}
	if ctx.Message.Info.Chat.Server == types.GroupServer {
		ownJID := ctx.Client.Store.ID
		extra.TargetJID = []types.JID{*ownJID}
	}

	_, err = ctx.Client.SendMessage(context.Background(), ctx.Message.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:        proto.String(finalText),
			ContextInfo: ctxInfo,
			MatchedText: proto.String(myrepo),
		},
	}, extra)
	if err != nil {
		return
	}
}

// TODO : nyieun menu.
// Bypass juga kek nya mau gw buat function aja di lib
// biar gak import package whatsmeow ama types mulu
