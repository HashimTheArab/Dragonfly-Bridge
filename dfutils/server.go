package dfutils

import (
	"VBridge/handlers"
	"VBridge/session"
	"fmt"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

var Srv *server.Server

func StartServer() {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel

	chat.Global.Subscribe(chat.StdoutSubscriber{})

	config, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}

	Srv = server.New(&config, log)
	Srv.CloseOnProgramEnd()
	if err := Srv.Start(); err != nil {
		log.Fatalln(err)
	}

	w := Srv.World()
	w.SetDefaultGameMode(world.GameModeSurvival{})

	w.Handle(handlers.NewWorldHandler(w)) // World Listener

	for {
		p, err := Srv.Accept()
		if err != nil {
			return
		}
		p.Handle(&handlers.PlayerHandler{Player: p, Session: session.New(p)}) // Player Listener and Session
	}
}


// readConfig reads the configuration from the config.toml file, or creates the file if it does not yet exist.
func readConfig() (server.Config, error) {
	c := server.DefaultConfig()
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return c, fmt.Errorf("failed encoding default config: %v", err)
		}
		if err := ioutil.WriteFile("config.toml", data, 0644); err != nil {
			return c, fmt.Errorf("failed creating config: %v", err)
		}
		return c, nil
	}
	data, err := ioutil.ReadFile("config.toml")
	if err != nil {
		return c, fmt.Errorf("error reading config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("error decoding config: %v", err)
	}
	return c, nil
}