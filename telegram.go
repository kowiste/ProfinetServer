package profinet

import (
	"encoding/binary"
	"fmt"
	"log"
)

func (s *Server) getTelegram(rq []byte) ([]byte, error) {
	out := make([]byte, 0)
	s.getCommand(rq) //Read the command, data lenght and other parameter

	out = append(out, s.getTPKHeader(rq)...) //b0-3//TPKT HEADER

	if s.command == ConnectionReq { //ISO CONNECION REQUEST
		log.Println("ISO REQUEST:      " + fmt.Sprint(rq))
		//ISO-COTP
		out = append(out, s.getISOCOPT()...)                           //b4-6// COPT for connection
		out = append(out, make([]byte, int(s.reqDataLen-isoHSize))...) //Add the number of 0 necessary
		return out, nil
	}
	out = append(out, s.getISOCOPT()...)     //b4-6	//ISO-COTP
	out = append(out, s.getS7PDUHeader()...) //S7PDU HEADER (10-12 bytes)

	if s.command == SetupCommunication { //NEGOTIATIN REQUEST
		log.Println("NEGOTIATION:      " + fmt.Sprint(rq)) //print
		out[3] = 27                                        //packetlenght
		out = append(out, make([]byte, 8)...)              //b17-24
		out = append(out, rq[23], rq[24])                  //b25-26 //PDULenght
		return out, nil
	}

	if s.command == ReadVariable { //command == ReadCounterReq || command == ReadTimerReq || command == ReadInputReq || command == ReadOutputReq || command == ReadMBReq || command == ReadDBReq || command == ReadMulti { //READ REQUET
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

	if s.command == WriteVariable { //command == WriteCounterReq || command == WriteTimerReq || command == WriteInputReq || command == WriteOutputReq || command == WriteMBReq || command == WriteDBReq || command == WriteMulti { //WRITE REQUEST
		out = append(out, 0x12)   //b17			//Specification type fr const 18 for read/write
		out = append(out, 0xFF)   //b18			// Lenght rest of byte
		out = append(out, 22)     //b19			//Syntax ID const 16 fr any typ addr
		out = append(out, rq[22]) //b20			//Variable type see table below
		out = append(out, 0xFF)   //b21			//
		out = append(out, 0x00)   //b22			//Count
		out[3] = 22
		return out, nil
	}
	if s.command == GetCPUInfo { //CPU INFO REQUEST
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
func (s *Server) getTPKHeader(rq []byte) (out []byte) {
	//TPKT HEADER RFC 100
	//b0-3		//pvrsn (always 3), Reserved,packet lenght [High,Low
	s.reqDataLen = int(binary.BigEndian.Uint16([]byte{rq[2], rq[3]}))
	out = append(out, ProtocolID, 0x00, rq[2], rq[3])
	return
}
func (s Server) getISOCOPT() (out []byte) {
	btemp := byte(0xF0)
	if s.command == ConnectionReq {
		btemp = 0xD0
	}
	out = append(out, 0x00, btemp, 0x00)
	return
}
func (s Server) getS7PDUHeader() (out []byte) {
	//S7PDU HEADER (10-12 bytes)
	out = append(out, ProtocolID)                      //b7		//Protocol id (Always 0x32:50)
	out = append(out, s.msgType)                       //b8		//Message type 0x01(job reques), 0x02 (Ack), 0x03 (Ack Data), 0x07 (UserData)
	out = append(out, 0x00, 0x00)                      //b9-10	//Reserved
	out = append(out, uint16ToByte(s.PDUReference)...) //b11-12	//PDU Reference  PLC just copy it to the reply Little endian
	out = append(out, 0x00, 0x00)                      //b13-14	//Parameter Length : quantity of bytes after the data lenght
	out = append(out, 0x00, 0x00)                      //b15-16	//Data Length
	return
}
func (s *Server) getCommand(rq []byte) {
	s.msgType = rq[8]
	s.PDUReference = binary.BigEndian.Uint16([]byte{rq[11], rq[12]})
	if rq[4] == 17 && rq[5] == 224 {
		s.command = ConnectionReq
	} else if s.msgType == JobRequest {
		s.command = rq[17]
		s.selectedJob(rq)
	} else if s.msgType == UserdataRequest {
		if rq[30] == 28 { //GET CPU INFO
			s.command = GetCPUInfo
			log.Println("CPU INFO REQUEST: " + fmt.Sprint(rq))
		} else if rq[30] == 49 { //GET CPU INFO
			s.command = GetCPInfo
			log.Println("CP INFO REQUEST:  " + fmt.Sprint(rq))
		} else {
			s.command = Unknown
			log.Println("UNKNOWN REQUEST:  " + fmt.Sprint(rq))
			return
		}
	} else {
		s.command = Unknown
		log.Println("UNKNOWN REQUEST:  " + fmt.Sprint(rq))
		return
	}
	return
}
func (s *Server) selectedJob(rq []byte) {
	req := ""
	switch s.command {
	case ReadVariable, WriteVariable:
		if s.command == ReadVariable {
			req = "READ "
		} else {
			req = "WRITE "
		}
		s.rwcomm = rq[27]
		req += s.printRWArea()
	case PLCStop:
		req = "PLC STOP REQUEST: "
	case PLCControl:
		req = "PLC START REQUEST:"
	default:
		req = "UNKNOWN REQUEST:  "
	}
	log.Println(req + fmt.Sprint(rq))
}
func (s *Server) printRWArea() (area string) {
	switch s.rwcomm {
	case MulipleRW:
		area = "MULTI:       "
	case CounterS7:
		area = "COUNTER:     "
	case TimerS7:
		area = "TIMER:       "
	case Input:
		area = "INPUT:       "
	case Output:
		area = "OUTPUT:      "
	case Marker:
		area = "MB:          "
	case DataBlock:
		area = "DB:          "
	default:
		area = "UNKNOWN      "
	}
	return
}
