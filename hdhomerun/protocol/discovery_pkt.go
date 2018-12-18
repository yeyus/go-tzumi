package protocol

import (
	"encoding/binary"
	"errors"
)

const HDHOMERUN_DEVICE_TYPE_WILDCARD = 0xFFFFFFFF
const HDHOMERUN_DEVICE_TYPE_TUNER = 0x00000001
const HDHOMERUN_DEVICE_TYPE_STORAGE = 0x00000005
const HDHOMERUN_DEVICE_ID_WILDCARD = 0xFFFFFFFF

type DiscoveryMsg struct {
	Type       PacketType
	DeviceType uint32
	DeviceId   uint32
	TunerCount uint8
}

func ParseDiscoveryMsg(packet *Packet) (*DiscoveryMsg, error) {
	if !(packet.Type == HDHOMERUN_TYPE_DISCOVER_REQ || packet.Type == HDHOMERUN_TYPE_DISCOVER_RPY) {
		return nil, errors.New("invalid package type")
	}

	deviceType := packet.GetTLVByTag(HDHOMERUN_TAG_DEVICE_TYPE)
	if deviceType == nil {
		return nil, errors.New("missing HDHOMERUN_TAG_DEVICE_TYPE")
	}

	deviceId := packet.GetTLVByTag(HDHOMERUN_TAG_DEVICE_ID)
	if deviceId == nil {
		return nil, errors.New("missing HDHOMERUN_TAG_DEVICE_TYPE")
	}

	tunerCount := uint8(0)
	tunerCountTlv := packet.GetTLVByTag(HDHOMERUN_TAG_TUNER_COUNT)
	if tunerCountTlv != nil {
		tunerCount = uint8(tunerCountTlv.Value[0])
	}

	msg := &DiscoveryMsg{
		Type:       packet.Type,
		DeviceType: binary.BigEndian.Uint32(deviceType.Value),
		DeviceId:   binary.BigEndian.Uint32(deviceId.Value),
		TunerCount: tunerCount,
	}

	return msg, nil
}

func (msg *DiscoveryMsg) Packet() Packet {
	deviceType := make([]byte, 4)
	binary.BigEndian.PutUint32(deviceType, msg.DeviceType)

	deviceId := make([]byte, 4)
	binary.BigEndian.PutUint32(deviceId, msg.DeviceId)

	tunerCount := make([]byte, 1)
	tunerCount[0] = msg.TunerCount
	p := Packet{
		Type: msg.Type,
		Tags: []*TLV{
			&TLV{
				Tag:    HDHOMERUN_TAG_DEVICE_TYPE,
				Length: 4,
				Value:  deviceType,
			},
			&TLV{
				Tag:    HDHOMERUN_TAG_DEVICE_ID,
				Length: 4,
				Value:  deviceId,
			},
			&TLV{
				Tag:    HDHOMERUN_TAG_TUNER_COUNT,
				Length: 1,
				Value:  tunerCount,
			},
		},
	}

	return p
}

func (msg *DiscoveryMsg) IsAddressedToSelf(ownDeviceId uint32) bool {
	return (msg.DeviceType == HDHOMERUN_DEVICE_TYPE_TUNER || msg.DeviceType == HDHOMERUN_DEVICE_TYPE_WILDCARD) && (msg.DeviceId == HDHOMERUN_DEVICE_ID_WILDCARD || msg.DeviceId == ownDeviceId)
}
