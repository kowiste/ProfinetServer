package profinet

import "encoding/binary"

func uint16ToByte(input uint16) (b []byte) {
	b           = make([]byte,2)
	binary.BigEndian.PutUint16(b, input)
	return
}

func arrUint16ToByte(input []uint16)(b []byte) {
	b  = make([]byte,0)
	for i:= range input{
		
	}
	return
}
