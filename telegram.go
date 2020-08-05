package profinet

func getTelegram(id int) ([]byte, error) {
	out := make([]byte, 0)
	out = append(out, 0x03, 0x00) //RFC 1006 ID (3) , Reserved
	out = append(out, 0x00, 22)   //packet lenght
	return out, nil
}
