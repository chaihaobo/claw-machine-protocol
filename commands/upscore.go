package commands

import (
	"github.com/chaihaobo/claw-machine-protocol/codec"
)

type (
	UpScoreInput struct {
		IncrementNo int8
		CoinCount   int16
	}
	UpScoreOutput struct {
		Success bool
	}
)

func (q UpScoreInput) Index() codec.Index {
	return codec.IndexMainboard
}

func (q *UpScoreOutput) Unmarshal(p codec.Packet) error {
	bytes := p.Data()
	q.Success = bytes[0] == 0x01
	return nil
}

func (q UpScoreInput) CMD() byte {
	//aa 08 01 03 01 01 00 00 00  0a dd
	return CommandUpScore
}

func (q UpScoreInput) Data() []byte {
	return []byte{
		byte(q.IncrementNo),
		byte(q.CoinCount),
		byte(q.CoinCount >> 8),
		0,
		0,
	}
}
