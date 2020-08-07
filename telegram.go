package profinet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
)

type Telegram struct {
	//TPKT HEADER RFC 100
	Rvsm         byte  //always 3
	Reserved     byte  //Reserve 0
	PacketLength int16 //
	// COPT for connection
	COPT0 byte
	COPT1 byte
	COPT3 byte
}

func getTelegram(rq []byte) ([]byte, error) {
	negReq := false
	readReq := false
	writeReq := false
	connReq := false
	out := make([]byte, 0)
	if rq[4] == 17 && rq[5] == 224 {
		connReq = true
	} else if rq[17] == 0xF0 {
		negReq = true
	} else if rq[17] == 0x04 {
		readReq = true
	} else if rq[17] == 0x05 {
		writeReq = true
	} else {
		log.Println("UNKNOWN REQUEST:  " + fmt.Sprint(rq))
		return nil, errors.New("REQUEST NOT IMPLEMENT")
	}

	//TPKT HEADER RFC 100
	out = append(out, 0x03, 0x00, 0x00, rq[3]) //b0-1		//pvrsn (always 3), Reserved,packet lenght [High,Low
	//ISO CONNECION REQUEST
	if connReq {
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
	out = append(out, 0x01)           //b8		//Message tye 0x01(job reques), 0x02 (Ack), 0x03 (Ack Data), 0x07 (UserData)
	out = append(out, 0x00, 0x00)     //b9-10	//Reserved
	out = append(out, rq[11], rq[12]) //b11-12	//PDU Refrence  PLC just copy it to the reply
	out = append(out, 0x00, 0x00)     //b13-14	//Parameter Lngth : quantity of bytes after te data lenght
	out = append(out, 0x00, 0x00)     //b15-16	//Data Length

	//NEGOTIATIN REQUEST
	if negReq {
		log.Println("NEGOTIATION:      " + fmt.Sprint(rq))
		out[3] = 27                                                       //packetlenght
		out = append(out, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) //b17-24
		out = append(out, rq[23], rq[24])                                 //b25-26
		return out, nil
	}

	//READ REQUET
	if readReq {

		data := binary.BigEndian.Uint16([]byte{rq[23], rq[24]})
		log.Println("READ REQUEST:     " + fmt.Sprint(rq))
		out = append(out, 0x00)           //b17
		out = append(out, 0x00)           //b18
		out = append(out, 0x00, 0x00)     //b19-20
		out = append(out, 0xFF)           //b21			//ErrorCode
		out = append(out, rq[22])         //b22 		//Variabltype see table below
		out = append(out, rq[23], rq[24]) //b23-24		//Count
		for i := 0; i < int(data); i++ {
			out = append(out, byte(i)) //data
		}
		out[3] = 25 + byte(data) //update packet lenght TODO both byte 2nd 3
		return out, nil
	}
	//READ REQUES
	if writeReq {
		log.Println("WRITE REQUEST:    " + fmt.Sprint(rq))
		out = append(out, 0x12)   //b17			//Specification type fr const 18 for read/write
		out = append(out, 0xFF)   //b18			// Lenght rest of byte
		out = append(out, 22)     //b19			//Syntax ID const 16 fr any typ addr
		out = append(out, rq[22]) //b20			//ariable type see table below
		out = append(out, 0xFF)   //b21			//
		out = append(out, 0x00)   //b22			//Count
		out[3] = 22
		return out, nil
	}
	return out, nil

}

/*
Vaiable tye
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
0x0C - S5TI
0x0F - DATE NDIME
0x1C - COUNTR
0x1D - TIMER
x1E - IEC TIER
x1F - IEC COUNTR
20 - HS COUNTER




rea table
==========
0x03 - Sysem info of S200 family
0x05 - System flags ofS200 family
0x06 - Analog inputsof S200 famil
0x07 - Analog outputs f S200 fail
0x1C - S7 counters ()
0x1D - S7 timers (T)
0x1E - IEC countes (00 family
0x1F - IEC timers 200 family)
0x80 - Direct perpheral cces (P)
0x81 - Inputs (I)
0x82 - Outputs (Q
0x83 - Flags (M) (erer
0x84 - Data blocks (DB
0x5 - Instance data bcks (DI)
0x86 - Local data (L)
0x7 - Unknown yet (V
*/
