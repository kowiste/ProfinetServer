package main

import (
	"encoding/binary"
	"log"
	"os"
	"time"

	"github.com/robinson/gos7"
)

func main() {
	//Testing the connection
	handler := gos7.NewTCPClientHandler("127.0.0.1", 0, 0)
	handler.Timeout = 5 * time.Millisecond
	handler.IdleTimeout = 500 * time.Millisecond
	handler.Logger = log.New(os.Stdout, "tcp: ", log.LstdFlags)
	// Connect manually so that multiple requests are handled in one connection session
	err := handler.Connect()
	defer handler.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	client := gos7.NewClient(handler)
	buf := make([]byte, 2)
	buf[0] = 24
	buf[1] = 34
	println("SEND ", binary.BigEndian.Uint16(buf))
	err = client.AGWriteDB(13, 4, 2, buf)
	buf = make([]byte, 2)
	println("DELETE  BUFFER ", binary.BigEndian.Uint16(buf))
	err = client.AGReadDB(13, 4, 2, buf)
	println("READ ", binary.BigEndian.Uint16(buf))

}
