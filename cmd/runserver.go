package main

import (
	"time"

	profinet "github.com/Kowiste/ProfinetServer"
)

func main() {
	server := profinet.NewServer()
	server.Listen("0.0.0.0:102", 0, 0)
	for {
		time.Sleep(250 * time.Millisecond)
	}
}
