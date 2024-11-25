package vpn

import (
	"log"
	"net"
	"time"
)

const BuffSize int = 65535 // UDP packet max size

type Connector struct {
	conns []*net.UDPConn
	Input chan []byte
}

func (c *Connector) SendAll(msg []byte) {
	for _, conn := range c.conns {
		_, err := conn.Write(msg)
		if err != nil {
			log.Printf("ERROR: Write message error to host %v: %v", conn.RemoteAddr(), err)
			continue
		}
	}
}

func (c *Connector) handler(conn *net.UDPConn) {
	time.Sleep(3 * time.Second)
	buff := make([]byte, BuffSize)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			log.Printf("ERROR: Read message error from %v: %v", conn.RemoteAddr(), err)
			//log.Printf("INFO: Connection to %v closed\n", conn.RemoteAddr())
			//return
		}
		c.Input <- buff[:n]
	}
}

func (c *Connector) connect(iaddr, raddr *net.UDPAddr) (*net.UDPConn, error) {
	conn, err := net.DialUDP("udp", iaddr, raddr)
	if err != nil {
		return nil, err
	}
	go c.handler(conn)
	c.conns = append(c.conns, conn)
	return conn, nil
}

func (c *Connector) Listen(iaddr *net.UDPAddr) (*net.UDPConn, error) {
	conn, err := net.ListenUDP("udp", iaddr)
	if err != nil {
		return nil, err
	}

	_, raddr, err := conn.ReadFromUDP(nil)
	if err != nil {
		return nil, err
	}

	err = conn.Close()
	if err != nil {
		return nil, err
	}

	return c.connect(iaddr, raddr)
}

func (c *Connector) Connect(raddr *net.UDPAddr) (*net.UDPConn, error) {
	return c.connect(nil, raddr)
}

func NewConnector() (*Connector, error) {
	return &Connector{
		conns: make([]*net.UDPConn, 0),
		Input: make(chan []byte, BuffSize),
	}, nil
}
