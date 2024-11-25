package vpn

import (
	"log"
	"net"
)

const BuffSize int = 65535 // UDP packet max size

type Connector struct {
	conns []*net.UDPConn
	Input chan []byte
}

func (c *Connector) SendAll(msg []byte) {
	log.Println("INFO: Sending started")
	for _, conn := range c.conns {
		log.Println("INFO: Sending message to ", conn.RemoteAddr())
		_, err := conn.Write(msg)
		if err != nil {
			log.Printf("ERROR: Write message error to host %v: %v", conn.RemoteAddr(), err)
			continue
		}
		log.Println("INFO: Sent message to ", conn.RemoteAddr())
	}
	log.Println("INFO: Sending finished")
}

func (c *Connector) handler(conn *net.UDPConn) {
	buff := make([]byte, BuffSize)
	log.Println("INFO: Handling from ", conn.RemoteAddr())
	for {
		n, err := conn.Read(buff)
		if err != nil {
			// TODO check if connection closed
			log.Printf("ERROR: Read message error from %v: %v", conn.RemoteAddr(), err)
			continue
			//log.Printf("INFO: Connection to %v closed\n", conn.RemoteAddr())
			//return
		}
		c.Input <- buff[:n]
		log.Println("INFO: Read message from ", conn.RemoteAddr())
	}
}

func (c *Connector) connect(iaddr, raddr *net.UDPAddr) (*net.UDPConn, error) {
	log.Printf("INFO: Dial from %v to %v\n", iaddr, raddr)
	conn, err := net.DialUDP("udp", iaddr, raddr)
	if err != nil {
		return nil, err
	}
	log.Println("INFO: Created dial to", conn.RemoteAddr())
	log.Println("INFO: Staring handler ", conn.RemoteAddr())
	go c.handler(conn)
	c.conns = append(c.conns, conn)
	return conn, nil
}

func (c *Connector) Listen(iaddr *net.UDPAddr) (*net.UDPConn, error) {
	log.Println("INFO: Listening on", iaddr)
	conn, err := net.ListenUDP("udp", iaddr)
	if err != nil {
		return nil, err
	}
	log.Println("INFO: Created listener on", conn.LocalAddr())

	log.Println("INFO: Waiting connection on", conn.LocalAddr())
	_, raddr, err := conn.ReadFromUDP(nil)
	if err != nil {
		return nil, err
	}
	log.Println("INFO: Received connection from", raddr)

	log.Println("INFO: Closing listener on", conn.LocalAddr())
	err = conn.Close()
	if err != nil {
		return nil, err
	}
	log.Println("INFO: Listener closed on", conn.LocalAddr())

	log.Printf("INFO: Connecting from %v to %v\n", iaddr, raddr)
	return c.connect(iaddr, raddr)
}

func (c *Connector) Connect(raddr *net.UDPAddr) (*net.UDPConn, error) {
	log.Println("INFO: Connecting to", raddr)
	return c.connect(nil, raddr)
}

func NewConnector() (*Connector, error) {
	return &Connector{
		conns: make([]*net.UDPConn, 0),
		Input: make(chan []byte, BuffSize),
	}, nil
}
