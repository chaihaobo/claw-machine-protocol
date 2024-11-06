package commands

import (
	"github.com/chaihaobo/claw-machine-protocol/codec"
)

type (
	QueryAccountInput struct {
	}
	QueryAccountOutput struct {
		// TotalCoinCount 总币数
		TotalCoinCount int32
		// TotalPrizeCount 出奖总数
		TotalPrizeCount int32
	}
)

func (q *QueryAccountOutput) Unmarshal(p codec.Packet) error {
	data := p.Data()
	if len(data) != 10 {
		return nil
	}
	q.TotalCoinCount = int32(data[3])<<24 | int32(data[2])<<16 | int32(data[1])<<8 | int32(data[0])
	q.TotalPrizeCount = int32(data[7])<<24 | int32(data[6])<<16 | int32(data[5])<<8 | int32(data[4])
	return nil
}

func (q QueryAccountInput) Index() codec.Index {
	return codec.IndexBoxHost
}

func (q QueryAccountInput) CMD() byte {
	return CommandQueryAccount
}

func (q QueryAccountInput) Data() []byte {
	return []byte{}
}
