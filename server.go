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
	onConnection (func(net.Addr))
}

//NewServer Create a new Profinet server
func NewServer() *Server {
	s := new(Server)
	go s.handler()
	return s
}

func (s Server) handler() {

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
	for {
		conn, err := listen.Accept()
		log.Println("Active conection")
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil
			}
			log.Printf("Unable to accept connections: %#v\n", err)
			return err
		}
		if s.onConnection != nil {
			go s.onConnection(conn.RemoteAddr())
		}
		go func(conn net.Conn) {
			defer conn.Close()

			for {
				packet := make([]byte, tcpMaxLength)
				nbytes, err := conn.Read(packet)
				if err != nil {
					if err != io.EOF {
						log.Printf("read error %v\n", err)
					}
					return
				}
				log.Println(packet[:nbytes])
				packet[3] = 22   //isoHSize
				packet[5] = 0xD0 // CC Connection confirm
				if _, err = conn.Write(packet); err != nil {
					return
				}
				println("send")
			}
		}(conn)
	}
}

//OnConnectionHandler Function that happend when there is a new conection
func (s *Server) OnConnectionHandler(function func(net.Addr)) {
	s.onConnection = function
}
