package protocol

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"

	"github.com/chaihaobo/claw-machine-protocol/codec"
	"github.com/chaihaobo/claw-machine-protocol/commands"
)

var (
	ErrDeviceNotfound = errors.New("device not found")
)

type (
	APICaller interface {
		// AllDevices returns all devices
		AllDevices(ctx context.Context) ([]*Device, error)
		// QueryStatus returns the status of the device
		QueryStatus(ctx context.Context, deviceID string, input *commands.QueryStatusInput) (*commands.QueryStatusOutput, error)
		// UpScore to one device
		UpScore(ctx context.Context, deviceID string, input *commands.UpScoreInput) (*commands.UpScoreOutput, error)
		// Control  one device
		Control(ctx context.Context, deviceID string, input *commands.ControlInput) error
		// LightControl one device
		LightControl(ctx context.Context, deviceID string, input *commands.LightControlInput) error
		// QuerySetting query the setting of the device
		QuerySetting(ctx context.Context, deviceID string) (*commands.QuerySettingOutput, error)
		// SetSetting set the setting of the device
		SetSetting(ctx context.Context, deviceID string, setting *commands.SetSettingInput) (*commands.SetSettingOutput, error)
		// QueryAccount query the account of the device
		QueryAccount(ctx context.Context, deviceID string) (*commands.QueryAccountOutput, error)
		// ResetPlays reset the play count of the device
		ResetPlays(ctx context.Context, deviceID string) error
		// Reboot one device
		Reboot(ctx context.Context, deviceID string) error
	}

	apiCaller struct {
		deviceRegistry Registry
		token          string
	}
)

func (a apiCaller) Reboot(ctx context.Context, deviceID string) error {
	if err := a.call(ctx, deviceID, &commands.RebootInput{}, nil); err != nil {
		return err
	}
	return nil
}

func (a apiCaller) ResetPlays(ctx context.Context, deviceID string) error {
	if err := a.call(ctx, deviceID, &commands.ResetPlays{}, nil); err != nil {
		return err
	}
	return nil
}

func (a apiCaller) QueryAccount(ctx context.Context, deviceID string) (*commands.QueryAccountOutput, error) {
	output := new(commands.QueryAccountOutput)
	if err := a.call(ctx, deviceID, &commands.QueryAccountInput{}, output); err != nil {
		return nil, err
	}
	return output, nil
}

func (a apiCaller) SetSetting(ctx context.Context, deviceID string, setting *commands.SetSettingInput) (*commands.SetSettingOutput, error) {
	output := new(commands.SetSettingOutput)
	if err := a.call(ctx, deviceID, setting, output); err != nil {
		return nil, err
	}
	return output, nil
}

func (a apiCaller) QuerySetting(ctx context.Context, deviceID string) (*commands.QuerySettingOutput, error) {
	output := new(commands.QuerySettingOutput)
	if err := a.call(ctx, deviceID, &commands.QuerySettingInput{}, output); err != nil {
		return nil, err
	}
	return output, nil
}

func (a apiCaller) LightControl(ctx context.Context, deviceID string, input *commands.LightControlInput) error {
	if err := a.call(ctx, deviceID, input, nil); err != nil {
		return err
	}
	return nil
}

func (a apiCaller) AllDevices(ctx context.Context) ([]*Device, error) {
	return a.deviceRegistry.DiscoveryAll()
}

func (a apiCaller) Control(ctx context.Context, deviceID string, input *commands.ControlInput) error {
	if err := a.call(ctx, deviceID, input, nil); err != nil {
		return err
	}
	return nil
}

func (a apiCaller) UpScore(ctx context.Context, deviceID string, input *commands.UpScoreInput) (*commands.UpScoreOutput, error) {
	output := new(commands.UpScoreOutput)
	if err := a.call(ctx, deviceID, input, output); err != nil {
		return nil, err
	}
	return output, nil
}

func (a apiCaller) QueryStatus(ctx context.Context, deviceID string, input *commands.QueryStatusInput) (*commands.QueryStatusOutput, error) {
	output := new(commands.QueryStatusOutput)
	if err := a.call(ctx, deviceID, input, output); err != nil {
		return nil, err
	}
	return output, nil
}

func (a apiCaller) call(ctx context.Context, deviceID string, input codec.Packet, output codec.PacketUnmarshaler) error {
	device, err := a.deviceRegistry.Discovery(deviceID)
	if err != nil {
		return err
	}
	if device == nil {
		return ErrDeviceNotfound
	}
	body, err := codec.Encode(input)
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://%s", device.ExternalControlAddr), bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/octet-stream")
	request.Header.Add(ExternalControlTokenHeader, a.token)
	request.Header.Add(ExternalControlDeviceIDHeader, deviceID)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("api call failed:invalid status code: %d", response.StatusCode))
	}
	defer response.Body.Close()
	responsePacket, err := codec.Decode(response.Body)
	if err != nil {
		return err
	}
	if output == nil {
		return nil
	}
	if err := output.Unmarshal(responsePacket); err != nil {
		return err
	}
	return nil
}

func NewAPICaller(redisClient *redis.Client, token string) APICaller {
	return &apiCaller{
		deviceRegistry: NewRegistry(redisClient),
		token:          token,
	}
}
