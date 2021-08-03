package cblock

import (
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

type EndPortal struct {

}

func (EndPortal) EncodeItem() (name string, meta int16) {
	return "minecraft:end_portal", 0
}

// EncodeBlock ...
func (b EndPortal) EncodeBlock() (name string, properties map[string]interface{}) {
	return "minecraft:end_portal", nil
}

func (b EndPortal) Model() world.BlockModel {
	return model.Solid{}
}

func (EndPortal) Hash() uint64 {
	return 3000
}