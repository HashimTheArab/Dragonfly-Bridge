package commands

import (
	"VBridge/session"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type Teleport struct {
	Player []cmd.Target `optional:""`
}

func (t Teleport) Run(source cmd.Source, output *cmd.Output){
	p := source.(*player.Player)
	if !session.Get(p).HasFlag(session.Staff) {
		p.Message(NoPermission)
		return
	}
	// todo
}