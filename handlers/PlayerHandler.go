package handlers

import (
	"VBridge/session"
	"fmt"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
)

type PlayerHandler struct {
	player.NopHandler
	Player *player.Player
	Session *session.Session
}

func (PlayerHandler) HandleBlockPlace(ctx *event.Context, pos cube.Pos, block world.Block) {

}

func (PlayerHandler) HandleBlockBreak(ctx *event.Context, pos cube.Pos) {

}

func (h PlayerHandler) HandleHurt(ctx *event.Context, _ *float64, source damage.Source) {
	if (source == damage.SourceVoid{}) {

	}
	if h.Player.World().Name() == "World" {
		ctx.Cancel()
		return
	}
}

func (PlayerHandler) HandleItemDrop(ctx *event.Context, _ *entity.Item) {
	ctx.Cancel()
}

func (h PlayerHandler) HandleChat(ctx *event.Context, message *string) {
	ctx.Cancel()
	_, _ = fmt.Fprintf(chat.Global, "§a%v§f: %v\n", h.Player.Name(), *message)
}

func (h PlayerHandler) HandleQuit() {
	h.Session.Close()
}