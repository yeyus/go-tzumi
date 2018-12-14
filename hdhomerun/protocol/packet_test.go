package protocol

import "testing"
import "github.com/stretchr/testify/assert"

func TestNewPacket_InvalidCRC(t *testing.T) {
	fakePkt := []byte{
		// Packet type
		0x00, 0x04,
		// Packet length
		0x00, 19,
		// Payload
		// -- TLV0: type
		0x03,
		// -- TLV0: length
		10,
		// -- TLV0: value
		'H', 'o', 'l', 'a', 'M', 'u', 'n', 'd', 'o', 0,
		// -- TLV1: type
		0x04,
		// -- TLV1: length
		0x05,
		// -- TLV1: value
		'1', '2', '3', '4', 0,
		// CRC32
		0x00, 0x00, 0x00, 0x00,
	}

	pkt := NewPacket(fakePkt[:])

	assert.Equal(t, pkt.Type, HDHOMERUN_TYPE_GETSET_REQ, "expected packet to be HDHOMERUN_TYPE_GETSET_REQ")
	assert.Equal(t, pkt.Length, uint16(19), "expected length to be 19")
	assert.Equal(t, pkt.Payload, fakePkt[4:23], "expected payload to match")
	assert.Equal(t, len(pkt.Tags), 2, "expected 2 tags")
	assert.Equal(t, pkt.IsRequest(), true, "expected to be a request")
	assert.Equal(t, pkt.IsReply(), false, "expected to not be a reply")
	assert.Equal(t, pkt.IsValid(), false, "expected packet to be not valid")
}

func TestNewPacket_ValidCRC(t *testing.T) {
	fakePkt := []byte{
		// Packet type
		0x00, 0x04,
		// Packet length
		0x00, 19,
		// Payload
		// -- TLV0: type
		0x03,
		// -- TLV0: length
		10,
		// -- TLV0: value
		'H', 'o', 'l', 'a', 'M', 'u', 'n', 'd', 'o', 0,
		// -- TLV1: type
		0x04,
		// -- TLV1: length
		0x05,
		// -- TLV1: value
		'1', '2', '3', '4', 0,
		// CRC32
		0xF0, 0xBC, 0x83, 0x2A,
	}

	pkt := NewPacket(fakePkt[:])

	assert.Equal(t, pkt.Type, HDHOMERUN_TYPE_GETSET_REQ, "expected packet to be HDHOMERUN_TYPE_GETSET_REQ")
	assert.Equal(t, pkt.Length, uint16(19), "expected length to be 19")
	assert.Equal(t, pkt.Payload, fakePkt[4:23], "expected payload to match")
	assert.Equal(t, len(pkt.Tags), 2, "expected 2 tags")
	assert.Equal(t, pkt.IsRequest(), true, "expected to be a request")
	assert.Equal(t, pkt.IsReply(), false, "expected to not be a reply")
	assert.Equal(t, pkt.IsValid(), true, "expected packet to be not valid")
}

