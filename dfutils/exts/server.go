package exts

import (
	"github.com/df-mc/dragonfly/server"
	"github.com/go-gl/mathgl/mgl64"
)

var Srv *server.Server // damn import cycles
var DefaultSpawn mgl64.Vec3