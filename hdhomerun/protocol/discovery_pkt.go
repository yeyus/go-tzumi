package protocol

import "encoding/binary"

type DiscoveryResponseMsg struct {
	DeviceType uint32
	DeviceId   uint32
}

func (msg *DiscoveryResponseMsg) Packet() Packet {
	deviceType := make([]byte, 4)
	binary.BigEndian.PutUint32(deviceType, msg.DeviceType)

	deviceId := make([]byte, 4)
	binary.BigEndian.PutUint32(deviceId, msg.DeviceId)

	p := Packet{
		Type: HDHOMERUN_TYPE_DISCOVER_RPY,
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
		},
	}

	return p
}
