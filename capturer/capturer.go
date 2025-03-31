package capturer

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

type CaptureManager struct {
	capDevice *pcap.Interface
	capture   *PacketCapture
	Running   bool
}

func (c *CaptureManager) Init() {
	c.setDeviceName()
	c.Running = false
}

func (c *CaptureManager) StartCapture() error {
	var err error
	c.capture, err = NewPacketCapture(c.capDevice.Name)
	if err != nil {
		return err
	}

	c.Running = true
	c.capture.Start()

	go c.WritePacketToFile("capture.pcap", c.capture.packetChan)
	return nil
}

func (c *CaptureManager) WritePacketToFile(filename string, packetChan <-chan gopacket.Packet) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := pcapgo.NewWriter(f)
	err = w.WriteFileHeader(1600, layers.LinkTypeEthernet)
	if err != nil {
		return err
	}

	for packet := range packetChan {
		err := w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *CaptureManager) EndCapture() {
	c.capture.Stop()
	c.Running = false
}

func (c *CaptureManager) PrintCli() {
	fmt.Println(" -> You're Within the Capture Menu")
	if c.Running {
		fmt.Println("Capture Running...")
	} else {
		fmt.Println("")
	}
	c.PrintTargetDevice()
}

func (c *CaptureManager) PrintTargetDevice() {
	fmt.Println("\nName:", c.capDevice.Name)
	fmt.Println("Description:", c.capDevice.Description)
	fmt.Println("- IP address:", c.capDevice.Addresses[0].IP)
	fmt.Println("- Subnet Mask: ", c.capDevice.Addresses[0].Netmask)
}

func (c *CaptureManager) PrintDevices() {
	devices, _ := pcap.FindAllDevs()

	for _, device := range devices {
		fmt.Println("\n Name:", device.Name, "\n Description:", device.Description, "\n Flags:", device.Flags)
		for _, address := range device.Addresses {
			fmt.Println("- IP address:", address.IP)
			fmt.Println("- Subnet Mask: ", address.Netmask)
		}
	}
}

func (c *CaptureManager) setDeviceName() {
	devices, _ := pcap.FindAllDevs()

	for _, device := range devices {
		if strings.Contains(device.Description, "Wireless") ||
			strings.Contains(device.Description, "Wi-Fi") {
			c.capDevice = &device
			return
		}
	}
}
