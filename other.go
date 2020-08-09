package profinet

import "encoding/binary"

func uint16ToByte(input uint16) (b []byte) {
	binary.BigEndian.PutUint16(b, input)
	return
}
