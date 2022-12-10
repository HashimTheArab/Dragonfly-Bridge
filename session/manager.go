package session

import "github.com/df-mc/dragonfly/server/player"

var sessions = map[string]*Session{}

func New(p *player.Player) *Session {
	session := &Session{Player: p}
	session.DefaultFlags()
	sessions[p.Name()] = session
	session.Scoreboard()
	return session
}

func Get(p *player.Player) *Session {
	return sessions[p.Name()]
}

func (s *Session) Close() {
	s.RemoveFromQueue()
	s.RemoveFromMatch()
	delete(sessions, s.Player.Name())
}
