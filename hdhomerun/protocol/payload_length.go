package protocol

/*
 * The length field can be one or two bytes long.
 * For a length <= 127 bytes the length is expressed as a single byte. The
 * most-significant-bit is clear indicating a single-byte length.
 * For a length >= 128 bytes the length is expressed as a sequence of two bytes as follows:
 * The first byte is contains the least-significant 7-bits of the length. The
 * most-significant bit is then set (add 0x80) to indicate that it is a two byte length.
 * The second byte contains the length shifted down 7 bits.
 */

func DecodeLength(d []byte) uint8 {
	length := d[0]
	if d[0]&0x80 > 0 {
		length &= 0x7F
		length |= d[1] << 7
	} else {
		length = d[0]
	}
	return length
}

func EncodeLength(length uint8) []byte {
	d := make([]byte, 2)
	if length <= 127 {
		d[0] = byte(length)
		return d[:1]
	} else {
		d[0] = (length | 0x80)
		d[1] = length >> 7
		return d[:2]
	}
}
