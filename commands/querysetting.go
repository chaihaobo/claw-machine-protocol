package commands

import (
	"github.com/chaihaobo/claw-machine-protocol/codec"
)

// PayoutMode 定义出奖模式的枚举类型
type PayoutMode int8

const (
	PayoutModeNoProbability PayoutMode = iota // 0 无概率
	PayoutModeRandom                          // 1 随机模式
	PayoutModeFixed                           // 2 固定模式
	PayoutModeGuanXing                        // 3 冠兴模式
)

func (p PayoutMode) String() string {
	switch p {
	case PayoutModeNoProbability:
		return "无概率"
	case PayoutModeRandom:
		return "随机模式"
	case PayoutModeFixed:
		return "固定模式"
	case PayoutModeGuanXing:
		return "冠兴模式"
	default:
		return "未知模式"
	}
}

type GuaranteeBonusMode int8

const (
	GuaranteeBonusModeGame         GuaranteeBonusMode = iota // 0 送游戏
	GuaranteeBonusModePrize                                  // 1 送中奖
	GuaranteeBonusModeGameAndPrize                           // 2 送游戏和中奖
)

func (g GuaranteeBonusMode) String() string {
	switch g {
	case GuaranteeBonusModeGame:
		return "送游戏"
	case GuaranteeBonusModePrize:
		return "送中奖"
	case GuaranteeBonusModeGameAndPrize:
		return "送游戏和中奖"
	default:
		return "未知模式"
	}
}

type (
	Setting struct {
		// CoinsPerPlay 完一把需要几个币
		CoinsPerPlay int16
		// GameTime 每把的游戏时间
		GameTime int8
		// PayoutMode 出奖模式
		PayoutMode PayoutMode
		// PayoutRate 出奖概率 范围1~9999
		PayoutRate int16
		// EnableAirGrab 是否开启空中抓物
		EnableAirGrab bool
		// continuousCoinBonus 连续投币赠送
		ContinuousCoinBonus int8
		// guaranteedCaptureCount 保夹次数
		GuaranteedCaptureCount int8
		// GuaranteeBonusMode 保夹赠送模式
		GuaranteeBonusMode GuaranteeBonusMode
		// StrongGripVoltage 强抓力电压
		StrongGripVoltage int16
		// MediumGripVoltage 中抓力电压
		MediumGripVoltage int16
		// WeakGripVoltage 弱抓力电压
		WeakGripVoltage int16
		// PrizeVoltage 中奖电压
		PrizeVoltage int16
		// StrongGripTime 强抓时间
		StrongGripTime int16
		// ReleaseGripTime 释放抓力时间
		ReleaseGripTime int8
		// ForwardBackwardSpeed 前后速度
		ForwardBackwardSpeed int8
		// LeftRightSpeed 左右速度
		LeftRightSpeed int8
		// UpDownSpeed 上下速度
		UpDownSpeed int8
		// LineReleaseTime 放线长度
		LineReleaseTime int16
		// PrizeDropHeight 礼品下放高度
		PrizeDropHeight int16
		// SwingGripLength 摆动抓力长度
		SwingGripLength int8
		// SwingGripProtection 摆动抓力保护
		SwingGripProtection int8
		// SwingGripVoltage 甩抓电压
		SwingGripVoltage int16
		// LiftProtection 上拉保护
		LiftProtection int8
		// CraneSelfRescueTime 天车自救时间
		CraneselfRescueTime int8
		// LowerGripDelay 下抓延时
		LowerGripDelay int16
		// GripObjectDelay 抓取物体延时
		GripObjectDelay int16
		// LiftStopDelay 上拉停止延时
		LiftStopDelay int16
		// JoystickDelay 摇杆延时
		JoystickDelay int16
		// GripObjectFirstRetrieval 抓物二收
		GripObjectSecondRetrieval int8
		// StandbyMusicSwitch 待机音乐开关
		StandbyMusicSwitch bool
		// VolumeLevel 音量大小
		VolumeLevel int8
		// StandbyMusicSelection 待机音乐
		StandbyMusicSelection int8
		// GameMusicSelection 游戏音乐
		GameMusicSelection int8
		// ProbabilityQueueAutoReset 概率队列自动重置 范围0~255分钟
		ProbabilityQueueAutoReset int8
	}
	QuerySettingInput struct {
	}
	QuerySettingOutput struct {
		Setting
	}
)

func (s *Setting) Index() codec.Index {
	return codec.IndexBoxHost
}

func (s *Setting) CMD() byte {
	return CommandSetSetting
}

