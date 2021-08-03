package handlers

import (
	"VBridge/cblock"
	"VBridge/session"
	"fmt"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
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

func (h *PlayerHandler) HandleHurt(ctx *event.Context, _ *float64, source damage.Source) {
	if source == (damage.SourceVoid{}) {
		ctx.Cancel()
		m := h.Session.Match
		if m == nil {
			h.Session.TeleportToSpawn()
		} else {
			m.BroadcastMessage(h.Session.MatchPlayer.Team.Color() + h.Player.Name() + " §7died in the void.")
			m.TeleportPlayerToSpawn(h.Player)
			if h.Session != nil && h.Session.MatchPlayer != nil {
				m.SendKit(h.Player, h.Session.MatchPlayer.Team.Color())
			}
		}
	} else if h.Player.World().Name() == "World" || (source == damage.SourceFall{}){
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

func (h PlayerHandler) HandleDeath(source damage.Source) {
	m := h.Session.Match
	if m != nil {
		if s, ok := source.(damage.SourceEntityAttack); ok {
			for _, p := range m.Players {
				if p.Player.Name() == s.Attacker.Name() {
					m.BroadcastMessage(h.Session.MatchPlayer.Team.Color() + h.Player.Name() + "§7 was killed by " + p.Team.Color() + p.Player.Name())
					p.Kills++
				}
			}
		}
	}
}

func (h PlayerHandler) HandleMove(_ *event.Context, pos mgl64.Vec3, _ float64, _ float64) {
	if h.Player.World().Block(cube.PosFromVec3(pos)) == (block.Cake{}) || h.Player.World().Block(cube.PosFromVec3(pos)) == (cblock.EndPortal{}) {
		if h.Player.World().Name() == "World" {
			h.Session.AddToQueue()
		} else {
			m := h.Session.Match
			if m != nil && !m.Over() && !m.Waiting.Load() {
				m.Score(h.Session.MatchPlayer)
				h.Player.Teleport(h.Session.MatchPlayer.SpawnPoint)
			}
		}
	}
}

func (h PlayerHandler) HandleQuit() {
	h.Session.Close()
}