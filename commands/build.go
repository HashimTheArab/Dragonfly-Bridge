package commands

import (
	"VBridge/session"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
)

type BuildCommand struct {
	Player []cmd.Target `optional:""`
}

func (t BuildCommand) Run(source cmd.Source, output *cmd.Output) {
	p := source.(*player.Player)
	if !session.Get(p).HasFlag(session.Admin) {
		p.Message(NoPermission)
		return
	}

	if len(t.Player) > 0 {
		if target, ok := t.Player[0].(*player.Player); ok {
			s := session.Get(target)
			if s.HasFlag(session.Builder) {
				target.Message("§b" + p.Name() + "§7 has set you out of §bbuilder §7mode!")
				p.Message("§b" + target.Name() + "§7 is out of §bbuilder §7mode!")
			} else {
				target.Message("§b" + p.Name() + "§7 has set you §ain §bbuilder §7mode!")
				p.Message("§b" + target.Name() + "§7 is now in §bbuilder §7mode!")
			}
			s.SetFlag(session.Builder)
			return
		}
		output.Errorf("§6%s §cis not online!", t.Player)
	}

	s := session.Get(p)
	if s.HasFlag(session.Builder) {
		p.Message("§7You are now out of §bbuilder §7mode!")
	} else {
		p.Message("§7You are now in §bbuilder §7mode!")
	}
	s.SetFlag(session.Builder)
}