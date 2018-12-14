package hdhomerun

// Discovery protocol uses UDP and control uses TCP
const PORT = 65001

const HDHOMERUN_DEVICE_TYPE_WILDCARD = 0xFFFFFFFF
const HDHOMERUN_DEVICE_TYPE_TUNER = 0x00000001
const HDHOMERUN_DEVICE_TYPE_STORAGE = 0x00000005
const HDHOMERUN_DEVICE_ID_WILDCARD = 0xFFFFFFFF

type HDHomerunEmulator struct {
	DeviceID uint32
}
