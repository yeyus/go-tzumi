package protocol

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
)

// Documentation
// https://github.com/Silicondust/libhdhomerun/blob/master/hdhomerun_pkt.h

type PacketType uint16

const HDHOMERUN_TYPE_DISCOVER_REQ PacketType = 0x0002
const HDHOMERUN_TYPE_DISCOVER_RPY PacketType = 0x0003
const HDHOMERUN_TYPE_GETSET_REQ PacketType = 0x0004
const HDHOMERUN_TYPE_GETSET_RPY PacketType = 0x0005
const HDHOMERUN_TYPE_UPGRADE_REQ PacketType = 0x0006
const HDHOMERUN_TYPE_UPGRADE_RPY PacketType = 0x0007

type Packet struct {
	Type          PacketType
	Length        uint16
	Payload       []byte
	Tags          []*TLV
	CRC           uint32
	calculatedCRC uint32
}

type PayloadTag uint8

const HDHOMERUN_TAG_DEVICE_TYPE PayloadTag = 0x01
const HDHOMERUN_TAG_DEVICE_ID PayloadTag = 0x02
const HDHOMERUN_TAG_GETSET_NAME PayloadTag = 0x03
const HDHOMERUN_TAG_GETSET_VALUE PayloadTag = 0x04
const HDHOMERUN_TAG_GETSET_LOCKKEY PayloadTag = 0x15
const HDHOMERUN_TAG_ERROR_MESSAGE PayloadTag = 0x05
const HDHOMERUN_TAG_TUNER_COUNT PayloadTag = 0x10
const HDHOMERUN_TAG_DEVICE_AUTH_BIN PayloadTag = 0x29
const HDHOMERUN_TAG_BASE_URL PayloadTag = 0x2A
const HDHOMERUN_TAG_DEVICE_AUTH_STR PayloadTag = 0x2B

type TLV struct {
	Tag    PayloadTag
	Length uint16
	Value  []byte
}

// Packet
func NewPacket(data []byte) *Packet {
	l := len(data)
	length := uint16(data[2]<<8 | data[3])
	p := &Packet{
		// are big endian
		Type:    PacketType(data[0]<<8 | data[1]),
		Length:  length,
		Payload: data[4 : 4+length],
		// is IEEE CRC32 (little endian)
		CRC:           binary.LittleEndian.Uint32(data[4+length : l]),
		calculatedCRC: crc32.ChecksumIEEE(data[0 : l-4]),
		// Derived payload
		Tags: ParseAllTLV(data[4 : 4+length]),
	}

	return p
}

func (p *Packet) IsValid() bool {
	return p.CRC == p.calculatedCRC
}

func (p *Packet) IsRequest() bool {
	return p.Type == HDHOMERUN_TYPE_UPGRADE_REQ || p.Type == HDHOMERUN_TYPE_GETSET_REQ || p.Type == HDHOMERUN_TYPE_DISCOVER_REQ
}

func (p *Packet) IsReply() bool {
	return p.Type == HDHOMERUN_TYPE_UPGRADE_RPY || p.Type == HDHOMERUN_TYPE_GETSET_RPY || p.Type == HDHOMERUN_TYPE_DISCOVER_RPY
}

func (p *Packet) GetTLVByTag(tag PayloadTag) (found *TLV) {
	for _, tlv := range p.Tags {
		if tlv.Tag == tag {
			found = tlv
			break
		}
	}

	return
}

func (p *Packet) Serialize() []byte {
	buf := make([]byte, 3072)
	length := uint16(0)

	// type
	binary.BigEndian.PutUint16(buf[:], uint16(p.Type))

	// tlvs
	for _, tlv := range p.Tags {
		tlvSlice := tlv.Serialize()
		copy(buf[length+4:], tlvSlice)
		length += uint16(len(tlvSlice))
	}

	// write final length
	binary.BigEndian.PutUint16(buf[2:], length)

	// CRC
	binary.LittleEndian.PutUint32(buf[length+2+2:], crc32.ChecksumIEEE(buf[0:length+2+2]))

	return buf[:2+2+length+4]
}

func (p *Packet) String() string {
	crc := "Incorrect"
	if p.IsValid() {
		crc = "OK"
	}

	return fmt.Sprintf("Packet{ Type: %s, Length: %d, Tags: %s, CRC: %s}", p.Type, p.Length, p.Tags, crc)
}

// Payloads
func NewTLV(payload []byte) *TLV {
	length := DecodeLength(payload[1:2])

	start := 2
	if length >= 128 {
		start = 3
	}

	end := start + int(length)
	return &TLV{
		Tag:    PayloadTag(payload[0]),
		Length: uint16(length),
		Value:  payload[start:end],
	}
}

func (t *TLV) Serialize() []byte {
	buf := make([]byte, 3072)
	l := 0

	buf[0] = uint8(t.Tag)
	l++

	lengthBytes := EncodeLength(uint8(len(t.Value)))
	copy(buf[1:], lengthBytes)
	l += len(lengthBytes)

	copy(buf[l:], t.Value)
	l += len(t.Value)

	return buf[:l]
}

func ParseAllTLV(payload []byte) []*TLV {
	if len(payload) == 0 {
		return []*TLV{}
	}

	tlv := NewTLV(payload)
	leap := 1 + 1 + tlv.Length
	// check if we use double byte length notation
	if tlv.Length >= 128 {
		leap++
	}

	return append([]*TLV{tlv}, ParseAllTLV(payload[leap:])...)
}
