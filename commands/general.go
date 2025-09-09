package commands

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sora/lib"
	"sort"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

func init() {
	Plugin(Cmd{
		Name:     "ping",
		Alias:    []string{"p"},
		Desc:     "Check Response.",
		Category: "general",
		Exec: func(ctx *CommandContext) {
			now := time.Now()
			latency := now.Sub(ctx.Message.Info.Timestamp)
			ctx.Reply(fmt.Sprintf("*Pong!* `%s`", latency.Round(time.Millisecond)))
		},
	})

	Plugin(Cmd{
		Name:     "info",
		Alias:    []string{"i"},
		Category: "general",
		Desc:     "Show Information.",
		Exec:     info,
	})

	Plugin(Cmd{
		Name:     "help",
		Alias:    []string{"h"},
		Category: "general",
		Desc:     "Show menu.",
		Exec:     help,
	})
}

func formatUptime(d time.Duration) string {
	d = d.Round(time.Second)
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	d -= minutes * time.Minute
	seconds := d / time.Second

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%dd", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}
	return strings.Join(parts, " ")
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
	totalCores, err := cpu.Counts(true)
	if err != nil {
		return
	}

	var uptimeStr string
	pid := os.Getpid()
	p, err := process.NewProcess(int32(pid))
	if err == nil {
		createTime, err := p.CreateTime()
		if err == nil {
			startTime := time.Unix(0, createTime*int64(time.Millisecond))
			uptime := time.Since(startTime)
			uptimeStr = formatUptime(uptime)
		}
	}

	if uptimeStr == "" {
		uptimeStr = "Calculating..."
	}
	var sb strings.Builder
	sb.WriteString("*Sora (空) is a Simple Base Bot WhatsApp written in Go using Whatsmeow Library.*\n")
	sb.WriteString(fmt.Sprintf("%s *OS:* %s\n", bullet, capitalize(osInfo.Platform)))
	sb.WriteString(fmt.Sprintf("%s *CPU:* %s\n", bullet, cpuInfo[0].ModelName))
	sb.WriteString(fmt.Sprintf("%s *Cores:* %d\n", bullet, totalCores))
	sb.WriteString(fmt.Sprintf("%s *Freq:* %.2f Mhz\n", bullet, cpuInfo[0].Mhz))
	sb.WriteString(fmt.Sprintf("%s *RAM:* %.2f GB / %.2f GB\n", bullet, float64(vmInfo.Used)/1024/1024/1024, float64(vmInfo.Total)/1024/1024/1024))
	sb.WriteString(fmt.Sprintf("%s *Platform:* %s\n", bullet, platform))
	sb.WriteString(fmt.Sprintf("%s *Memory Usage:* %s (Sora Alloc)\n", bullet, memUsage))
	sb.WriteString(fmt.Sprintf("%s *Uptime:* %s\n", bullet, uptimeStr))
	sb.WriteString(fmt.Sprintf("%s *Commands:* %d Registered\n", bullet, registeredCmds))
	sb.WriteString(fmt.Sprintf("> *Source Code : %s.*\n> *Don't forget to give a star ✨٩(ˊᗜˋ**)و✨", myrepo))
	thumb := "https://i.pinimg.com/originals/7e/2b/fb/7e2bfb2629b8e72826b818a5e749839b.jpg"
	finalText := sb.String()
	ctxInfo := &waE2E.ContextInfo{
		Participant: proto.String(ctx.Message.Info.Sender.String()),
		ExternalAdReply: &waE2E.ContextInfo_ExternalAdReplyInfo{
			Title:                 proto.String("✨ System Information ✨"),
			Body:                  proto.String("A lightweight & high-performance base bot."),
			ThumbnailURL:          proto.String(thumb),
			RenderLargerThumbnail: proto.Bool(true),
			SourceURL:             proto.String(myrepo),
			WtwaAdFormat:          proto.Bool(true),
			MediaType:             waE2E.ContextInfo_ExternalAdReplyInfo_IMAGE.Enum(),
		},
	}
	if exp, omkeh := lib.GetEphemeralDuration(ctx.Message); omkeh {
		ctxInfo.Expiration = proto.Uint32(exp)
	}
	// bypass participant sekarang lebih gampang dipake, udah gw jadiin function di lib
	bypass := lib.Bypass(ctx.Client, ctx.Message.Info.Chat)
	_, err = ctx.Client.SendMessage(context.Background(), ctx.Message.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:        proto.String(finalText),
			ContextInfo: ctxInfo,
			MatchedText: proto.String(myrepo),
		},
	}, bypass)
	if err != nil {
		return
	}
}

