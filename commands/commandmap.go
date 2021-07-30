package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
)

func RegisterAll(){
	cmd.Register(cmd.New("gamemode", "§bChange your gamemode!", []string{"gm"}, GamemodeCommand{}))
}