package handlers

import (
	"VBridge/dfutils/exts"
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

func (h PlayerHandler) HandleBlockPlace(ctx *event.Context, _ cube.Pos, _ world.Block) {
	if h.Player.World().Name() == "World" {
		if !h.Session.HasFlag(session.Builder) {
			ctx.Cancel()
		}
	} else {
		// determine if they can place blocks in the match
	}
}

func (h PlayerHandler) HandleBlockBreak(ctx *event.Context, _ cube.Pos) {
	if h.Player.World().Name() == "World" {
		if !h.Session.HasFlag(session.Builder) {
			ctx.Cancel()
		}
	} else {
		// determine if they can break blocks in the match
	}
}

func (h PlayerHandler) HandleHurt(ctx *event.Context, _ *float64, source damage.Source) {
	if (source == damage.SourceVoid{}) {
		if h.Player.World().Name() == "World" {
			ctx.Cancel()
			h.Player.Teleport(exts.Srv.World().Spawn().Vec3())
		} else {
			// session.game.get their spawn
		}
	}
	if h.Player.World().Name() == "World" {
		ctx.Cancel()
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