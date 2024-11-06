package protocol

import (
	"context"
	"encoding/hex"
	"net"
	"sync"
	"time"

	"github.com/chaihaobo/claw-machine-protocol/codec"
	"github.com/chaihaobo/claw-machine-protocol/commands"
)

// DefaultReadStatusPacketTimeout 读取设备状态包的超时时间
const (
	DefaultReadStatusPacketTimeout = time.Second
	DefaultReadStatusInterval      = time.Second
)

var (
	gameOverAckCommand = []byte{0xaa, 0x04, 0x01, 0x50, 0x09, 0x5C, 0xdd}
)

type Connection struct {
	deviceID           string
	con                net.Conn
	mutex              sync.Mutex
	logger             Logger
	statusEventHandler statusEventHandler
}

func (c *Connection) RemoteAddr() string {
	return c.con.RemoteAddr().String()
}

func (c *Connection) Send(input codec.Packet, output codec.PacketUnmarshaler) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	inputPayload, err := codec.Encode(input)
	if err != nil {
		return err
	}
	c.logger.Printf("send to device %s: %s\n", c.deviceID, hex.EncodeToString(inputPayload))
	_, err = c.con.Write(inputPayload)
	if err != nil {
		return err
	}
	outputPacket, err := codec.Decode(c.con)
	if err != nil {
		return err
	}
	if output != nil {
		return output.Unmarshal(outputPacket)
	}
	return nil
}

func (c *Connection) pollingReadStatus() {
	ticker := time.NewTicker(DefaultReadStatusInterval)
	for {
		select {
		case <-ticker.C:
			if err := c.handleReadStatus(); err != nil {
				c.logger.Printf("read status from device %s failed: %s, then will stop polling status\n", c.deviceID, err)
				return
			}
		}
	}
}

func (c *Connection) handleReadStatus() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.con.SetReadDeadline(time.Now().Add(DefaultReadStatusPacketTimeout))
	defer c.con.SetReadDeadline(time.Time{})
	outputPacket, err := codec.Decode(c.con)
	if err != nil && !isDeadlineExceededError(err) {
		c.logger.Printf("read status from device %s failed: %s\n", c.deviceID, err)
		return err
	}
	if isDeadlineExceededError(err) {
		return nil
	}
	// 处理主动上报状态
	if outputPacket.CMD() == 0x50 && outputPacket.Index() == codec.IndexMainboard {
		var statusEvent commands.ReadStatusOutput
		_ = statusEvent.Unmarshal(outputPacket)

		c.logger.Printf("receive status from device %s: %s grabbed:%t \n", c.deviceID,
			statusEvent.StatusType.Value(),
			statusEvent.Grabbed,
		)
		if statusEventHandler := c.statusEventHandler; statusEventHandler != nil {
			statusEventHandler(context.Background(), c.deviceID, &statusEvent)
			// if receive the game over event, then ask this message
			c.con.Write(gameOverAckCommand)
		}
	}

	return nil
}
func isDeadlineExceededError(err error) bool {
	if ne, ok := err.(*net.OpError); ok && ne.Timeout() {
		return true
	}
	return false
}

func (c *Connection) Close() error {
	return c.con.Close()
}

func NewConnection(logger Logger, con net.Conn) *Connection {
	c := &Connection{
		con:    con,
		logger: logger,
	}
	go c.pollingReadStatus()
	return c
}
