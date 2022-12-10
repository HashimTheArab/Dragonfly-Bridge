package dfutils

import (
	"VBridge/cblock"
	"VBridge/dfutils/exts"
	"VBridge/handlers"
	"VBridge/session"
	"VBridge/utils"
	"VBridge/worldmanager"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	config, err := utils.ReadDragonflyConfig(log)
	if err != nil {
		log.Fatalln(err)
	}

	Srv := config.New()
	Srv.CloseOnProgramEnd()

	Srv.Listen()

	world.RegisterBlock(cblock.EndPortal{})

	utils.WorldManager = worldmanager.New(Srv, "worlds/", log)
	w := Srv.World()
	w.SetDefaultGameMode(world.GameModeSurvival)
	w.SetTime(0)
	w.StopTime()

	utils.VelvetConfig()
	utils.Srv = Srv
	utils.DefaultSpawn = Srv.World().Spawn().Vec3().Add(mgl64.Vec3{0, 1, 0})
	exts.TrackData()

	for Srv.Accept(func(p *player.Player) {
		p.Handle(&handlers.PlayerHandler{Player: p, Session: session.New(p)}) // Player Listener and Session
		p.ShowCoordinates()
	}) {
	}
}
