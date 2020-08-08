package profinet

import (
	"encoding/binary"
	"fmt"
	"log"
)

//Cosntant
const (
	Unknown = iota
	ConnectionReq
	NegotiationReq
	ReadCounterReq
	ReadTimerReq
	ReadInputReq
	ReadOutputReq
	ReadMBReq
	ReadDBReq
	ReadMulti
	WriteCounterReq
	WriteTimerReq
	WriteInputReq
	WriteOutputReq
	WriteMBReq
	WriteDBReq
	WriteMulti
	PLCStartReq
	PLCStopReq
	GetCPUInfo
	GetCPInfo
)

//Telegram structure
type Telegram struct {
	//TPKT HEADER RFC 100 (4 bytes)
	Rvsm         byte   //always 3
	Reserved1    byte   //Reserve 0
	PacketLength uint16 //0x00
	// COPT for connection (3 bytes)
	CoptLenght  byte
	CoptPDUType byte
	CoptTPDU    byte
	//S7PDU HEADER (10-12 bytes)
	ProtocolID      byte   //(Always 0x32:50)
	MessageType     byte   //Message type
	Reserve2        byte   //0x00
	Reserve3        byte   //0x00
	PDUReference    uint16 //PDU Refrence  PLC just copy it to the reply
	ParameterLength uint16 //Parameter Lngth : quantity of bytes after the data lenght
	DataLenght      uint16 //Data Length

}

//NewTelegram Create a new telegram
func NewTelegram() (t *Telegram) {
	t = new(Telegram)
	t.Rvsm = 3
	t.ProtocolID = 50
	return
}

