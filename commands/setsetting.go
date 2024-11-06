package commands

import (
	"github.com/chaihaobo/claw-machine-protocol/codec"
)

type (
	SetSettingInput struct {
		Setting
	}
	SetSettingOutput struct {
		Result bool
	}
)

func (s *SetSettingOutput) Unmarshal(p codec.Packet) error {
	s.Result = p.Data()[0] == 0x01
	return nil
}
