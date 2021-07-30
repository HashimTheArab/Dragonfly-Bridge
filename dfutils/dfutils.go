package dfutils

import (
	"VBridge/dfutils/exts"
	"VBridge/handlers"
	"VBridge/session"
	"VBridge/utils"
	"github.com/df-mc/dragonfly/server"
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

	config, err := utils.ReadDragonflyConfig()
	if err != nil {
		log.Fatalln(err)
	}

	Srv := server.New(&config, log)
	Srv.CloseOnProgramEnd()
	if err := Srv.Start(); err != nil {
		log.Fatalln(err)
	}

	w := Srv.World()
	w.SetDefaultGameMode(world.GameModeSurvival{})
	w.SetTime(0)
	w.StopTime()

	utils.VelvetConfig()
	exts.Srv = Srv
	exts.DefaultSpawn = Srv.World().Spawn().Vec3().Add(mgl64.Vec3{0, 1, 0})

	for {
		p, err := Srv.Accept()
		if err != nil {
			return
		}
		p.Handle(&handlers.PlayerHandler{Player: p, Session: session.New(p)}) // Player Listener and Session
	}
}