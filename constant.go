package profinet

//Cosntant
const (
	//Protocol ID:
	ProtocolID byte = 0x32

	//MESSAGE TYPES
	JobRequest      byte = 0x01
	AckRequest      byte = 0x02
	AckDataRequest  byte = 0x03
	UserdataRequest byte = 0x07

	//REQUEST CONSTANT
	ConnectionReq byte = 0x54
	GetCPUInfo byte = 0x55
	GetCPInfo byte = 0x56

	//MEMORY AREA CONSTANT
	INFOs200                 byte = 0x03
	FLAGs200                 byte = 0x05
	AnalogInputs200          byte = 0x06
	MulipleRW                byte = 0x07 //In the constant below say is AnalogOutput200 but when read multiple send 7
	CounterS7                byte = 0x1C
	TimerS7                  byte = 0x1D
	CounterIEC               byte = 0x1E
	TimerIEC                 byte = 0x1F
	DirectPeriphericalAccess byte = 0x80
	Input                    byte = 0x81
	Output                   byte = 0x82
	Marker                   byte = 0x83
	DataBlock                byte = 0x84
	InstanceData             byte = 0x85
	LocalData                byte = 0x86
	Unknown                  byte = 0x99

	//JOB REQUEST

	CPUService         byte = 0x00
	SetupCommunication byte = 0xF0
	ReadVariable       byte = 0x04
	WriteVariable      byte = 0x05
	DownloadReq        byte = 0x1A
	DownloadBlock      byte = 0x1B
	DownloadEnd        byte = 0x1C
	StartUpload        byte = 0x1D
	Upload             byte = 0x1E
	UplodaEn           byte = 0x1F
	PLCControl         byte = 0x28
	PLCStop            byte = 0x29
)

