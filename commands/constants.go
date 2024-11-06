package commands

// CommandQueryLink 查询状态
// CommandUpScore 上分
// CommandControl 控制
// CommandLightControl 灯光控制
// CommandQuerySetting 查询设置
// CommandSetSetting 设置
// CommandQueryAccount 查账
// CommandResetPlays 重置局数
const (
	CommandQueryLink    = 0x01
	CommandUpScore      = 0x03
	CommandControl      = 0x50
	CommandLightControl = 0x60
	CommandQuerySetting = 0x05
	CommandSetSetting   = 0x06
	CommandQueryAccount = 0x30
	CommandResetPlays   = 0x32
	CommandReboot       = 0x39
)
