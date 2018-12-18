package protocol

type Property int

const UnknownProperty = -1
const TunerChannel Property = 0
const TunerChannelMap Property = 1
const TunerFilter Property = 2
const TunerProgram Property = 3
const TunerTarget Property = 4
const TunerStatus Property = 5
const TunerStreamInfo Property = 6
const TunerDebug Property = 7
const TunerLockKey Property = 8
const IrTarget Property = 9
const LineupLocation Property = 10
const SysModel Property = 11
const SysFeatures Property = 12
const SysVersion Property = 13
const SysCopyright Property = 14
const SysDebug Property = 15
const SysHwModel Property = 16

var Properties []Property = []Property{
	UnknownProperty,
	TunerChannel,
	TunerChannelMap,
	TunerFilter,
	TunerProgram,
	TunerTarget,
	TunerStatus,
	TunerStreamInfo,
	TunerDebug,
	TunerLockKey,
	IrTarget,
	LineupLocation,
	SysModel,
	SysFeatures,
	SysVersion,
	SysCopyright,
	SysDebug,
	SysHwModel,
}

func (p Property) GetValue() string {
	switch p {
	case TunerChannel:
		return "/tuner0/channel\x00"
	case TunerChannelMap:
		return "/tuner0/channelmap\x00"
	case TunerFilter:
		return "/tuner0/filter\x00"
	case TunerProgram:
		return "/tuner0/program\x00"
	case TunerTarget:
		return "/tuner0/targer\x00"
	case TunerStatus:
		return "/tuner0/status\x00"
	case TunerStreamInfo:
		return "/tuner0/streaminfo\x00"
	case TunerDebug:
		return "/tuner0/debug\x00"
	case TunerLockKey:
		return "/tuner0/lockkey\x00"
	case IrTarget:
		return "/ir/target\x00"
	case LineupLocation:
		return "/lineup/location\x00"
	case SysModel:
		return "/sys/model\x00"
	case SysFeatures:
		return "/sys/features\x00"
	case SysVersion:
		return "/sys/version\x00"
	case SysCopyright:
		return "/sys/copyright\x00"
	case SysDebug:
		return "/sys/debug\x00"
	case SysHwModel:
		return "/sys/hwmodel\x00"
	}
	return "Unknown property"
}
