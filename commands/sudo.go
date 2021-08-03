package commands

import (
	"VBridge/session"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"strings"
)

type Sudo struct {
	Targets []cmd.Target
	Message string
}

func (t Sudo) Run(source cmd.Source, _ *cmd.Output){
	p := source.(*player.Player)
	if !session.Get(p).HasFlag(session.Admin) {
		p.Message(NoPermission)
		return
	}
	if len(t.Targets) > 0 {
		var names []string
		var command bool
		if strings.HasPrefix(t.Message, "/") {
			command = true
		}
		for _, target := range t.Targets {
			if tar, ok := target.(*player.Player); ok {
				if command {
					tar.ExecuteCommand(t.Message)
				} else {
					tar.Chat(t.Message)
				}
				names = append(names, tar.Name())
			}
		}
		p.Message("§aMessage sent as §e" + strings.Join(names, ", ") + "!")
	}
}