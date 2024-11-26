package tun

import (
	"fmt"
	"github.com/songgao/water"
	"log"
	"os/exec"
)

type Interface struct {
	*water.Interface
	Ip string
}

// NewInterface creates TUN-interface. Requires sudo rights
func NewInterface() (*Interface, error) {
	// Create TUN-interface
	log.Println("INFO: Creating (parent) TUN-interface")
	tun, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		return nil, err
	}
	log.Println("INFO: (parent) TUN-interface created")

	return &Interface{
		Interface: tun,
	}, nil
}

// SetIp adds IP to TUN-interface
func (i *Interface) SetIp(ip, mask string) error {
	// IP doesn't change
	if i.Ip == ip {
		log.Println("INFO: IP is the same")
		return nil
	}

	// Delete IP if exists
	if i.Ip != "" {
		log.Println("INFO: Deleting previous ip")
		err := exec.Command("ip", "addr", "del", i.Ip+mask, "dev", i.Interface.Name()).Run()
		if err != nil {
			return fmt.Errorf("failed to remove ip from TUN-interface: %v", err)
		}
		i.Ip = ""
		log.Println("INFO: Removed previous ip")
	}

	// Add IP
	log.Println("INFO: Add ip")
	err := exec.Command("ip", "addr", "add", ip+mask, "dev", i.Interface.Name()).Run()
	if err != nil {
		return fmt.Errorf("failed to add ip to TUN-interface: %v", err)
	}
	i.Ip = ip
	log.Println("INFO: Added ip")
	return nil
}

// Up TUN-interface
func (i *Interface) Up() error {
	// Up interface
	log.Println("INFO: Up interface")
	return exec.Command("ip", "link", "set", "dev", i.Interface.Name(), "up").Run()
}