func TestNewPacket_DiscoveryRequest(t *testing.T) {
	fakePkt := []byte{
		0, 2, 0, 12, 1, 4, 0, 0, 0, 1, 2, 4, 255, 255, 255, 255, 78, 80, 127, 53,
	}

	pkt := NewPacket(fakePkt[:])

	assert.Equal(t, pkt.Type, HDHOMERUN_TYPE_DISCOVER_REQ, "expected packet to be HDHOMERUN_TYPE_DISCOVER_REQ")
	assert.Equal(t, uint16(12), pkt.Length, "expected length to be 12")
	assert.Equal(t, true, pkt.IsRequest(), "expected to be a request")
	assert.Equal(t, false, pkt.IsReply(), "expected to not be a reply")
	assert.Equal(t, true, pkt.IsValid(), "expected packet to be valid")
	assert.Equal(t, 2, len(pkt.Tags), "expected 2 tags")

	// first tag
	assert.Equal(t, HDHOMERUN_TAG_DEVICE_TYPE, pkt.Tags[0].Tag, "first tag should be HDHOMERUN_TAG_DEVICE_TYPE")
	assert.Equal(t, uint16(4), pkt.Tags[0].Length, "first tag should be 4 bytes long")
	assert.Equal(t, []byte{0x00, 0x00, 0x00, 0x01}, pkt.Tags[0].Value, "first tag should have value 0x00000001")

	// second tag
	assert.Equal(t, HDHOMERUN_TAG_DEVICE_ID, pkt.Tags[1].Tag, "second tag should be HDHOMERUN_TAG_DEVICE_ID")
	assert.Equal(t, uint16(4), pkt.Tags[1].Length, "second tag should be 4 bytes long")
	assert.Equal(t, []byte{0xFF, 0xFF, 0xFF, 0xFF}, pkt.Tags[1].Value, "second tag should have value 0xFFFFFFFF")
}

func TestPacketSerialize(t *testing.T) {
	fakePkt := []byte{
		0, 2, 0, 12, 1, 4, 0, 0, 0, 1, 2, 4, 255, 255, 255, 255, 78, 80, 127, 53,
	}

	pkt := NewPacket(fakePkt[:])
	ser := pkt.Serialize()

	assert.Equal(t, fakePkt, ser, "going full circle on decode/encode should give equal byte arrays")
}

func TestParseAllTLV_SingleTLV(t *testing.T) {
	fakePkt := []byte{
		// Packet type
		0x00, 0x04,
		// Packet length
		0x00, 12,
		// Payload
		// -- TLV0: type
		0x03,
		// -- TLV0: length
		10,
		// -- TLV0: value
		'H', 'o', 'l', 'a', 'M', 'u', 'n', 'd', 'o', 0,
		// CRC32
		0x00, 0x00, 0x00, 0x00,
	}

	tlvs := ParseAllTLV(fakePkt[4 : len(fakePkt)-4])
	assert.Equal(t, tlvs[0].Tag, HDHOMERUN_TAG_GETSET_NAME, "expected tag to be HDHOMERUN_TAG_GETSET_NAME")
	assert.Equal(t, tlvs[0].Length, uint16(10), "expected length to be 10")
	assert.Equal(t, tlvs[0].Value, fakePkt[6:16], "expected value to match fake packet")
}

func TestParseAllTLV_MultipleTLV(t *testing.T) {
	fakePkt := []byte{
		// Packet type
		0x00, 0x04,
		// Packet length
		0x00, 19,
		// Payload
		// -- TLV0: type
		0x03,
		// -- TLV0: length
		10,
		// -- TLV0: value
		'H', 'o', 'l', 'a', 'M', 'u', 'n', 'd', 'o', 0,
		// -- TLV1: type
		0x04,
		// -- TLV1: length
		0x05,
		// -- TLV1: value
		'1', '2', '3', '4', 0,
		// CRC32
		0x00, 0x00, 0x00, 0x00,
	}

	tlvs := ParseAllTLV(fakePkt[4 : len(fakePkt)-4])
	assert.Equal(t, tlvs[0].Tag, HDHOMERUN_TAG_GETSET_NAME, "expected tag to be HDHOMERUN_TAG_GETSET_NAME")
	assert.Equal(t, tlvs[0].Length, uint16(10), "expected length to be 10")
	assert.Equal(t, tlvs[0].Value, fakePkt[6:16], "expected value to match fake packet")

	assert.Equal(t, tlvs[1].Tag, HDHOMERUN_TAG_GETSET_VALUE, "expected tag to be HDHOMERUN_TAG_GETSET_VALUE")
	assert.Equal(t, tlvs[1].Length, uint16(5), "expected length to be 5")
	assert.Equal(t, tlvs[1].Value, fakePkt[18:23], "expected value to match fake packet")
}
