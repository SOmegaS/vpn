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
	log.Println("INFO: Create TUN-interface")
	tun, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		return nil, err
	}
	return &Interface{
		Interface: tun,
	}, nil
}

// SetIp adds IP to TUN-interface
func (i *Interface) SetIp(ip, mask string) error {
	// IP doesn't change
	if i.Ip == ip {
		return nil
	}

	// Delete IP if exists
	if i.Ip != "" {
		err := exec.Command("ip", "addr", "del", i.Ip+mask, "dev", i.Interface.Name()).Run()
		if err != nil {
			return fmt.Errorf("failed to remove ip from TUN-interface: %v", err)
		}
		i.Ip = ""
	}

	// Add IP
	log.Println("INFO: Add ip")
	err := exec.Command("ip", "addr", "add", ip+mask, "dev", i.Interface.Name()).Run()
	if err != nil {
		return fmt.Errorf("failed to add ip to TUN-interface: %v", err)
	}
	i.Ip = ip
	return nil
}

// Up TUN-interface
func (i *Interface) Up() error {
	// Up interface
	log.Println("INFO: Up interface")
	return exec.Command("ip", "link", "set", "dev", i.Interface.Name(), "up").Run()
}
