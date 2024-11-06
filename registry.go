package protocol

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	deviceCacheKey = "claw:machine:devices"
)

type (
	// Registry 注册中心. 用于注册设备和发现设备
	Registry interface {
		Register(deviceID string, device *Device) error
		Discovery(deviceID string) (*Device, error)
		DiscoveryAll() ([]*Device, error)
	}

	Device struct {
		ID   string `json:"id"`
		Addr string `json:"addr"`
		// ExternalControlAddr 外部控制地址 将来注册到注册中心后 客户端会从注册中心获取外部控制地址进行调用
		ExternalControlAddr string    `json:"external_control_host"`
		LastHeartBeatAt     time.Time `json:"last_heart_beat_at"`
	}
	registry struct {
		redisClient *redis.Client
	}
)

func (r registry) DiscoveryAll() ([]*Device, error) {
	devices := make([]*Device, 0)
	if err := r.redisClient.HVals(context.Background(), deviceCacheKey).ScanSlice(&devices); err != nil {
		return nil, err
	}
	return devices, nil
}

func NewRegistry(redisClient *redis.Client) Registry {
	return &registry{redisClient}
}

func (r registry) Register(deviceID string, device *Device) error {
	return r.redisClient.HSet(context.Background(), deviceCacheKey, deviceID, device).Err()
}

func (r registry) Discovery(deviceID string) (*Device, error) {
	var device Device
	err := r.redisClient.HGet(context.Background(), deviceCacheKey, deviceID).Scan(&device)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	return &device, nil
}

func NewDevice(id, addr, externalControlHost string) *Device {
	return &Device{
		ID:                  id,
		Addr:                addr,
		ExternalControlAddr: externalControlHost,
		LastHeartBeatAt:     time.Now(),
	}
}

func (d Device) MarshalBinary() (data []byte, err error) {
	return json.Marshal(d)
}

func (d *Device) UnmarshalBinary(data []byte) (err error) {
	return json.Unmarshal(data, d)
}
