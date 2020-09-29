# ProfinetServer PLC Simulator Siemens

This program simulate a Profinet connection, I create it to test with a program that use the gos7 library https://github.com/robinson/gos7.
On cmd/client there is a small example of client connect to this server, try cmd/runserver.go and later uncomment the line 15 to see how he value is overwrite.


Is not completely functional.

In this structure are the "Event" that should call when try to read or write in  that part of the memory, or the timer to modify every second or whatever time.

	onConnection   (func(net.Addr))                //On Connection handler
	onCounterRead  (func(*Server) ([]byte, error)) //On Read Counter handler
	onTimerRead    (func(*Server) ([]byte, error)) //On Read Timer handler
	onInputRead    (func(*Server) ([]byte, error)) //On Read Input handler
	onOutputRead   (func(*Server) ([]byte, error)) //On Read Output handler
	onMBRead       (func(*Server) ([]byte, error)) //On Read MB handler
	onDBRead       (func(*Server) ([]byte, error)) //On Read DB handler
	onMultiRead    (func())                        //On Multi Read handler
	onCounterWrite (func())                        //On Write Counter handler
	onTimerWrite   (func())                        //On Write Timer handler
	onInputWrite   (func())                        //On Write Input handler
	onOutputWrite  (func())                        //On Write Output handler
	onMBWrite      (func())                        //On Write MB handler
	onDBWrite      (func())                        //On Write DB handler
	onMultiWrite   (func())                        //On Multi Write handler
	onTimer        (func(*Server))                 //On time handler
  
  
The array of byte have  a specific lenght which you can read on s.PDULenght.
