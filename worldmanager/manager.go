package worldmanager

import (
	"fmt"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/mcdb"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

// WorldManager manages multiple worlds, dragonfly does not have multi-world management itself,
// so we must implement it ourselves.
type WorldManager struct {
	s *server.Server

	folderPath string

	log *logrus.Logger

	worldsMu sync.RWMutex
	worlds   map[string]*world.World
}

// New ...
func New(server *server.Server, folderPath string, log *logrus.Logger) *WorldManager {
	_ = os.Mkdir(folderPath, 0644)
	return &WorldManager{
		s:          server,
		folderPath: folderPath,
		log:        log,
		worlds:     make(map[string]*world.World),
	}
}

// DefaultWorld ...
func (m *WorldManager) DefaultWorld() *world.World {
	return m.s.World()
}

// Worlds ...
func (m *WorldManager) Worlds() []*world.World {
	m.worldsMu.RLock()
	defer m.worldsMu.RUnlock()

	worlds := make([]*world.World, 0, len(m.worlds))
	for _, w := range m.worlds {
		worlds = append(worlds, w)
	}
	return worlds
}

// World ...
func (m *WorldManager) World(name string) (*world.World, bool) {
	m.worldsMu.RLock()
	w, ok := m.worlds[name]
	m.worldsMu.RUnlock()
	return w, ok
}

// LoadWorld ...
func (m *WorldManager) LoadWorld(folderName, worldName string, simulationDistance int) error {
	w := world.New(m.log, simulationDistance)
	p, err := mcdb.New(m.folderPath + "/" + folderName)
	if err != nil {
		return fmt.Errorf("error loading world: %v", err)
	}

	settings := p.Settings()
	settings.Name = worldName
	p.SaveSettings(settings)

	w.Provider(p)
	if _, ok := m.World(w.Name()); ok {
		return fmt.Errorf("world is already loaded")
	}

	m.worldsMu.Lock()
	m.worlds[w.Name()] = w
	m.worldsMu.Unlock()
	return nil
}

// UnloadWorld ...
func (m *WorldManager) UnloadWorld(w *world.World) error {
	if w == m.DefaultWorld() {
		return fmt.Errorf("the default world cannot be unloaded")
	}

	if _, ok := m.World(w.Name()); !ok {
		return fmt.Errorf("world isn't loaded")
	}

	m.log.Debugf("Unloading world '%v'\n", w.Name())
	for _, p := range m.s.Players() {
		if p.World() == w {
			m.DefaultWorld().AddEntity(p)
			p.Teleport(m.DefaultWorld().Spawn().Vec3Middle())
		}
	}

	m.worldsMu.Lock()
	delete(m.worlds, w.Name())
	m.worldsMu.Unlock()

	if err := w.Close(); err != nil {
		return fmt.Errorf("error closing world: %v", err)
	}
	m.log.Debugf("Unloaded world '%v'\n", w.Name())
	return nil
}

// Close ...
func (m *WorldManager) Close() error {
	m.worldsMu.Lock()
	for _, w := range m.worlds {
		m.log.Debugf("Closing world '%v'\n", w.Name())
		if err := w.Close(); err != nil {
			return err
		}
	}
	m.worlds = map[string]*world.World{}
	m.worldsMu.Unlock()
	return nil
}