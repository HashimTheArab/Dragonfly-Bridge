package commands

import "github.com/df-mc/dragonfly/server/cmd"

type Duel struct {
	Player []cmd.Target
}

func (t Duel) Run(source cmd.Source, output *cmd.Output){
	// todo
}