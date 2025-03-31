package capturer

import (
	"fmt"
	"strings"

	"github.com/google/gopacket/pcap"
)

type Capturer struct {
	capDevice *pcap.Interface
	Running   bool
}

func (c *Capturer) Init() {
	c.setDeviceName()
	c.Running = false
}

func (c *Capturer) StartCapture() {
	c.Running = true
}

func (c *Capturer) EndCapture() {
	c.Running = false
}

func (c *Capturer) PrintCli() {
	fmt.Println(" -> You're Within the Capture Menu")
	if c.Running {
		fmt.Println("Capture Running...")
	} else {
		fmt.Println("")
	}
	c.PrintTargetDevice()
}

func (c *Capturer) PrintTargetDevice() {
	fmt.Println("\nName:", c.capDevice.Name)
	fmt.Println("Description:", c.capDevice.Description)
	fmt.Println("- IP address:", c.capDevice.Addresses[0].IP)
	fmt.Println("- Subnet Mask: ", c.capDevice.Addresses[0].Netmask)
}

func (c *Capturer) PrintDevices() {
	devices, _ := pcap.FindAllDevs()

	for _, device := range devices {
		fmt.Println("\n Name:", device.Name, "\n Description:", device.Description, "\n Flags:", device.Flags)
		for _, address := range device.Addresses {
			fmt.Println("- IP address:", address.IP)
			fmt.Println("- Subnet Mask: ", address.Netmask)
		}
	}
}

func (c *Capturer) setDeviceName() {
	devices, _ := pcap.FindAllDevs()

	for _, device := range devices {
		if strings.Contains(device.Description, "Wireless") ||
			strings.Contains(device.Description, "Wi-Fi") {
			c.capDevice = &device
			return
		}
	}
}
