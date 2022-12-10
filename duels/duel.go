package duels

import (
	"VBridge/utils"
	"fmt"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/player/title"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"go.uber.org/atomic"
	"os"
	"strconv"
	"sync"
	"time"
)

type Duel struct {
	Players     []*DuelPlayer
	Arena       *world.World
	UUID        string
	ElapsedTime atomic.Uint32
	Waiting     atomic.Bool
}

type DuelPlayer struct {
	Kills      uint32
	SpawnPoint mgl64.Vec3
	Player     *player.Player
	Team       *Team
}

type Team struct {
	Name   string
	Points uint32
}

type matches struct {
	List  map[string]*Duel
	mutex sync.RWMutex
}

const (
	TeamRed  = "Red"
	TeamBlue = "Blue"
)

var Matches = matches{List: map[string]*Duel{}}
var Queue []*player.Player

func (m *matches) Add(d *Duel) {
	m.mutex.Lock()
	m.List[d.UUID] = d
	m.mutex.Unlock()
}

func (m *matches) Get(d *Duel) *Duel {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.List[d.UUID]
}

func (m *matches) Delete(d *Duel) {
	m.mutex.Lock()
	delete(m.List, d.UUID)
	m.mutex.Unlock()
}

func (m *matches) Amount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.List)
}

func (*Duel) SendKit(p *player.Player, ColorString string) {
	var ItemColor item.Colour
	if ColorString == "§c" {
		ItemColor = item.ColourRed()
	} else {
		ItemColor = item.ColourBlue()
	}
	inv := p.Inventory()
	inv.Clear()
	_ = inv.SetItem(0, item.NewStack(item.Sword{Tier: item.ToolTierIron}, 1))
	// _ = inv.SetItem(1, item.NewStack(item.Bow, 1)) projectiles wen
	_ = inv.SetItem(1, item.NewStack(item.Pickaxe{Tier: item.ToolTierDiamond}, 1))
	_ = inv.SetItem(2, item.NewStack(block.Concrete{Colour: ItemColor}, 64))
	_ = inv.SetItem(8, item.NewStack(item.GoldenApple{}, 4))
	// _ = inv.SetItem(8, item.NewStack(item.Arrow, 1)) projectiles wen

	a := p.Armour()
	a.SetHelmet(item.NewStack(item.Helmet{Tier: item.ArmourTierLeather{}}, 1))
	a.SetChestplate(item.NewStack(item.Chestplate{Tier: item.ArmourTierLeather{}}, 1))
	a.SetLeggings(item.NewStack(item.Leggings{Tier: item.ArmourTierLeather{}}, 1))
	a.SetBoots(item.NewStack(item.Boots{Tier: item.ArmourTierLeather{}}, 1))
}

func (d *Duel) TeleportPlayerToSpawn(p *player.Player) {
	d.Arena.AddEntity(p)
	for _, pl := range d.Players {
		if pl.Player.Name() == p.Name() {
			p.Teleport(pl.SpawnPoint)
		}
	}
}

func (d *Duel) RemovePlayer(p *player.Player) {
	go func() {
		var winner *DuelPlayer
		var loser *DuelPlayer
		for _, pl := range d.Players {
			if p.Name() == pl.Player.Name() {
				loser = pl
			} else {
				winner = pl
			}
			pl.Player.SetMobile()
		}
		if winner == nil {
			winner = d.GetDuelPlayer(p)
		}
		t := title.New(winner.Team.Color()+winner.Team.Name+" won").WithSubtitle(winner.Team.Color()+strconv.FormatUint(uint64(winner.Team.Points), 10), "§7-", loser.Team.Color()+strconv.FormatUint(uint64(loser.Team.Points), 10))
		for _, pl := range []*DuelPlayer{winner, loser} {
			pl.Player.SendTitle(t)
			pl.Player.Message(winner.Team.Color() + winner.Player.Name() + "§e has won!")
		}
		time.Sleep(3 * time.Second)
		for _, pl := range []*DuelPlayer{winner, loser} {
			if pl.Player != nil && !pl.Player.Dead() {
				utils.Srv.World().AddEntity(pl.Player)
				pl.Player.Teleport(utils.DefaultSpawn)
			}
		}
		_, _ = fmt.Fprintf(chat.Global, "§a%v §7won a match against §c%v!\n", winner.Player.Name(), p.Name())
		uuid := d.UUID
		Matches.Delete(d)
		go func() {
			if w, ok := utils.WorldManager.World(uuid); ok {
				_ = utils.WorldManager.UnloadWorld(w)
			}
			if err := os.RemoveAll("worlds/" + uuid); err != nil {
				panic(err)
			}
		}()
	}()
}

