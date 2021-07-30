package session

import "github.com/df-mc/dragonfly/server/player"

type Session struct {
	Player *player.Player
	Flags uint32
}

func (s Session) Join() {
}

func (s *Session) SetFlag(flag uint32) {
	s.Flags ^= 1 << flag
}

func (s Session) HasFlag(flag uint32) bool {
	return s.Flags & (1 << flag) == 1
}

func (s Session) DefaultFlags() *Session {
	return &s
}