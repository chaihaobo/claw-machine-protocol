package commands

import (
	"github.com/chaihaobo/claw-machine-protocol/codec"
)

type (
	LightControlInput struct {
		Open bool
	}
)

func (l LightControlInput) Index() codec.Index {
	return codec.IndexBoxHost
}

func (l LightControlInput) CMD() byte {
	return CommandLightControl
}

func (l LightControlInput) Data() []byte {
	var data byte = 0x00
	if l.Open {
		data = 0x01
	}
	return []byte{0x01, data}
}
