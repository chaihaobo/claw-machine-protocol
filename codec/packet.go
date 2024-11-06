package codec

import (
	"errors"
)

var (
	ErrPacketCheckSumInvalid = errors.New("invalid packet check sum")
)

type MutablePacket interface {
	Packet
	PacketUnmarshaler
}

type Packet interface {
	Index() Index
	CMD() byte
	Data() []byte
}

type PacketUnmarshaler interface {
	Unmarshal(p Packet) error
}

type packet struct {
	index byte
	cmd   byte
	data  []byte
}

func (p *packet) Unmarshal(src Packet) error {
	p.index = byte(src.Index())
	p.cmd = src.CMD()
	p.data = src.Data()
	return nil
}

func (p packet) Index() Index {
	return Index(p.index)
}

func (p packet) CMD() byte {
	return p.cmd
}

func (p packet) Data() []byte {
	return p.data
}

func newPacket(p Packet) packet {
	if !IsValidIndex(p.Index()) {
		panic("invalid index")
	}
	if p == nil {
		panic("nil packet")
	}
	return packet{
		index: byte(p.Index()),
		cmd:   p.CMD(),
		data:  p.Data(),
	}
}

// parsePacket 解析数据包
// Index + CMD + Data + Check
func parsePacket(data []byte) (packet, error) {
	length := len(data)
	checksum := data[length-1]
	index := data[0]
	cmd := data[1]
	cmdData := data[2 : length-1]

	payload := []byte{byte(length), index, cmd}
	payload = append(payload, cmdData...)
	// 如果是灯光控制响应，就不检查checksum了。因为他们厂家的协议设计的太烂了!
	if index == byte(IndexMainboard) && cmd == 0x60 && len(cmdData) > 0 && cmdData[0] == 0x01 {
		return packet{index, cmd, cmdData}, nil
	}
	caculatedCheckSum := calculateCheckSum(payload...)
	if checksum != caculatedCheckSum {
		return packet{}, ErrPacketCheckSumInvalid
	}
	return packet{index, cmd, cmdData}, nil
}

func EmptyPacket() MutablePacket {
	return &packet{}
}

func Unmarshal(src Packet, dst PacketUnmarshaler) error {
	return dst.Unmarshal(src)
}
