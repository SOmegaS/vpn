package client

import (
	"fmt"
	"log"
	"net"
	"time"
	"vpn/internal/tun"
	"vpn/internal/vpn"
)

type Client struct {
	iface *tun.Interface
	vpn   *vpn.Connector
	sym   bool
}

func (c *Client) upIface() error {
	// Add IP from stdin to TUN-interface
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

	// Up TUN-interface
	err = c.iface.Up()
	if err != nil {
		return fmt.Errorf("failed to up VPN net interface: %v", err)
	}
	return nil
}

func (c *Client) Init() error {
	// Up interface
	err := c.upIface()
	if err != nil {
		return fmt.Errorf("failed to init VPN interface: %v", err)
	}

	// Resolve external ip
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

	return nil
}

func (c *Client) Listen() error {
	// Specify local port
	var iaddr *net.UDPAddr
	err := fmt.Errorf("")
	for err != nil {
		fmt.Print("Specify port or empty for random: ")
		var port string
		_, _ = fmt.Scanln(&port)
		iaddr, err = net.ResolveUDPAddr("udp", ":"+port)
	}

	conn, err := c.vpn.Listen(iaddr)
	if err != nil {
		return fmt.Errorf("failed to connect to host: %v", err)
	}

	fmt.Printf("Host %v connected to port %v\n", conn.RemoteAddr(), conn.LocalAddr().(*net.UDPAddr).Port)

	return nil
}

func (c *Client) Connect() error {
	// Connect to another host
	var raddr *net.UDPAddr
	err := fmt.Errorf("")
	for err != nil {
		fmt.Print("Type address to connect to (ip:port): ")
		var ip string
		_, _ = fmt.Scanln(&ip)
		raddr, err = net.ResolveUDPAddr("udp", ip)
	}

	conn, err := c.vpn.Connect(raddr)
	if err != nil {
		return fmt.Errorf("failed to connect to host: %v", err)
	}

	fmt.Printf("Connected to port %v\n", conn.LocalAddr().(*net.UDPAddr).Port)
	return nil
}

func (c *Client) Serve() error {
	for range 10 {
		c.vpn.SendAll([]byte("HellO!!"))
		select {
		case buff := <-c.vpn.Input:
			log.Println("INFO: Received message")
			fmt.Println(string(buff))
		default:
			log.Println("INFO: Nothing received")
		}
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func NewClient() (*Client, error) {
	// Create TUN-interface
	iface, err := tun.NewInterface()
	if err != nil {
		return nil, fmt.Errorf("unable to create TUN-interface: %v", err)
	}
	// Create vpn connector
	connector, err := vpn.NewConnector()
	if err != nil {
		return nil, err
	}
	return &Client{
		iface: iface,
		vpn:   connector,
	}, nil
}
