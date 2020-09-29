package profinet

import (
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
)

func (s *Server) getTelegram(rq []byte) ([]byte, error) {
	out := make([]byte, 0)
	s.getCommand(rq) //Read the command, data lenght and other parameter

	out = append(out, s.getTPKHeader(rq)...) //b0-3//TPKT HEADER

	if s.command == ConnectionReq { //ISO CONNECION REQUEST
		log.Println("ISO REQUEST:      " + fmt.Sprint(rq))
		out = append(out, s.handlerConnection(rq)...) //che
		return out, nil
	}

	out = append(out, s.getS7PDUHeader()...) //S7PDU HEADER (10-12 bytes)

	if s.command == SetupCommunication { //NEGOTIATIN REQUEST
		out[3] = 27                                     //packetlenght
		out = append(out, make([]byte, 11)...)          //b17-24
		out = append(out, uint16ToByte(s.PDULenght)...) //b25-26 //PDULenght
		return out, nil
	}

	if s.command == ReadVariable { //command == ReadCounterReq || command == ReadTimerReq || command == ReadInputReq || command == ReadOutputReq || command == ReadMBReq || command == ReadDBReq || command == ReadMulti { //READ REQUET
		out = append(out, s.getVariable()...)
		out[20] = 1
		out[3] = 25 + byte(s.PDULenght) //update packet lenght TODO both byte 2nd 3
		return out, nil
	}

	if s.command == WriteVariable { //command == WriteCounterReq || command == WriteTimerReq || command == WriteInputReq || command == WriteOutputReq || command == WriteMBReq || command == WriteDBReq || command == WriteMulti { //WRITE REQUEST
		out = append(out, s.setVariable(rq)...)
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
	out = append(out, 0x03, 0x00, rq[2], rq[3])
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
		s.PDULenght = binary.BigEndian.Uint16([]byte{rq[23], rq[24]})
		s.varType = rq[22]
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
	case SetupCommunication:
		req = "NEGOTIATION:      "
	case ReadVariable, WriteVariable:
		if s.command == ReadVariable {
			req = "READ "
		} else {
			req = "WRITE "
		}
		s.rwcomm = rq[27]
		s.addressReq.DB = binary.BigEndian.Uint16([]byte{rq[25], rq[26]})
		s.addressReq.Address = (binary.BigEndian.Uint32([]byte{0x00, rq[28], rq[29], rq[30]})) >> 3
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
func (s *Server) handlerConnection(rq []byte) []byte {

	DstTSAP := binary.BigEndian.Uint16([]byte{rq[20], rq[21]}) // Dst TSAP
	dstslot := DstTSAP & 0x001F                                //getting slot
	dstrack := ((DstTSAP >> 5) - 8) & 0x001F                   //getting rack
	log.Println("slot ", strconv.Itoa(int(dstslot)))
	log.Println("rack ", strconv.Itoa(int(dstrack)))
	out := make([]byte, 0)
	out = append(out, s.getISOCOPT()...)
	if s.rack != dstrack || s.slot != dstslot {
		out[1] = 0x00
	}
	out = append(out, make([]byte, int(s.reqDataLen-isoHSize))...)
	return out
}

func (s *Server) getVariable() []byte {
	out := make([]byte, 0)
	out = append(out, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00) //b17
	out = append(out, 0xFF)                                     //b21			//ErrorCode
	out = append(out, s.varType)                                //b22 			//Variable type see table in constant
	out = append(out, uint16ToByte(s.PDULenght)...)             //b23-24		//Count
	switch s.rwcomm {
	case Input:
		out = append(out, s.getMemory(s.input)...)
	case Output:
		out = append(out, s.getMemory(s.output)...)
	case Marker:
		if s.onMBRead != nil {
			outEvent, _ := s.onMBRead(s)
			out = append(out, outEvent...)
		} else {
			out = append(out, s.getMemory(s.marker)...)
		}
	case DataBlock:
		if s.onDBRead != nil {
			outEvent, _ := s.onDBRead(s)
			out = append(out, outEvent...)
		} else {
			out = append(out, s.getMemory(s.db[int(s.addressReq.DB)])...)
		}

	}
	return out
}
func (s *Server) setVariable(rq []byte) []byte {
	out := make([]byte, 0)
	out = append(out, 0x00)      //b17			//Specification type fr const 18 for read/write
	out = append(out, 0x00)      //b18			// Lenght rest of byte
	out = append(out, 22)        //b19			//Syntax ID const 16 for any typ addr
	out = append(out, s.varType) //b20			//Variable type see table below
	out = append(out, 0x00)      //b21			//
	out = append(out, 0x00, 0x00)
	if s.setMemory(rq) {
		out = append(out, 0xFF) //b22			ok write
	} else {
		out = append(out, 0x00) //b22			Bad write
	}

	return out
}

func (s *Server) getMemory(mem []byte) (out []byte) {
	out = make([]byte, 0)
	for element := 0; element < int(s.PDULenght); element++ {
		out = append(out, mem[int(s.addressReq.Address*2)+element])
	}
	return
}
func (s *Server) setMemory(rq []byte) bool {
	storeAddr := int(s.addressReq.Address * 2)
	for addr := 0; addr < int(s.PDULenght); addr++ {
		switch s.rwcomm {
		case Input:
			s.input[storeAddr+addr] = rq[35+addr]
		case Output:
			s.output[storeAddr+addr] = rq[35+addr]
		case Marker:
			s.marker[storeAddr+addr] = rq[35+addr]
		case DataBlock:
			s.db[int(s.addressReq.DB)][storeAddr+addr] = rq[35+addr]
		}
	}
	return true
}
