package profinet

import (
	"net"
	"time"
)

//SetDB set DB data
func (s *Server) SetDB(DBNum int, data []uint16) {

	s.db[DBNum] = arrUint16ToByte(data)
}

//SetMB set MB data
func (s *Server) SetMB(mark []byte) {
	s.marker = mark
}

//SetTimer set MB data
func (s *Server) SetTimer(input []uint16) {
	s.timer = arrUint16ToByte(input)
}

//SetCounter set MB data
func (s *Server) SetCounter(input []uint16) {
	s.counter = arrUint16ToByte(input)
}

//SetInput set input
func (s *Server) SetInput(data []uint16) {
	s.input = arrUint16ToByte(data)
}

//SetOutput set input
func (s *Server) SetOutput(data []uint16) {
	s.output = arrUint16ToByte(data)
}

//OnConnectionHandler Function that happend when there is a new nection
func (s *Server) OnConnectionHandler(function func(net.Addr)) {
	s.onConnection = function
}

//OnCounterReadHandler Function that happend when read a counter
func (s *Server) OnCounterReadHandler(function func(*Server) ([]byte, error)) {
	s.onCounterRead = function
}

//OnTimerReadHandler Function that happend when read a timer
func (s *Server) OnTimerReadHandler(function func(*Server) ([]byte, error)) {
	s.onTimerRead = function
}

//OnInputReadHandler Function that happend when read a input
func (s *Server) OnInputReadHandler(function func(*Server) ([]byte, error)) {
	s.onTimerRead = function
}

//OnOutputReadHandler Function that happend when read a output
func (s *Server) OnOutputReadHandler(function func(*Server) ([]byte, error)) {
	s.onOutputRead = function
}

//OnMBReadHandler Function that happend when read a MB
func (s *Server) OnMBReadHandler(function func(*Server) ([]byte, error)) {
	s.onMBRead = function
}

//OnDBReadHandler Function that happend when read a DB
func (s *Server) OnDBReadHandler(function func(*Server) ([]byte, error)) {
	s.onDBRead = function
}

//OnMultiReadHandler Function that happend when multi read
func (s *Server) OnMultiReadHandler(function func()) {
	s.onMultiRead = function
}

//OnCounterWriteHandler Function that happend when write a counter
func (s *Server) OnCounterWriteHandler(function func()) {
	s.onCounterWrite = function
}

//OnTimerWriteHandler Function that happend when write a timer
func (s *Server) OnTimerWriteHandler(function func()) {
	s.onTimerWrite = function
}

//OnInputWriteHandler Function that happend when write a input
func (s *Server) OnInputWriteHandler(function func()) {
	s.onInputWrite = function
}

//OnOutputWriteHandler Function that happend when write a output
func (s *Server) OnOutputWriteHandler(function func()) {
	s.onOutputWrite = function
}

//OnMBWriteHandler Function that happend when write a MB
func (s *Server) OnMBWriteHandler(function func()) {
	s.onMBWrite = function
}

//OnDBWriteHandler Function that happend when write a DB
func (s *Server) OnDBWriteHandler(function func()) {
	s.onDBWrite = function
}

//OnMultiWriteHandler Function that happend when multi rite
func (s *Server) OnMultiWriteHandler(function func()) {
	s.onMultiWrite = function
}

//OnTimerHandler Function that happend when the interval occure
func (s *Server) OnTimerHandler(function func(*Server), tick time.Duration) {
	s.onTimer = function
	go s.timing(tick)
}
func (s *Server) timing(tick time.Duration) {
	t := time.NewTicker(tick)
	for {
		<-t.C //Wait for the time
		s.onTimer(s)
	}
}