func (t Telegram) getTelegram(rq []byte) ([]byte, error) {
	out := make([]byte, 0)
	command := getCommand(rq)
	//TPKT HEADER RFC 100
	out = append(out, 0x03, 0x00, 0x00, rq[3]) //b0-1		//pvrsn (always 3), Reserved,packet lenght [High,Low

	if command == ConnectionReq { //ISO CONNECION REQUEST
		log.Println("ISO REQUEST:      " + fmt.Sprint(rq))
		//ISO-COTP
		out = append(out, 0x00, 0xD0, 0x00) //b4-6// COPT for connection
		for i := 0; i < int(rq[3]-isoHSize); i++ {
			out = append(out, 0)
		}
		return out, nil
	}
	//ISO-COTP
	out = append(out, 0x02, 0xF0, 0x80) //b4-6	// seen to be always 2,240,128

	//S7PDU HEADER (10-12 bytes)
	out = append(out, 0x32)           //b7		//Protocol id (Always 0x32:50)
	out = append(out, rq[8])          //b8		//Message type 0x01(job reques), 0x02 (Ack), 0x03 (Ack Data), 0x07 (UserData)
	out = append(out, 0x00, 0x00)     //b9-10	//Reserved
	out = append(out, rq[11], rq[12]) //b11-12	//PDU Refrence  PLC just copy it to the reply Little endian
	out = append(out, 0x00, 0x00)     //b13-14	//Parameter Length : quantity of bytes after the data lenght
	out = append(out, 0x00, 0x00)     //b15-16	//Data Length

	if command == NegotiationReq { //NEGOTIATIN REQUEST
		log.Println("NEGOTIATION:      " + fmt.Sprint(rq))
		out[3] = 27                                                       //packetlenght
		out = append(out, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) //b17-24
		out = append(out, rq[23], rq[24])                                 //b25-26 //PDULenght
		return out, nil
	}

	if command == ReadCounterReq || command == ReadTimerReq || command == ReadInputReq || command == ReadOutputReq || command == ReadMBReq || command == ReadDBReq || command == ReadMulti { //READ REQUET
		data := binary.BigEndian.Uint16([]byte{rq[23], rq[24]})
		out = append(out, 0x00)           //b17
		out = append(out, 0x00)           //b18
		out = append(out, 0x00, 0x00)     //b19-20
		out = append(out, 0xFF)           //b21			//ErrorCode
		out = append(out, rq[22])         //b22 		//Variabltype see table below
		out = append(out, rq[23], rq[24]) //b23-24		//Count
		for i := 0; i < int(data); i++ {
			out = append(out, byte(i)) //data
		}
		out[20] = 1
		out[3] = 25 + byte(data) //update packet lenght TODO both byte 2nd 3
		return out, nil
	}

	if command == WriteCounterReq || command == WriteTimerReq || command == WriteInputReq || command == WriteOutputReq || command == WriteMBReq || command == WriteDBReq || command == WriteMulti { //WRITE REQUEST
		out = append(out, 0x12)   //b17			//Specification type fr const 18 for read/write
		out = append(out, 0xFF)   //b18			// Lenght rest of byte
		out = append(out, 22)     //b19			//Syntax ID const 16 fr any typ addr
		out = append(out, rq[22]) //b20			//Variable type see table below
		out = append(out, 0xFF)   //b21			//
		out = append(out, 0x00)   //b22			//Count
		out[3] = 22
		return out, nil
	}
	if command == GetCPUInfo { //CPU INFO REQUEST
		log.Println("CPU REQUEST:    " + fmt.Sprint(rq))
		out = append(out, 0x12)   //b17			//Specification type fr const 18 for read/write
		out = append(out, 0xFF)   //b18			// Lenght rest of byte
		out = append(out, 22)     //b19			//Syntax ID const 16 fr any typ addr
		out = append(out, rq[22]) //b20			//Variable type see table below
		out = append(out, 0xFF)   //b21			//
		out = append(out, 0x00)   //b22			//Count
		out = append(out, 0x12)   //b23			//Specification type fr const 18 for read/write
		out = append(out, 0xFF)   //b24			// Lenght rest of byte
		out = append(out, 22)     //b25			//Syntax ID const 16 fr any typ addr
		out = append(out, rq[22]) //b26			//Variable type see table below
		out = append(out, 0xFF)   //b27		//
		out = append(out, 0x00)   //b28			//Count
		out = append(out, 0xFF)   //b29			//Specification type fr const 18 for read/write
		out = append(out, 0xFF)   //b30			// Lenght rest of byte
		out = append(out, 22)     //b31			//Syntax ID const 16 fr any typ addr
		out = append(out, 22)     //b32			//Syntax ID const 16 fr any typ addr
		for i := 0; i < 20; i++ {
			out = append(out, byte(i)) //data
		}
		out[3] = byte(len(out))
		return out, nil
	}
	return out, nil
}
func getCommand(rq []byte) (comm int) {

	if rq[4] == 17 && rq[5] == 224 {
		comm = ConnectionReq
	} else if rq[8] == 0x01 {
		if rq[17] == 0xF0 { //NEGOTIATION REQUEST
			comm = NegotiationReq
		} else if rq[17] == 0x04 { //READ REQUEST
			if rq[27] == 0x07 {
				comm = ReadMulti
				log.Println("READ MULTI:       " + fmt.Sprint(rq))
			} else if rq[27] == 0x1C {
				comm = ReadCounterReq
				log.Println("READ COUNTER:     " + fmt.Sprint(rq))
			} else if rq[27] == 0x1D {
				comm = ReadTimerReq
				log.Println("READ TIMER:       " + fmt.Sprint(rq))
			} else if rq[27] == 0x81 {
				comm = ReadInputReq
				log.Println("READ INPUT:       " + fmt.Sprint(rq))
			} else if rq[27] == 0x82 {
				comm = ReadOutputReq
				log.Println("READ OUTPUT:      " + fmt.Sprint(rq))
			} else if rq[27] == 0x83 {
				comm = ReadMBReq
				log.Println("READ MB:          " + fmt.Sprint(rq))
			} else if rq[27] == 0x84 {
				comm = ReadDBReq
				log.Println("READ DB:          " + fmt.Sprint(rq))
			} else {
				comm = Unknown
				log.Println("UNKNOWN REQUEST:  " + fmt.Sprint(rq))
				return
			}
		} else if rq[17] == 0x05 { //WRITE REQUEST
			if rq[27] == 0x07 {
				comm = WriteMulti
				log.Println("WRITE MULTI:      " + fmt.Sprint(rq))
			} else if rq[27] == 0x1C {
				comm = WriteCounterReq
				log.Println("WRITE COUNTER:    " + fmt.Sprint(rq))
			} else if rq[27] == 0x1D {
				comm = WriteTimerReq
				log.Println("WRITE TIMER:      " + fmt.Sprint(rq))
			} else if rq[27] == 0x81 {
				comm = WriteInputReq
				log.Println("WRITE INPUT:      " + fmt.Sprint(rq))
			} else if rq[27] == 0x82 {
				comm = WriteOutputReq
				log.Println("WRITE OUTPUT:     " + fmt.Sprint(rq))
			} else if rq[27] == 0x83 {
				comm = WriteMBReq
				log.Println("WRITE MB:         " + fmt.Sprint(rq))
			} else if rq[27] == 0x84 {
				comm = WriteDBReq
				log.Println("WRITE DB:         " + fmt.Sprint(rq))
			} else {
				comm = Unknown
				log.Println("UNKNOWN REQUEST:  " + fmt.Sprint(rq))
				return
			}
		} else if rq[17] == 0x29 { //PLC STOP
			comm = PLCStopReq
			log.Println("PLC STOP REQUEST: " + fmt.Sprint(rq))
		} else if rq[17] == 0x28 { //PLC START
			comm = PLCStartReq
			log.Println("PLC START REQUEST:" + fmt.Sprint(rq))
		} else {
			comm = Unknown
			log.Println("UNKNOWN REQUEST:  " + fmt.Sprint(rq))
			return
		}
	} else if rq[8] == 0x07 {
		if rq[30] == 28 { //GET CPU INFO
			comm = GetCPUInfo
			log.Println("CPU INFO REQUEST: " + fmt.Sprint(rq))
		} else if rq[30] == 49 { //GET CPU INFO
			comm = GetCPInfo
			log.Println("CP INFO REQUEST:  " + fmt.Sprint(rq))
		} else {
			comm = Unknown
			log.Println("UNKNOWN REQUEST:  " + fmt.Sprint(rq))
			return
		}
	} else {
		comm = Unknown
		log.Println("UNKNOWN REQUEST:  " + fmt.Sprint(rq))
		return
	}
	return
}

