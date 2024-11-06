package codec

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type leftCommandPacket int

func (l leftCommandPacket) CMD() byte {
	return 0x50
}

func (l leftCommandPacket) Data() []byte {
	return []byte{0x01, 0x01}
}

type downCommandPacket int

func (d downCommandPacket) CMD() byte {
	return 0x50
}

func (d downCommandPacket) Data() []byte {
	return []byte{0x02, 0x01}
}

func TestEncode(t *testing.T) {
	testcases := []struct {
		name    string
		command Packet
		want    []byte
	}{
		{
			name:    "when run left command",
			command: leftCommandPacket(0),
			want:    []byte{0xaa, 0x05, 0x01, 0x50, 0x01, 0x01, 0x54, 0xdd},
		},
		{
			name:    "when run down command",
			command: downCommandPacket(0),
			want:    []byte{0xaa, 0x05, 0x01, 0x50, 0x02, 0x01, 0x57, 0xdd},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			rawPacket, err := Encode(IndexBoxHost, tc.command)
			assert.NoError(t, err)
			assert.Equal(t, tc.want, rawPacket)
		})
	}

}

func TestDecode(t *testing.T) {
	testcases := []struct {
		name   string
		reader io.Reader
		want   Packet
	}{
		{
			name:   "when decode let command success",
			reader: bytes.NewReader([]byte{0xaa, 0x05, 0x01, 0x50, 0x01, 0x01, 0x54, 0xdd}),
			want:   leftCommandPacket(0),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			packet, err := Decode(tc.reader)
			assert.NoError(t, err)
			assert.Equal(t, tc.want.CMD(), packet.CMD())
			assert.Equal(t, tc.want.Data(), packet.Data())

		})
	}

}
