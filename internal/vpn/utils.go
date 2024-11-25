package vpn

import (
	"crypto/rand"
	"net"
)

func makeSTUNMsg() ([]byte, error) {
	message := make([]byte, 20)
	header := []byte{0, 1, 0, 0, 33, 18, 164, 66}
	copy(message, header)
	_, err := rand.Read(message[len(header):])
	if err != nil {
		return nil, err
	}
	return message, nil
}

// ResolveNatIP Creates STUN request from specified port
func ResolveNatIP(stunUri string) (eaddr *net.IPAddr, sym bool, err error) {
	msg, err := makeSTUNMsg()
	if err != nil {
		return nil, false, err
	}
	raddr, err := net.ResolveUDPAddr("udp", stunUri)
	if err != nil {
		return nil, false, err
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, false, err
	}
	defer func(conn *net.UDPConn) {
		err = conn.Close()
	}(conn)
	_, err = conn.Write(msg)
	if err != nil {
		return nil, false, err
	}
	// TODO верификацию ответа: код ответа, считанная длина, тот же transaction id и т.д.
	buff := make([]byte, 32)
	_, err = conn.Read(buff)
	if err != nil {
		return nil, false, err
	}
	// XOR with magic cookie (stun)
	eaddr = &net.IPAddr{
		IP: []byte{0, 0, 0, 0},
	}
	port := int(buff[26]^buff[4])<<8 + int(buff[27]^buff[5])
	eaddr.IP[0] = buff[28] ^ buff[4]
	eaddr.IP[1] = buff[29] ^ buff[5]
	eaddr.IP[2] = buff[30] ^ buff[6]
	eaddr.IP[3] = buff[31] ^ buff[7]
	return eaddr, port != conn.LocalAddr().(*net.UDPAddr).Port, err
}
