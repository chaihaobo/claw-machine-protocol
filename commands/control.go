package commands

import (
	"github.com/chaihaobo/claw-machine-protocol/codec"
)

const (
	ControlTypeForward  ControlType = 0x01 // 前
	ControlTypeBackward ControlType = 0x02 // 后
	ControlTypeLeft     ControlType = 0x03 // 左
	ControlTypeRight    ControlType = 0x04 // 右
	ControlTypeStop     ControlType = 0x05 // 停
	ControlTypeGrab     ControlType = 0x06 // 抓
)

const (
	StrengthWeak   Strength = 0x00 // 弱
	StrengthStrong Strength = 0x01 // 强
)

type (
	ControlInput struct {
		ControlType ControlType
		Strength    Strength
	}
	ControlType int8
	Strength    int8
)

func (c ControlInput) Index() codec.Index {
	return codec.IndexBoxHost
}

func (c ControlInput) CMD() byte {
	return CommandControl
}

func (c ControlInput) Data() []byte {
	return []byte{byte(c.ControlType), byte(c.Strength)}
}
