package codec

import (
	"errors"
	"io"
)

var (
	ErrInvalidFrameHead = errors.New("invalid frame head")
)

type (
	packetEncoder packet
	packetDecoder struct {
		reader io.Reader
	}
)

func (p packetEncoder) encode() ([]byte, error) {
	var (
		index  = p.index
		cmd    = p.cmd
		data   = p.data
		length = byte(3 + len(data))
	)
	payload := []byte{length, index, cmd}
	payload = append(payload, data...)
	checkSum := calculateCheckSum(payload...)
	rawPacket := []byte{FrameHead, length, index, cmd}
	rawPacket = append(rawPacket, data...)
	rawPacket = append(rawPacket, checkSum, FrameTail)
	return rawPacket, nil
}

func (p packetDecoder) decode() (Packet, error) {
	frameHeadBytes := make([]byte, 1)
	_, err := p.reader.Read(frameHeadBytes)
	if err != nil {
		return nil, err
	}
	if frameHeadBytes[0] != FrameHead {
		return nil, ErrInvalidFrameHead
	}
	// 剩下的字节长度
	length := make([]byte, 1)
	_, err = p.reader.Read(length)
	if err != nil {
		return nil, err
	}
	mainPayload := make([]byte, length[0])
	// 读取剩下的字节 Index + CMD + Data + Check
	_, err = p.reader.Read(mainPayload)
	if err != nil {
		return nil, err
	}
	// 读取尾部字节
	tail := make([]byte, 1)
	_, err = p.reader.Read(tail)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	if tail[0] != FrameTail {
		return nil, ErrInvalidFrameHead
	}

	return parsePacket(mainPayload)
}

func Encode(packet Packet) ([]byte, error) {
	return packetEncoder(newPacket(packet)).encode()
}

func Decode(reader io.Reader) (Packet, error) {
	return packetDecoder{reader: reader}.decode()
}
