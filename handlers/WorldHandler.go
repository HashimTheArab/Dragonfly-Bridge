package handlers

import (
	"github.com/df-mc/dragonfly/server/world"
)

type WorldHandler struct {
	world.NopHandler
	World *world.World
}

func NewWorldHandler(World *world.World) *WorldHandler {
	return &WorldHandler{World: World}
}