//GetTPKTHeader get the tpkt header
func (t Telegram) GetTPKTHeader() (b []byte) {
	b = []byte{t.Rvsm, t.Reserved1}
	b = append(b, uint16ToByte(t.PacketLength)...)
	return
}

func (t *Telegram) setPAcketLength(input []byte) (b []byte) {
	b = []byte{t.Rvsm, t.Reserved1}
	b = append(b, uint16ToByte(t.PacketLength)...)
	return
}

func uint16ToByte(input uint16) (b []byte) {
	binary.BigEndian.PutUint16(b, input)
	return
}

/*
Vaable type
===========
0x01 - BIT
0x02 - BYT
0x03 - CHA
0x04 - WOR
0x05 - INT
0x06 - DWO
0x07 - DIN
0x08 - REA
0x09 - DAT
0x0A - TOD
0x0B - TIM
0x0C - S5T
0x0F - DATENDIME
0x1C - COUNTR
0x1D - TIMER
x1E - IEC TIR
x1F - IEC COUNR
20 - HS COUNTER




ea table
=========
0x03 - Sysm info of S200 family
0x05 - System flags ofS200 famil
0x06 - Analog inputsof S200 famil
0x07 - Analog outputs f S200 fail
0x1C - S7 counters ()
0x1D - S7 timers (T)
0x1E - IEC countes (0 family
0x1F - IEC timers 200 family)
0x80 - Direct perpheral cces P)
0x81 - Inputs (I)
0x82 - Outputs (Q
0x83 - Flags (M) erer
0x84 - Data blocks (DB
0x5 - Instance data bcs (DI)
0x86 - Local data (L)
0x7 - Unknown yet (V
*/
