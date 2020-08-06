package profinet

import (
	"io"
	"log"
	"net"
	"strings"
)

const (
	tcpMaxLength = 2084
)

//Server structure for the server
type Server struct {
	listeners    []net.Listener
	conn         net.Conn
	onConnection (func(net.Addr))
}

//NewServer Create a new Profinet server
func NewServer() *Server {
	s := new(Server)
	return s
}

//Listen Listen profinet
func (s Server) Listen(IP string) error {
	listen, err := net.Listen("tcp", IP)
	if err != nil {
		log.Printf("Failed to Listen: %v\n", err)
		return err
	}
	log.Println("Server Active on: ", IP)
	s.listeners = append(s.listeners, listen)
	go s.accept(listen)
	return err
}
func (s *Server) accept(listen net.Listener) error {
	var err error
	for {
		s.conn, err = listen.Accept()
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
		go s.handler()
	}
}

func (s Server) handler() {
	defer s.conn.Close()
	packet := make([]byte, tcpMaxLength)
	for {
		//log.Println("Wait command")
		nbytes, err := s.conn.Read(packet)
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error %v\n", err)
			}
			return
		}
		log.Println(packet[:nbytes])
		if _, err = s.conn.Write(createbytes(packet)); err != nil {
			log.Println("Error writing")
			return
		}
		//log.Println("SEND PACKET")
	}
}

//OnConnectionHandler Function that happend when there is a new nection
func (s *Server) OnConnectionHandler(function func(net.Addr)) {
	s.onConnection = function
}
func createbytes(r []byte) []byte {
	out := []byte{3, 0, 0}
	for i := 0; i < int(r[3]-3); i++ {
		out = append(out, 0)
	}
	out[3] = r[3]
	out[5] = 0xD0
	if r[5] == 240 && r[11] != 5 {
		out[3] = 27
		out = append(out, 0, 0)
		out[25] = 0
		out[26] = 20
		log.Println("Send negotation ack")
	} else if r[11] == 5 {
		out[21] = 0xFF
		out[25] = 9
		out[26] = 10
		//log.Println("Send read")
	} else {
		log.Println("Send connection ack")
	}
	return out
}