/*
http://gmiru.com/article/s7comm-part2/
##
# Most of this is extracted from s7comm
# wireshark dissector plugin sources
# created by Thomas Wiens <th.wiens[AT]gmx.de>
# Date: 2016-15-03
# Author: Gyorgy Miru
# Version: 0.2
##

#Protocol ID:
0x32 - Protocol ID

#Message Types:
0x01 - Job Request
0x02 - Ack
0x03 - Ack-Data
0x07 - Userdata

#Header Error Class:
0x00 - No error
0x81 - Application relationship error
0x82 - Object definition error
0x83 - No ressources available error
0x84 - Error on service processing
0x85 - Error on supplies
0x87 - Access error

#Header Error Codes: (Further refines error)

#Parameter Error Codes:
0x0000 - No error
0x0110 - Invalid block type number
0x0112 - Invalid parameter
0x011A - PG ressource error
0x011B - PLC ressource error
0x011C - Protocol error
0x011F - User buffer too short
0x0141 - Request error
0x01C0 - Version mismatch
0x01F0 - Not implemented
0x8001 - L7 invalid CPU state
0x8500 - L7 PDU size error
0xD401 - L7 invalid SZL ID
0xD402 - L7 invalid index
0xD403 - L7 DGS Connection already announced
0xD404 - L7 Max user NB
0xD405 - L7 DGS function parameter syntax error
0xD406 - L7 no info
0xD601 - L7 PRT function parameter syntax error
0xD801 - L7 invalid variable address
0xD802 - L7 unknown request
0xD803 - L7 invalid request status

#Return value of item response
0x00 - Reserved
0x01 - Hardware fault
0x03 - Accessing the object not allowed
0x05 - Address out of range
0x06 - Data type not supported
0x07 - Data type inconsistent
0x0a - Object does not exist
0xff - Success

#Job Request/Ack-Data function codes
0x00 - CPU services
0xF0 - Setup communication
0x04 - Read Variable
0x05 - Write Variable
0x1A - Request download
0x1B - Download block
0x1C - Download ended
0x1D - Start upload
0x1E - Upload
0x1F - End upload
0x28 - PLC Control
0x29 - PLC Stop

#Memory Areas
0x03 - System info of S200 family
0x05 - System flags of S200 family
0x06 - Analog inputs of S200 family
0x07 - Analog outputs of S200 family
0x1C - S7 counters (C)
0x1D - S7 timers (T)
0x1E - IEC counters (200 family)
0x1F - IEC timers (200 family)
0x80 - Direct peripheral access (P)
0x81 - Inputs (I)
0x82 - Outputs (Q)
0x83 - Flags (M) (Merker)
0x84 - Data blocks (DB)
0x85 - Instance data blocks (DI)
0x86 - Local data (L)
0x87 - Unknown yet (V)

#Transport size (variable Type) in Item data
0x01 - BIT
0x02 - BYTE
0x03 - CHAR
0x04 - WORD
0x05 - INT
0x06 - DWORD
0x07 - DINT
0x08 - REAL
0x09 - DATE
0x0A - TOD
0x0B - TIME
0x0C - S5TIME
0x0F - DATE AND TIME
0x1C - COUNTER
0x1D - TIMER
0x1E - IEC TIMER
0x1F - IEC COUNTER
0x20 - HS COUNTER

#Variable ddressing mode
0x10 - S7-Any pointer (regular addressing) memory+variable length+offset
0xa2 - Drive-ES-Any seen on Drive ES Starter with routing over S7
0xb2 - S1200/S1500? Symbolic addressing mode
0xb0 - Special DB addressing for S400 (subitem read/write)

#Transport size in data
0x00 - NULL
0x03 - BIT
0x04 - BYTE/WORD/DWORD
0x05 - INTEGER
0x07 - REAL
0x09 - OCTET STRING

#Block type constants
'08' - OB
'0A' - DB
'0B' - SDB
'0C' - FC
'0D' - SFC
'0E' - FB
'0F' - SFB

#Sub block types
0x08 - OB
0x0a - DB
0x0b - SDB
0x0c - FC
0x0d - SFC
0x0e - FB
0x0f - SFB

#Block security mode
0 - None
3 - Kow How Protect

#Block Language
0x00 - Not defined
0x01 - AWL
0x02 - KOP
0x03 - FUP
0x04 - SCL
0x05 - DB
0x06 - GRAPH
0x07 - SDB
0x08 - CPU-DB DB was created from Plc programm (CREAT_DB)
0x11 - SDB (after overall reset) another SDB, don't know what it means, in SDB 1 and SDB 2, uncertain
0x12 - SDB (Routing) another SDB, in SDB 999 and SDB 1000 (routing information), uncertain
0x29 - ENCRYPT  block is encrypted (encoded?) with S7-Block-Privacy

#Userdata transmission type
0x0 - Push cyclic data push by the PLC
0x4 - Request by the master
0x8 - Response by the slave

#Userdata last PDU
0x00 - Yes
0x01 - No

#Userdata Functions
0x1 - Programmer commands
0x2 - Cyclic data
0x3 - Block functions
0x4 - CPU functions
0x5 - Security
0x7 - Time functions

#Variable table type of data
0x14 - Request
0x04 - Response

#VAT area and length type
0x01 - MB
0x02 - MW
0x03 - MD
0x11 - IB
0x12 - IW
0x13 - ID
0x21 - QB
0x22 - QW
0x23 - QD
0x31 - PIB
0x32 - PIW
0x33 - PID
0x71 - DBB
0x72 - DBW
0x73 - DBD
0x54 - TIMER
0x64 - COUNTER

#Userdata programmer subfunctions
0x01 - Request diag data (Type 1)
0x02 - VarTab
0x0c - Erase
0x0e - Read diag data
0x0f - Remove diag data
0x10 - Forces
0x13 - Request diag data (Type2)

#Userdata cyclic data subfunctions
0x01 - Memory
0x04 - Unsubscribe

#Userdata block subfunctions
0x01 - List blocks
0x02 - List blocks of type
0x03 - Get block info

#Userdata CPU subfunctions
0x01 - Read SZL
0x02 - Message service
0x03 - Transition to stop
0x0b - Alarm was acknowledged in HMI/SCADA 1
0x0c - Alarm was acknowledged in HMI/SCADA 2
0x11 - PLC is indicating a ALARM message
0x13 - HMI/SCADA initiating ALARM subscription


#Userdata security subfunctions
0x01 - PLC password

#Userdata time subfunctions
0x01 - Read clock
0x02 - Set clock
0x03 - Read clock (following)
0x04 - Set clock

#Flags for LID access
0x2 - Encapsulated LID
0x3 - Encapsulated Index
0x4 - Obtain by LID
0x5 - Obtain by Index
0x6 - Part Start Address
0x7 - Part Length

#TIA 1200 area names
0x8a0e - DB
0x0000 - IQMCT
0x50 - Inputs (I)
0x51 - Outputs (Q)
0x52 - Flags (M)
0x53 - Counter (C)
0x54 - Timer (T)
*/
