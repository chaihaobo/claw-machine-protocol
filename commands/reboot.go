package commands

import (
	"github.com/chaihaobo/claw-machine-protocol/codec"
)

type (
	RebootInput struct {
	}
)

func (r RebootInput) Index() codec.Index {
	return codec.IndexBoxHost
}

func (r RebootInput) CMD() byte {
	return CommandReboot
}

func (r RebootInput) Data() []byte {
	return []byte{}
}
