package protocol

import "errors"

type GetSetMsg struct {
	Type       PacketType
	TunerIndex uint8
	Property   Property
	Value      []byte
}

func ParseGetSetMsg(packet *Packet) (*GetSetMsg, error) {
	if !(packet.Type == HDHOMERUN_TYPE_GETSET_REQ || packet.Type == HDHOMERUN_TYPE_GETSET_RPY) {
		return nil, errors.New("invalid package type")
	}

	propertyTag := packet.GetTLVByTag(HDHOMERUN_TAG_GETSET_NAME)
	if propertyTag == nil {
		return nil, errors.New("invalid getset name")
	}
	property := parseProperty(string(propertyTag.Value))

	valueTag := packet.GetTLVByTag(HDHOMERUN_TAG_GETSET_VALUE)
	var value []byte
	if valueTag != nil {
		value = valueTag.Value
	}

	msg := &GetSetMsg{
		Type:     packet.Type,
		Property: property,
		Value:    value,
	}

	return msg, nil
}

func (msg *GetSetMsg) Packet() *Packet {
	// TODO support tuner index encoding
	var tlvs []*TLV
	if msg.Type == HDHOMERUN_TYPE_GETSET_REQ {
		tlvs = []*TLV{
			&TLV{
				Tag:    HDHOMERUN_TAG_GETSET_NAME,
				Length: uint16(len(msg.Property.GetValue())),
				Value:  []byte(msg.Property.GetValue()),
			},
		}
	} else {
		tlvs = []*TLV{
			&TLV{
				Tag:    HDHOMERUN_TAG_GETSET_NAME,
				Length: uint16(len(msg.Property.GetValue())),
				Value:  []byte(msg.Property.GetValue()),
			},
			&TLV{
				Tag:    HDHOMERUN_TAG_GETSET_VALUE,
				Length: uint16(len(msg.Value)),
				Value:  msg.Value,
			},
		}
	}
	p := Packet{
		Type: msg.Type,
		Tags: tlvs,
	}

	return &p
}

func parseProperty(p string) Property {
	// TODO parse regexps for tuner index
	for _, prop := range Properties {
		if prop.GetValue() == p {
			return prop
		}
	}

	return UnknownProperty
}
