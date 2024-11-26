package vpn

import (
	"fmt"
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
	log.Println("INFO: Sending started")
	for _, conn := range c.conns {
		log.Println("INFO: Sending message to", conn.RemoteAddr())
		_, err := conn.Write(msg)
		if err != nil {
			log.Printf("ERROR: Write message error to host %v: %v", conn.RemoteAddr(), err)
			continue
		}
		log.Println("INFO: Sent message to", conn.RemoteAddr())
	}
	log.Println("INFO: Sending finished")
}

func (c *Connector) handler(conn *net.UDPConn) {
	buff := make([]byte, BuffSize)
	log.Println("INFO: Handling from", conn.RemoteAddr())
	for {
		n, err := conn.Read(buff)
		if err != nil {
			log.Printf("ERROR: Read message error from %v: %v", conn.RemoteAddr(), err)
			log.Printf("INFO: Connection to %v closed\n", conn.RemoteAddr())
			return
		}
		if string(buff[:n]) == "keep-alive" {
			log.Println("INFO: Received keep-alive from", conn.RemoteAddr())
			continue
		}
		c.Input <- buff[:n]
		log.Println("INFO: Read message from", conn.RemoteAddr())
	}
}

func (c *Connector) handshake(conn *net.UDPConn) error {
	var n int
	err := fmt.Errorf("")
	buff := make([]byte, BuffSize)
	for err != nil {
		log.Println("INFO: Sending hello to", conn.RemoteAddr())
		_, err = conn.Write([]byte("hello"))
		if err != nil {
			return err
		}
		log.Println("INFO: Sent hello to", conn.RemoteAddr())

		log.Println("INFO: Waiting hello or ready from", conn.RemoteAddr())
		n, err = conn.Read(buff)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
		}
	}
	log.Println("INFO: Read message from", conn.RemoteAddr())

	if string(buff[:n]) == "hello" {
		log.Println("INFO: Received hello from", conn.RemoteAddr())
	} else if string(buff[:n]) == "ready" {
		log.Println("INFO: Received ready from", conn.RemoteAddr())
	} else {
		return fmt.Errorf("unrecognized hanshake command: %v from %v", string(buff[:n]), conn.RemoteAddr())
	}

	log.Println("INFO: Sending ready to", conn.RemoteAddr())
	_, err = conn.Write([]byte("ready"))
	if err != nil {
		return err
	}
	log.Println("INFO: Sent ready to", conn.RemoteAddr())

	if string(buff[:n]) == "ready" {
		return nil
	}

	log.Println("INFO: Waiting ready from", conn.RemoteAddr())
	err = fmt.Errorf("")
	for err != nil {
		n, err = conn.Read(buff)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
		}
	}
	log.Println("INFO: Read message from", conn.RemoteAddr())

	if string(buff[:n]) == "ready" {
		log.Println("INFO: Received ready from", conn.RemoteAddr())
	} else {
		return fmt.Errorf("unrecognized hanshake command: %v from %v", string(buff[:n]), conn.RemoteAddr())
	}

	return nil
}

func (c *Connector) Connect(iaddr, raddr *net.UDPAddr) (*net.UDPConn, error) {
	log.Println("INFO: Connecting to", raddr)
	log.Printf("INFO: Dial from %v to %v\n", iaddr, raddr)
	conn, err := net.DialUDP("udp", iaddr, raddr)
	if err != nil {
		return nil, err
	}
	log.Println("INFO: Created dial to", conn.RemoteAddr())

	log.Println("INFO: Handshaking with", conn.RemoteAddr())
	err = c.handshake(conn)
	if err != nil {
		return nil, err
	}
	log.Println("INFO: Successful handshake with", conn.RemoteAddr())

	log.Println("INFO: Staring handler", conn.RemoteAddr())
	go c.handler(conn)
	log.Println("INFO: Starting keep-alive ping")
	go func() {
		for {
			time.Sleep(15 * time.Second) // TODO подобрать время
			log.Println("INFO: Sending keep-alive to", conn.RemoteAddr())
			_, err := conn.Write([]byte("keep-alive"))
			if err != nil {
				return
			}
			log.Println("INFO: Sent keep-alive to", conn.RemoteAddr())
		}
	}()
	c.conns = append(c.conns, conn)
	return conn, nil
}

func NewConnector() (*Connector, error) {
	return &Connector{
		conns: make([]*net.UDPConn, 0),
		Input: make(chan []byte, BuffSize),
	}, nil
}
