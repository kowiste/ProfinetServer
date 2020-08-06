package profinet

import (
	"fmt"
	"log"
)

func getTelegram(rq []byte) ([]byte, error) {
	//[3 0 0 31 2 240 128 50 1 0 0 5 0 0 14 0 0 4 1 18 10 16 2 0 2 1 1 132 0 0 80]
	out := make([]byte, 0)
	//TPKT HEADER RFC 1006
	out = append(out, 0x03, 0x00)  //b0-1		//pvrsn (always 3), Reserved
	out = append(out, 0x00, rq[3]) //b2-3		//packet lenght [High,Low]
	//ISO-COTP
	out = append(out, 0x02, 0xF0, 0x80) //b4-6	// seen to be always 2,240,128
	if rq[4] == 17 && rq[5] == 224 {    //ISO CONNECTION REQUEST
		log.Println("ISO REQUEST:  " + fmt.Sprint(rq))
		//log.Println(rq)
		for i := 0; i < int(rq[3]-3); i++ {
			out = append(out, 0)
		}
		out[3] = rq[3]
		out[5] = 0xD0
		return out, nil
	}
	//S7PDU
	//HEADER (10-12 bytes)
	out = append(out, 0x32)           //b7		//Protocol id (Always 0x32:50)
	out = append(out, 0x01)           //b8		//Message type 0x01(job request), 0x02 (Ack), 0x03 (Ack Data), 0x07 (UserData)
	out = append(out, 0x00, 0x00)     //b9-10	//Reserved
	out = append(out, rq[11], rq[12]) //b11-12	//PDU Reference  PLC just copy it to the reply
	out = append(out, 0x00, 0x00)     //b13-14	//Parameter Length : quantity of bytes after the data lenght
	out = append(out, 0x00, 0x00)     //b15-16	//Data Length

	switch rq[17] {

	case 0xF0: //Negogiation PDU
		log.Println("NEGOTIATION:  " + fmt.Sprint(rq))
		out[3] = 27
		out = append(out, 0xF0)           //b17
		out = append(out, 0x00)           //b18
		out = append(out, 0x00, 0x00)     //b19-20
		out = append(out, 0xFF)           //b21
		out = append(out, 0x00)           //b22
		out = append(out, 0x00, 0x00)     //b23-24
		out = append(out, rq[23], rq[24]) //b25-26

	default:

		if out[6] == 0x02 {
			out = append(out, 0x00) //b15			//Error Class 	[optional only present in the Ack-Data message]
			out = append(out, 0x00) //b16			//Error Code	[optional only present in the Ack-Data message]
		}
		out = append(out, 0x00)             //b15
		out = append(out, 0x00)             //b16
		out = append(out, 0x12)             //b17				//Specification type for const 18 for read/write
		out = append(out, 0x00)             //b18				// Lenght rest of byte
		out = append(out, 0x10)             //b19				//Syntax ID const 16 for any type addr
		out = append(out, 0x00)             //b20				//Variable type see tble below
		out = append(out, 0x00)             //b21-22			//Count
		out = append(out, 0x00)             //b25-26 			//DB number
		out = append(out, 0x00)             //b27 				//Area ee the table below
		out = append(out, 0x00, 0x00, 0x00) //b28-29-30 		//Addres
	}

	return out, nil
}

/*
Variable tye
============
0x01 - BIT
0x02 - BYT
0x03 - CHAR
0x04 - WORD
0x05 - INT
0x06 - DWOR
0x07 - DIN
0x08 - REAL
0x09 - DATE
0x0A - TOD
0x0B - TIME
0x0C - S5TIM
0x0F - DATE AND IME
0x1C - COUNTER
0x1D - TIMER
x1E - IEC TIMER
x1F - IEC COUNTER
x20 - HS COUNTER




Area table
==========
0x03 - System info of S200 family
0x05 - System flags ofS200 family
0x06 - Analog inputsof S200 family
0x07 - Analog outputs of S200 faily
0x1C - S7 counters (C)
0x1D - S7 timers (T)
0x1E - IEC countes (200 family)
0x1F - IEC timers 200 family)
0x80 - Direct peripheral ccess (P)
0x81 - Inputs (I)
0x82 - Outputs (Q)
0x83 - Flags (M) (Merer)
0x84 - Data blocks (DB
0x5 - Instance data blocks (DI)
0x86 - Local data (L)
0x87 - Unknown yet (V)
*/
