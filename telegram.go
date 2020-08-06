package profinet

func getTelegram( rq []byte) ([]byte, error) {
	out := make([]byte, 0)
	//TPKT HEADER RFC 1006
	out = append(out, 0x03, 0x00) //b0-1		//pvrsn (always 3), Reserved
	out = append(out, 0x00, 0x00) //b2-3		//packet lenght [High,Low]
	//ISO-COTP
	out = append(out, 0x02, 0xF0, 0x80) //b2-4	// seen to be always 2,240,128
	//S7PDU
	//HEADER (10-12 bytes)
	out = append(out, 0x32)           //b5		//Protocol id (Always 0x32:50)
	out = append(out, 0x01)           //b6		//Message type 0x01(job request), 0x02 (Ack), 0x03 (Ack Data), 0x07 (UserData)
	out = append(out, 0x00, 0x00)     //b7-8	//Reserved
	out = append(out, rq[11], rq[12]) //b9-10	//PDU Reference  PLC just copy it to the reply
	out = append(out, 0x00, 0x00)     //b11-12	//Parameter Length
	out = append(out, 0x00, 0x00)     //b13-14	//Data Length
	out = append(out, 0x00)           //b15		//Error Class 	[optional only present in the Ack-Data message]
	out = append(out, 0x00)           //b16		//Error Code	[optional only present in the Ack-Data message]

	return out, nil
}