func (d *Duel) GetDuelPlayer(p *player.Player) *DuelPlayer {
	for _, pl := range d.Players {
		if pl.Player.Name() == p.Name() {
			return pl
		}
	}
	return nil
}

func (d *Duel) BroadcastMessage(message string) {
	for _, p := range d.Players {
		p.Player.Message(message)
	}
}

func (t *Team) Color() string {
	if t.Name == TeamRed {
		return "§c"
	}
	return "§b"
}

func (d *Duel) Score(p *DuelPlayer) {
	t := p.Team
	t.Points++
	if t.Points == 3 {
		for _, pl := range d.Players {
			if pl.Player.Name() != p.Player.Name() {
				d.RemovePlayer(p.Player)
				break
			}
		}
		return
	}
	go func() {
		if d == nil {
			return
		}
		ttl := title.New(t.Color() + p.Player.Name() + " scored")
		d.Waiting.Store(true)
		for _, pl := range d.Players {
			pl.Player.Message(t.Color() + p.Player.Name() + "§e has scored. §d" + strconv.FormatUint(uint64(3-t.Points), 10) + "§e more to win!")
			pl.Player.Teleport(pl.SpawnPoint)
			pl.Player.SetImmobile()
			d.SendKit(pl.Player, pl.Team.Color())
		}
		for i := 5; i > 0; i-- {
			if d == nil {
				return
			}
			for _, pl := range d.Players {
				pl.Player.SendTitle(ttl.WithSubtitle("§7Round starts in §a" + strconv.Itoa(i) + "s§7..."))
			}
			if i <= 3 {
				for _, pl := range d.Players {
					pl.Player.PlaySound(sound.Note{
						Instrument: sound.Piano(),
						Pitch:      0,
					})
				}
			}
			time.Sleep(1 * time.Second)
		}
		for _, pl := range d.Players {
			pl.Player.SetMobile()
			pl.Player.PlaySound(sound.Note{
				Instrument: sound.Piano(),
				Pitch:      1,
			})
		}
		d.Waiting.Store(false)
	}()
}

func (d *Duel) GetElapsedTime() string {
	t := d.ElapsedTime.Load()
	return strconv.FormatUint(uint64(t/60), 10) + ":" + strconv.FormatUint(uint64(t%60), 10)
}

func (d *Duel) BroadcastScoreboard() {
	if len(d.Players) > 1 {
		p1 := d.Players[0]
		p2 := d.Players[1]

		sb1 := scoreboard.New("§3Velvet Bridge")
		sb2 := scoreboard.New("§3Velvet Bridge")

		var rs string
		var bs string
		if p1.Team.Name == TeamRed {
			rs = p1.Team.FormatSquares()
			bs = p2.Team.FormatSquares()
		} else {
			rs = p2.Team.FormatSquares()
			bs = p1.Team.FormatSquares()
		}

		_, _ = sb1.WriteString("§bDuration: §3" + d.GetElapsedTime() + "\n\n§b[B] " + bs + "\n§c[R] " + rs + "\n\n§aKills: §f" + strconv.FormatUint(uint64(p1.Kills), 10) + "\n§aGoals: §f" + strconv.FormatUint(uint64(p1.Team.Points), 10) + "\n\n§aYour Ping: §f" + strconv.Itoa(int(p1.Player.Latency().Seconds())) + "\n§cTheir Ping: §f" + strconv.Itoa(int(p2.Player.Latency().Seconds())) + "\n\n§bvelvetpractice.live")
		_, _ = sb2.WriteString("§bDuration: §3" + d.GetElapsedTime() + "\n\n§b[B] " + bs + "\n§c[R] " + rs + "\n\n§aKills: §f" + strconv.FormatUint(uint64(p2.Kills), 10) + "\n§aGoals: §f" + strconv.FormatUint(uint64(p2.Team.Points), 10) + "\n\n§aYour Ping: §f" + strconv.Itoa(int(p2.Player.Latency().Seconds())) + "\n§cTheir Ping: §f" + strconv.Itoa(int(p1.Player.Latency().Seconds())) + "\n\n§bvelvetpractice.live")

		p1.Player.SendScoreboard(sb1)
		p2.Player.SendScoreboard(sb2)
	}
}

func (t *Team) FormatSquares() string {
	a := t.Color()
	switch t.Points {
	case 1:
		return a + "█" + "§7██"
	case 2:
		return a + "██" + "§7█"
	case 3:
		return a + "███"
	}
	return "§7███"
}

func (d *Duel) Over() bool {
	for _, pl := range d.Players {
		if pl.Team.Points == 3 {
			return true
		}
	}
	return false
}
