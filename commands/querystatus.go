package commands

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/chaihaobo/claw-machine-protocol/codec"
)

type (
	QueryStatusInput struct {
	}
	QueryStatusOutput struct {
		DeviceID string
	}
)

func (q QueryStatusInput) Index() codec.Index {
	return codec.IndexBoxHost
}

func (q *QueryStatusOutput) Unmarshal(p codec.Packet) error {
	bytes := p.Data()
	q.DeviceID = hex.EncodeToString(bytes)
	return nil
}

func (q QueryStatusInput) CMD() byte {
	return CommandQueryLink
}

func (q QueryStatusInput) Data() []byte {
	return GenerateDeviceID()
}

// GenerateDeviceID 生成一个 8 字节的随机设备唯一码
func GenerateDeviceID() []byte {
	id := make([]byte, 8) // 创建一个 8 字节的切片
	_, err := rand.Read(id)
	if err != nil {
		return nil
	}
	return id
}
