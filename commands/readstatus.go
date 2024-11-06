package commands

import (
	"github.com/chaihaobo/claw-machine-protocol/codec"
)

const (
	StatusTypeAutoStart StatusType = 0x07
	StatusTypeAutoGrab  StatusType = 0x08
	StatusTypeGameOver  StatusType = 0x09
)

var statusTypeValues = map[StatusType]string{
	StatusTypeAutoStart: "自动开始",
	StatusTypeAutoGrab:  "自动抓取",
	StatusTypeGameOver:  "游戏结束",
}

type (
	ReadStatusOutput struct {
		// StatusType 状态类型
		StatusType StatusType
		// 是否抓中礼品
		Grabbed bool
	}

	StatusType byte
)

func (s StatusType) Value() string {
	return statusTypeValues[s]
}

func (r *ReadStatusOutput) Unmarshal(p codec.Packet) error {
	data := p.Data()
	if len(data) != 2 {
		return nil
	}
	r.StatusType = StatusType(data[0])
	r.Grabbed = data[1] == 0x01
	return nil
}
