package dfutils

import (
	"VBridge/dfutils/exts"
	"VBridge/handlers"
	"VBridge/session"
	"VBridge/utils"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
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

	utils.VelvetConfig()

	w := Srv.World()
	w.SetDefaultGameMode(world.GameModeSurvival{})

	exts.Srv = Srv
	for {
		p, err := Srv.Accept()
		if err != nil {
			return
		}
		p.Handle(&handlers.PlayerHandler{Player: p, Session: session.New(p)}) // Player Listener and Session
	}
}