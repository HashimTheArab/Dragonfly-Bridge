package session

import (
	"VBridge/dfutils/exts"
	"VBridge/duels"
	"VBridge/utils"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
	"math"
	"strconv"
	"time"
)

type Session struct {
	Player      *player.Player
	Match       *duels.Duel
	MatchPlayer *duels.DuelPlayer
	MatchInvite *duels.Invite
	Flags       uint32
}

func (s *Session) SetFlag(flag uint32) {
	s.Flags ^= 1 << flag
}

func (s *Session) HasFlag(flag uint32) bool {
	return s.Flags&(1<<flag) > 0
}

func (s *Session) IsStaff(CheckAdmin bool) bool {
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

func (s *Session) AddToQueue() {
	if s.InQueue() || s.Match != nil {
		return
	}
	duels.Queue = append(duels.Queue, s.Player)
	s.Player.Message("§aYou are now queued for Unranked §dBridge.")
	if len(duels.Queue) > 1 {
		p := duels.Queue[0]
		t := duels.Queue[1]
		var sessionList []*Session
		for _, plr := range duels.Queue {
			ses := Get(plr)
			ses.RemoveFromQueue()
			sessionList = append(sessionList, ses)
		}
		s.StartMatch(p, t, sessionList)
	}
}

func (s *Session) StartMatch(p, t *player.Player, sessionList []*Session) {
	go func() {
		id := uuid.NewString()
		if err := utils.Unzip("worlds/bridge-1.zip", "worlds/"+id); err != nil {
			panic(err)
		}
		if err := utils.WorldManager.LoadWorld(id, id); err != nil {
			panic(err)
		}
		w, _ := utils.WorldManager.World(id)
		duel := &duels.Duel{
			Arena: w,
			UUID:  id,
		}
		tPlayer := &duels.DuelPlayer{
			SpawnPoint: mgl64.Vec3{30, 69, 0},
			Player:     t,
			Team:       &duels.Team{Name: duels.TeamRed},
		}
		pPlayer := &duels.DuelPlayer{
			SpawnPoint: mgl64.Vec3{-30, 69, 0},
			Player:     p,
			Team:       &duels.Team{Name: duels.TeamBlue},
		}
		duel.Players = append(duel.Players, tPlayer)
		duel.Players = append(duel.Players, pPlayer)
		sessionList[0].MatchPlayer = tPlayer
		sessionList[1].MatchPlayer = pPlayer
		for num, pl := range duel.Players {
			w.AddEntity(pl.Player)
			pl.Player.Teleport(pl.SpawnPoint)
			pl.Player.SetGameMode(world.GameModeSurvival)
			duel.SendKit(pl.Player, pl.Team.Color())
			pl.Player.SetImmobile()
			sessionList[num].Match = duel
		}
		duels.Matches.Add(duel)
		p.Message("\n§l§aMatch Found!\n\n§rOpponent: §b" + t.Name() + "\nPing: §b" + strconv.FormatFloat(math.Round(p.Latency().Seconds()), 'f', 0, 64) + "ms\n")
		t.Message("\n§l§aMatch Found!\n\n§rOpponent: §b" + p.Name() + "\nPing: §b" + strconv.FormatFloat(t.Latency().Seconds(), 'f', 0, 64) + "ms\n")

		duel.BroadcastScoreboard()

		dv := duels.Matches.Get(duel)
		dv.Waiting.Store(true)
		for i := 5; i > 0; i-- {
			if dv == nil {
				return
			}
			for _, pl := range duel.Players {
				pl.Player.Message("Starting in §b" + strconv.Itoa(i) + "§f...")
			}
			if i <= 3 {
				for _, pl := range duel.Players {
					pl.Player.PlaySound(sound.Note{
						Instrument: sound.Piano(),
						Pitch:      0,
					})
				}
			}
			time.Sleep(1 * time.Second)
		}
		for _, pl := range duel.Players {
			pl.Player.Message("Match Started!")
			duel.SendKit(pl.Player, pl.Team.Color())
			pl.Player.SetMobile()
			pl.Player.PlaySound(sound.Note{
				Instrument: sound.Piano(),
				Pitch:      1,
			})
		}
		dv.Waiting.Store(false)
		go func() {
			for {
				if duel == nil {
					return
				}
				duel.ElapsedTime.Add(1)
				duel.BroadcastScoreboard()
				time.Sleep(1 * time.Second)
			}
		}()
	}()
}

func (s *Session) InQueue() bool {
	for _, p := range duels.Queue {
		if p.Name() == s.Player.Name() {
			return true
		}
	}
	return false
}

func (s *Session) RemoveFromQueue() {
	if s.InQueue() {
		var New []*player.Player
		for _, p := range duels.Queue {
			if p.Name() != s.Player.Name() {
				New = append(New, p)
			}
		}
		duels.Queue = New
	}
}

func (s *Session) RemoveFromMatch() {
	if s.Match != nil {
		for _, p := range s.Match.Players {
			ses := Get(p.Player)
			if ses != nil {
				ses.MatchPlayer = nil
				ses.Match = nil
			}
		}
		s.Match.RemovePlayer(s.Player)
	}
}

func (s *Session) TeleportToSpawn() {
	utils.Srv.World().AddEntity(s.Player)
	s.Player.Teleport(utils.DefaultSpawn)
}

func (s *Session) Scoreboard() {
	go func() {
		for {
			if s == nil || s.Player == nil {
				return
			}
			if s.Match == nil {
				s.DefaultScoreboard()
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func (s *Session) DefaultScoreboard() {
	sb := scoreboard.New("§3Velvet Bridge")
	_, _ = sb.WriteString("§bOnline: §3" + exts.Online.String() + "\n§bMatches: §3" + exts.Matches.String() + "\n\n§bvelvetpractice.live")
	s.Player.SendScoreboard(sb)
}
