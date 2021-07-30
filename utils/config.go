package utils

import (
	"fmt"
	"github.com/df-mc/dragonfly/server"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
)

type Config struct {
	Staff struct {
		Admins []string
		Staff []string
	}
}

// ReadDragonflyConfig reads the configuration from the dragonfly.toml file, or creates the file if it does not yet exist.
func ReadDragonflyConfig() (server.Config, error) {
	c := server.DefaultConfig()
	if _, err := os.Stat("config/dragonfly.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return c, fmt.Errorf("failed encoding default config: %v", err)
		}
		if err := ioutil.WriteFile("config/dragonfly.toml", data, 0644); err != nil {
			return c, fmt.Errorf("failed creating config: %v", err)
		}
		return c, nil
	}
	data, err := ioutil.ReadFile("config/dragonfly.toml")
	if err != nil {
		return c, fmt.Errorf("error reading config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("error decoding config: %v", err)
	}
	return c, nil
}

// VelvetConfig reads the configuration from the velvet.toml file and sets the proper values.
func VelvetConfig() {
	c := Config{}
	data, err := ioutil.ReadFile("config/velvet.toml")
	if err != nil {
		fmt.Printf("error reading config: %v", err)
		return
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		fmt.Printf("error decoding config: %v", err)
		return
	}
	Admins = c.Staff.Admins
	Staff = c.Staff.Staff
}