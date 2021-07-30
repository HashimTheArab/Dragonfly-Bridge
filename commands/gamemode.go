package commands

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"strings"
)

type GamemodeCommand struct {
	GameMode string
	Player []cmd.Target `optional:""`
}

func (t GamemodeCommand) Run(source cmd.Source, output *cmd.Output) {
	var gm world.GameMode
	switch strings.ToLower(t.GameMode) {
		case "survival", "0":
			gm = world.GameModeSurvival{}
		case "creative", "1":
			gm = world.GameModeCreative{}
		case "adventure", "2":
			gm = world.GameModeAdventure{}
		case "spectator", "3":
			gm = world.GameModeSpectator{}
		default:
			output.Error("§cInvalid Gamemode!")
			return
	}

	if len(t.Player) > 0 {
		if target, ok := t.Player[0].(*player.Player); ok {
			target.SetGameMode(gm)
			target.Message("§bYour gamemode has been set to §f" + t.GameMode + "!")
			target.Message(fmt.Sprintf("§bYou have set §f%s's §bgamemode to §f%s!", target.Name(), t.GameMode))
			return
		}
		output.Errorf("§6%s §cis not online!", t.Player)
	}

	p := source.(*player.Player)
	p.SetGameMode(gm)
	output.Print("§bYour gamemode has been set to §f" + t.GameMode + "!")
}