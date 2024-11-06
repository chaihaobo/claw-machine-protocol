package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/chaihaobo/claw-machine-protocol"
	"github.com/chaihaobo/claw-machine-protocol/commands"
)

func main() {
	opt := &redis.Options{
		Network:  "tcp",
		Addr:     "0",
		Password: "0",
		DB:       3,
	}
	client := redis.NewClient(opt)
	apicaller := protocol.NewAPICaller(client, "")
	ctx := context.Background()
	devices, _ := apicaller.AllDevices(ctx)
	for _, device := range devices {
		println(device.ID)
	}
	if len(devices) == 0 {
		return
	}
	output, err := apicaller.QueryStatus(ctx, devices[0].ID, &commands.QueryStatusInput{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v \n", output)
	// 重置局数
	if err := apicaller.ResetPlays(ctx, devices[0].ID); err != nil {
		panic(err)
	}
	// 开灯
	err = apicaller.LightControl(ctx, devices[0].ID, &commands.LightControlInput{true})
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 2)
	// 关灯
	apicaller.LightControl(ctx, devices[0].ID, &commands.LightControlInput{false})
	// 加币
	apicaller.UpScore(ctx, devices[0].ID, &commands.UpScoreInput{1, 2})
	setting, err := apicaller.QuerySetting(ctx, devices[0].ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v \n", setting)
	setting.CoinsPerPlay = 201 // 2币1玩
	result, err := apicaller.SetSetting(ctx, devices[0].ID, &commands.SetSettingInput{setting.Setting})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v \n", result)
	// 查账
	accountOutput, err := apicaller.QueryAccount(ctx, devices[0].ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v \n", accountOutput)

	apicaller.Control(ctx, devices[0].ID, &commands.ControlInput{commands.ControlTypeRight, commands.StrengthStrong})
	time.Sleep(time.Second * 2)
	apicaller.Control(ctx, devices[0].ID, &commands.ControlInput{commands.ControlTypeStop, commands.StrengthStrong})
	apicaller.Control(ctx, devices[0].ID, &commands.ControlInput{commands.ControlTypeBackward, commands.StrengthStrong})
	time.Sleep(time.Second * 2)
	apicaller.Control(ctx, devices[0].ID, &commands.ControlInput{commands.ControlTypeStop, commands.StrengthStrong})
	apicaller.Control(ctx, devices[0].ID, &commands.ControlInput{commands.ControlTypeLeft, commands.StrengthStrong})
	time.Sleep(time.Second * 2)
	apicaller.Control(ctx, devices[0].ID, &commands.ControlInput{commands.ControlTypeStop, commands.StrengthStrong})
	apicaller.Control(ctx, devices[0].ID, &commands.ControlInput{commands.ControlTypeForward, commands.StrengthStrong})
	time.Sleep(time.Second * 2)
	apicaller.Control(ctx, devices[0].ID, &commands.ControlInput{commands.ControlTypeStop, commands.StrengthStrong})
	apicaller.Control(ctx, devices[0].ID, &commands.ControlInput{commands.ControlTypeGrab, commands.StrengthStrong})

	time.Sleep(time.Second * 2)
	//if err := apicaller.Reboot(ctx, devices[0].ID); err != nil {
	//	panic(err)
	//}
	//println(output.DeviceID)
}