// finally akhirnya menu gweh jadi
func help(ctx *CommandContext) {
	categorized := make(map[string][]Cmd)

	// biar gak duplikat
	uniqueCmds := make(map[string]bool)
	for _, cmd := range Commands {
		if _, exists := uniqueCmds[cmd.Name]; !exists {
			categorized[cmd.Category] = append(categorized[cmd.Category], cmd)
			uniqueCmds[cmd.Name] = true
		}
	}

	var sb strings.Builder
	sb.WriteString("*Sora (空) - Command Menu*\n")
	sb.WriteString("*━━━━━━━━━━━━━━━━━━*\n\n")

	var categories []string
	for cat := range categorized {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	maxLen := 0
	for _, cmds := range categorized {
		for _, cmd := range cmds {
			var aliasStr string
			if len(cmd.Alias) > 0 {
				aliasStr = fmt.Sprintf(" (%s)", strings.Join(cmd.Alias, ", "))
			}
			currentLen := len(fmt.Sprintf("✦ /%s%s", cmd.Name, aliasStr))
			if currentLen > maxLen {
				maxLen = currentLen
			}
		}
	}

	for _, category := range categories {
		sb.WriteString(fmt.Sprintf("> *%s*\n", strings.ToUpper(category)))
		sort.Slice(categorized[category], func(i, j int) bool {
			return categorized[category][i].Name < categorized[category][j].Name
		})

		for _, cmd := range categorized[category] {

			var aliasStr string
			if len(cmd.Alias) > 0 {
				aliasStr = fmt.Sprintf(" (%s)", strings.Join(cmd.Alias, ", "))
			}

			cmdLine := fmt.Sprintf("✦ /%s%s", cmd.Name, aliasStr)

			padding := strings.Repeat(" ", maxLen-len(cmdLine)+2)

			sb.WriteString(fmt.Sprintf("`%s%s– %s`\n", cmdLine, padding, cmd.Desc))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("*━━━━━━━━━━━━━━━━━━*\n")
	sb.WriteString("_Check out Sora, a simple, lightweight, and high-performance WhatsApp base bot on GitHub:_\n")
	sb.WriteString("_https://github.com/hllstr/sora_")
	// gw copas dari githuh mwhehe
	thumb := "https://camo.githubusercontent.com/059157854c0fdb6d3f3976443bdf4b20439b0533c086535c0cdebb3055c25de0/68747470733a2f2f692e70696e696d672e636f6d2f6f726967696e616c732f33622f61312f62312f33626131623162656334383433376630373430373463653730653935613036662e6a7067"
	newText := sb.String()
	ctxInfo := &waE2E.ContextInfo{
		Participant: proto.String(ctx.Message.Info.Sender.String()),
		ExternalAdReply: &waE2E.ContextInfo_ExternalAdReplyInfo{
			Title:                 proto.String("✨ Commands Menu ✨"),
			Body:                  proto.String("A lightweight & high-performance base bot."),
			ThumbnailURL:          proto.String(thumb),
			RenderLargerThumbnail: proto.Bool(true),
			SourceURL:             proto.String("https://githuh.com/hllstr/sora"),
			WtwaAdFormat:          proto.Bool(true),
			MediaType:             waE2E.ContextInfo_ExternalAdReplyInfo_IMAGE.Enum(),
		},
	}
	if exp, omkeh := lib.GetEphemeralDuration(ctx.Message); omkeh {
		ctxInfo.Expiration = proto.Uint32(exp)
	}
	bypass := lib.Bypass(ctx.Client, ctx.Message.Info.Chat)
	_, err := ctx.Client.SendMessage(context.Background(), ctx.Message.Info.Chat, &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:        proto.String(newText),
			ContextInfo: ctxInfo,
		},
	}, bypass)
	if err != nil {
		return
	}
}
