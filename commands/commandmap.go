package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
)

const (
	NoPermission = "§cYou do not have permission to use this command!"
)

func RegisterAll(){
	cmd.Register(cmd.New("gamemode", "§bChange your gamemode!", []string{"gm"}, Gamemode{}))
	cmd.Register(cmd.New("build", "§bManage builder mode", []string{}, Build{}))
	cmd.Register(cmd.New("teleport", "§bTeleport a player to another location", []string{"tp"}, Teleport{}))
	cmd.Register(cmd.New("duel", "§bRequest a duel with another player", []string{}, Duel{}))
	cmd.Register(cmd.New("queue", "§bJoin the queue to duel!", []string{}, Queue{}))
	cmd.Register(cmd.New("sudo", "§bExecute a command or send a message as another player", []string{}, Sudo{}))
}