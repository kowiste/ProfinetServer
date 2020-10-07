package main

import (
	"log"
	"time"

	profinet "github.com/Kowiste/ProfinetServer"
)

func main() {
	server := profinet.NewServer()
	server.SetOutput([]uint16{23, 563, 464, 5, 7856, 45, 2345, 6, 7, 535})
	server.SetInput([]uint16{2456, 876, 23, 2245, 675, 86, 97, 2134, 5, 47})
	server.SetDB(13, []uint16{11, 22, 33, 44, 55, 66, 77, 88, 99, 100})
	//server.OnDBReadHandler(handlerDB)
	server.OnTimerHandler(handlerTimer, 1*time.Second)
	server.Listen("0.0.0.0:102", 0, 1)
	for {
		time.Sleep(1000 * time.Millisecond)
	}
}

func handlerDB(s *profinet.Server) ([]byte, error) {
	return []byte{0, 34}, nil
}

func handlerTimer(s *profinet.Server) {
	log.Println("the time pass...")
}
