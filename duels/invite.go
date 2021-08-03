package duels

import (
	"github.com/df-mc/dragonfly/server/player"
)

type Invite struct {
	From *player.Player
	To *player.Player
}