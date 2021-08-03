package exts

import (
	"VBridge/duels"
	"VBridge/utils"
	"go.uber.org/atomic"
	"time"
)

var Matches atomic.Int32
var Online atomic.Int32

func TrackData(){
	go func(){
		for {
			Matches.Store(int32(len(duels.Matches.List)))
			Online.Store(int32(len(utils.Srv.Players())))
			time.Sleep(1 * time.Second)
		}
	}()
}
