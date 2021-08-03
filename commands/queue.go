package commands

import (
	"VBridge/session"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type Queue struct {}

func (t Queue) Run(source cmd.Source, output *cmd.Output){
	s := session.Get(source.(*player.Player))
	if s.Match != nil {
		output.Error("You cannot use this command at the moment.")
		return
	}
	s.AddToQueue()
}