package commands

import (
	"os/exec"
	"strings"
)

func init() {
	Plugin(Cmd{
		Name:       "exec",
		Alias:      []string{"$"},
		Desc:       "Execute command.",
		Category:   "owner",
		Permission: Owner,
		Exec:       execute,
	})
}

func execute(ctx *CommandContext) {
	if ctx.RawArgs == "" {
		return
	}

	cmd := exec.Command("bash", "-c", ctx.RawArgs)
	output, err := cmd.CombinedOutput()
	if err != nil {
		errjir := strings.TrimSpace(string(output)) + "\n" + error.Error(err)
		ctx.Reply(errjir)
		return
	}
	ctx.Reply(strings.TrimSpace(string(output)))
}
