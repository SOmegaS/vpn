package client

import (
	"fmt"
	"log"
	"net"
	"vpn/internal/tun"
	"vpn/internal/vpn"
)

const BuffSize int = 65535 // UDP packet max size

type Client struct {
	iface *tun.Interface
	vpn   *vpn.Connector
	sym   bool
}

func (c *Client) setupIface() error {
	// Add IP from stdin to TUN-interface
	log.Println("INFO: Choose IP for interface")
	err := fmt.Errorf("")
	for err != nil {
		fmt.Print("Type IP for VPN net interface (default 192.168.13.1): ")
		var ip string
		_, _ = fmt.Scanln(&ip)
		if ip == "" {
			ip = "192.168.13.1"
		}
		err = c.iface.SetIp(ip, "/24")
	}
	log.Println("INFO: IP chosen")

	// Up TUN-interface
	log.Println("INFO: Upping interface")
	err = c.iface.Up()
	if err != nil {
		return fmt.Errorf("failed to up VPN net interface: %v", err)
	}
	log.Println("INFO: Upped interface")
	return nil
}

func (c *Client) Init() error {
	// Up interface
	log.Println("INFO: Setting up interface")
	err := c.setupIface()
	if err != nil {
		return fmt.Errorf("failed to init VPN interface: %v", err)
	}
	log.Println("INFO: Set up interface")

	// Resolve external ip
	log.Println("INFO: Resolving NAT type and external IP")
	eaddr, sym, err := vpn.ResolveNatIP("stun.l.google.com:19302")
	if err != nil {
		return fmt.Errorf("failed to resolve NAT type: %v", err)
	}
	c.sym = sym
	if c.sym {
		fmt.Println("You have symmetric nat (((")
	} else {
		fmt.Println("Congrats! You have not Symmetric NAT")
	}
	fmt.Printf("Your external IP is %v\n", eaddr)
	log.Printf("INFO: Resolved symmetric NAT %v, external IP %v\n", sym, eaddr)
	return nil
}

func (c *Client) Connect() error {
	log.Println("INFO: Choosing local port")
	var iaddr *net.UDPAddr
	err := fmt.Errorf("")
	for err != nil {
		fmt.Print("Specify port: ")
		var port string
		_, _ = fmt.Scanln(&port)
		iaddr, err = net.ResolveUDPAddr("udp", ":"+port)
	}
	log.Println("INFO: Chosen local port")

	// Connect to another host
	log.Println("INFO: Choosing host address")
	var raddr *net.UDPAddr
	err = fmt.Errorf("")
	for err != nil {
		fmt.Print("Type address to connect to (ip:port): ")
		var ip string
		_, _ = fmt.Scanln(&ip)
		raddr, err = net.ResolveUDPAddr("udp", ip)
	}
	log.Println("INFO: Chosen host address")

	log.Println("INFO: Connecting to host")
	conn, err := c.vpn.Connect(iaddr, raddr)
	if err != nil {
		return fmt.Errorf("failed to connect to host: %v", err)
	}
	log.Println("INFO: Connected host")

	fmt.Printf("Host %v connected to port %v\n", conn.RemoteAddr(), conn.LocalAddr().(*net.UDPAddr).Port)
	log.Printf("INFO: Host %v connected to port %v\n", conn.RemoteAddr(), conn.LocalAddr().(*net.UDPAddr).Port)
	return nil
}

func (c *Client) Serve() error {
	finish := make(chan error)

	go func() {
		buff := make([]byte, BuffSize)
		log.Println("INFO: Serve TUN")
		for {
			log.Println("INFO: Reading from TUN")
			n, err := c.iface.Read(buff)
			if err != nil {
				finish <- fmt.Errorf("failed to read from TUN: %v", err)
				return
			}
			log.Println("INFO: Read from TUN")

			log.Println("INFO: Sending to all")
			c.vpn.SendAll(buff[:n])
			log.Println("INFO: Sent to all")
		}
	}()

	go func() {
		log.Println("INFO: Serving connection")
		for {
			log.Println("INFO: Waiting message")
			buff := <-c.vpn.Input
			log.Println("INFO: Received message")

			log.Println("INFO: Writing to TUN")
			_, err := c.iface.Write(buff)
			if err != nil {
				finish <- fmt.Errorf("failed to write to TUN: %v", err)
				return
			}
			log.Println("INFO: Written to TUN")
		}
	}()

	err := <-finish
	return err
}

func NewClient() (*Client, error) {
	// Create TUN-interface
	log.Println("INFO: Creating TUN interface")
	iface, err := tun.NewInterface()
	if err != nil {
		return nil, fmt.Errorf("unable to create TUN-interface: %v", err)
	}
	log.Println("INFO: TUN interface created")

	// Create vpn connector
	log.Println("INFO: Creating VPN connector")
	connector, err := vpn.NewConnector()
	if err != nil {
		return nil, err
	}
	log.Println("INFO: VPN connector created")

	return &Client{
		iface: iface,
		vpn:   connector,
	}, nil
}
