package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
)

const (
	NoPermission = "§cYou do not have permission to use this command!"
)

func RegisterAll(){
	cmd.Register(cmd.New("gamemode", "§bChange your gamemode!", []string{"gm"}, GamemodeCommand{}))
	cmd.Register(cmd.New("build", "§bManage builder mode", []string{}, BuildCommand{}))
}