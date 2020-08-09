package profinet

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const (
	tcpMaxLength     = 2084
	isoHSize         = 7
	minPduSize       = 16
	pduSizeRequested = 480
)

//Server structure for the server
type Server struct {
	listeners      []net.Listener
	conn           net.Conn
	State          bool
	msgType        byte
	command        byte
	rwcomm         byte
	PDUReference   uint16
	reqDataLen     int
	onConnection   (func(net.Addr)) //On Connection handler
	onCounterRead  (func())         //On Read Counter handler
	onTimerRead    (func())         //On Read Timer handler
	onInputRead    (func())         //On Read Input handler
	onOutputRead   (func())         //On Read Output handler
	onMBRead       (func())         //On Read MB handler
	onDBRead       (func())         //On Read DB handler
	onMultiRead    (func())         //On Multi Read handler
	onCounterWrite (func())         //On Write Counter handler
	onTimerWrite   (func())         //On Write Timer handler
	onInputWrite   (func())         //On Write Input handler
	onOutputWrite  (func())         //On Write Output handler
	onMBWrite      (func())         //On Write MB handler
	onDBWrite      (func())         //On Write DB handler
	onMultiWrite   (func())         //On Multi Write handler
	input          []uint16
	ouput          []uint16
	marker         []byte
	counter        []uint16
	timer          []uint16
	db             map[int][]uint16
}

//Address structure to save the address
type Address struct {
	DB     int
	Memory int
	Size   int
}

//NewServer Create a new Profinet server
func NewServer() *Server {
	s := new(Server)
	s.configureMemory()
	return s
}

//Listen Listen profinet
func (s Server) Listen(IP string) error {
	listen, err := net.Listen("tcp", IP)
	if err != nil {
		log.Printf("Failed to Listen: %v\n", err)
		return err
	}
	log.Println("Profinet Server Active on: ", IP)
	s.listeners = append(s.listeners, listen)
	go s.accept(listen)
	return err
}

func (s *Server) configureMemory() {
	s.counter = make([]uint16, ^uint16(0))
	s.db = make(map[int][]uint16)
	for i := 1; i < 100; i++ {
		s.db[i] = make([]uint16, ^uint16(0))
	}
	s.input = make([]uint16, ^uint16(0))
	s.ouput = make([]uint16, ^uint16(0))
	s.marker = make([]byte, ^byte(0))
	s.counter = make([]uint16, ^uint16(0))
	s.timer = make([]uint16, ^uint16(0))
}

func (s *Server) accept(listen net.Listener) error {
	var err error
	for {
		s.conn, err = listen.Accept() //Wait for a connection
		if s.onConnection != nil {
			go s.onConnection(s.conn.RemoteAddr())
			log.Println("Connection Handler")
		}
		log.Println("Active conection")
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil
			}
			log.Printf("Unable to accept connections: %#v\n", err)
			return err
		}
		if s.onConnection != nil {
			go s.onConnection(s.conn.RemoteAddr())
		}
		go s.handler() //Handler the connection
	}
}

func (s Server) handler() {
	defer s.conn.Close()
	for {
		recv, err := s.readConn()
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error %v\n", err)
			}
			return
		}
		send, _ := s.getTelegram(recv)

		if _, err = s.conn.Write(send); err != nil {
			log.Println("Error writing")
			return
		}
		log.Println("SEND PACKET:      " + fmt.Sprint(send))
	}
}

func (s *Server) readConn() (response []byte, err error) {
	done := false
	data := make([]byte, tcpMaxLength)
	length := 0
	for !done && err == nil {
		// Get TPKT (4 bytes)
		if _, err = io.ReadFull(s.conn, data[:4]); err != nil {
			log.Printf("%T %+v", err, err)
			return nil, err
		}
		// Read length, ignore transaction & protocol id (4 bytes)
		length = int(binary.BigEndian.Uint16(data[2:]))
		if length == isoHSize {
			_, err = io.ReadFull(s.conn, data[4:7])
			if err != nil { // Skip remaining 3 bytes and Done is still false
				return nil, err
			}
		} else {
			if length > pduSizeRequested+isoHSize || length < minPduSize {
				err = fmt.Errorf("s7: invalid pdu")
				return nil, err
			}
			done = true
		}
	}
	// Skip remaining 3 COTP bytes
	_, err = io.ReadFull(s.conn, data[4:7])
	if err != nil {
		return nil, err
	}
	// Receives the S7 Payload
	_, err = io.ReadFull(s.conn, data[7:length])
	if err != nil {
		return nil, err
	}
	response = data[0:length]
	return
}
