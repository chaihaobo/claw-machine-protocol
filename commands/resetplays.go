package commands

import "github.com/chaihaobo/claw-machine-protocol/codec"

type (
	ResetPlays struct {
	}
)

func (r ResetPlays) Index() codec.Index {
	return codec.IndexBoxHost
}

func (r ResetPlays) CMD() byte {
	return CommandResetPlays
}

func (r ResetPlays) Data() []byte {
	return []byte{}
}
