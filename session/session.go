package session

import (
	"VBridge/utils"
	"github.com/df-mc/dragonfly/server/player"
)

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

func (s Session) IsStaff(CheckAdmin bool) bool {
	name := s.Player.Name()
	xuid := s.Player.XUID()
	if CheckAdmin {
		for _, v := range utils.Admins {
			if v == name || v == xuid {
				return true
			}
		}
		return false
	}
	for _, v := range utils.Staff {
		if v == name || v == xuid {
			return true
		}
	}
	return false
}

func (s *Session) DefaultFlags() {
	if s.IsStaff(true) {
		s.SetFlag(Admin)
		s.SetFlag(Staff)
	} else if s.IsStaff(false) {
		s.SetFlag(Staff)
	}
}