func (s *Setting) Data() []byte {
	return []byte{
		byte(s.CoinsPerPlay),
		byte(s.CoinsPerPlay >> 8),
		byte(s.GameTime),
		byte(s.PayoutMode),
		byte(s.PayoutRate),
		byte(s.PayoutRate >> 8),
		boolToByte(s.EnableAirGrab),
		byte(s.ContinuousCoinBonus),
		byte(s.GuaranteedCaptureCount),
		byte(s.GuaranteeBonusMode),
		byte(s.StrongGripVoltage),
		byte(s.StrongGripVoltage >> 8),
		byte(s.MediumGripVoltage),
		byte(s.MediumGripVoltage >> 8),
		byte(s.WeakGripVoltage),
		byte(s.WeakGripVoltage >> 8),
		byte(s.PrizeVoltage),
		byte(s.PrizeVoltage >> 8),
		byte(s.StrongGripTime),
		byte(s.StrongGripTime >> 8),
		byte(s.ReleaseGripTime),
		byte(s.ForwardBackwardSpeed),
		byte(s.LeftRightSpeed),
		byte(s.UpDownSpeed),
		byte(s.LineReleaseTime),
		byte(s.LineReleaseTime >> 8),
		byte(s.PrizeDropHeight),
		byte(s.PrizeDropHeight >> 8),
		byte(s.SwingGripLength),
		byte(s.SwingGripProtection),
		byte(s.SwingGripVoltage),
		byte(s.SwingGripVoltage >> 8),
		byte(s.LiftProtection),
		byte(s.CraneselfRescueTime),
		byte(s.LowerGripDelay),
		byte(s.LowerGripDelay >> 8),
		byte(s.GripObjectDelay),
		byte(s.GripObjectDelay >> 8),
		byte(s.LiftStopDelay),
		byte(s.LiftStopDelay >> 8),
		byte(s.JoystickDelay),
		byte(s.JoystickDelay >> 8),
		byte(s.GripObjectSecondRetrieval),
		boolToByte(s.StandbyMusicSwitch),
		byte(s.VolumeLevel),
		byte(s.StandbyMusicSelection),
		byte(s.GameMusicSelection),
		byte(s.ProbabilityQueueAutoReset),
	}
}

// bool to byte
func boolToByte(b bool) byte {
	if b {
		return 0x01
	}
	return 0x00
}

func (s *Setting) Unmarshal(p codec.Packet) error {
	data := p.Data()
	if len(data) != 48 {
		return nil
	}
	// 解析
	s.CoinsPerPlay = int16(data[1]<<8) | int16(data[0])
	s.GameTime = int8(data[2])
	s.PayoutMode = PayoutMode(data[3])
	s.PayoutRate = int16(data[5]<<8) | int16(data[4])
	s.EnableAirGrab = data[6] == 0x01
	s.ContinuousCoinBonus = int8(data[7])
	s.GuaranteedCaptureCount = int8(data[8])
	s.GuaranteeBonusMode = GuaranteeBonusMode(data[9])
	s.StrongGripVoltage = int16(data[11]<<8) | int16(data[10])
	s.MediumGripVoltage = int16(data[13]<<8) | int16(data[12])
	s.WeakGripVoltage = int16(data[15]<<8) | int16(data[14])
	s.PrizeVoltage = int16(data[17]<<8) | int16(data[16])
	s.StrongGripTime = int16(data[19]<<8) | int16(data[18])
	s.ReleaseGripTime = int8(data[20])
	s.ForwardBackwardSpeed = int8(data[21])
	s.LeftRightSpeed = int8(data[22])
	s.UpDownSpeed = int8(data[23])
	s.LineReleaseTime = int16(data[25])<<8 | int16(data[24])
	s.PrizeDropHeight = int16(data[27]<<8) | int16(data[26])
	s.SwingGripLength = int8(data[28])
	s.SwingGripProtection = int8(data[29])
	s.SwingGripVoltage = int16(data[31]<<8) | int16(data[30])
	s.LiftProtection = int8(data[32])
	s.CraneselfRescueTime = int8(data[33])
	s.LowerGripDelay = int16(data[35]<<8) | int16(data[35])
	s.GripObjectDelay = int16(data[37]<<8 | data[36])
	s.LiftStopDelay = int16(data[39]<<8) | int16(data[38])
	s.JoystickDelay = int16(data[41]<<8) | int16(data[40])
	s.GripObjectSecondRetrieval = int8(data[42])
	s.StandbyMusicSwitch = data[43] == 0x01
	s.VolumeLevel = int8(data[44])
	s.StandbyMusicSelection = int8(data[45])
	s.GameMusicSelection = int8(data[46])
	s.ProbabilityQueueAutoReset = int8(data[47])
	return nil
}

func (q QuerySettingInput) Index() codec.Index {
	return codec.IndexBoxHost
}

func (q QuerySettingInput) CMD() byte {
	return CommandQuerySetting
}

func (q QuerySettingInput) Data() []byte {
	return []byte{}
}
