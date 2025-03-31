package capturer

import (
	"fmt"
	"strings"
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type Capturer struct {
	capDevice *pcap.Interface
	Running   bool
	capture   *PacketCapture
}

type PacketCapture struct {
	handle     *pcap.Handle
	stopCh     chan struct{}
	wg         sync.WaitGroup
	packetChan chan gopacket.Packet
}

func (pc *PacketCapture) Start() {
	pc.wg.Add(1)
	go func() {
		defer pc.wg.Done()
		defer close(pc.packetChan)
		source := gopacket.NewPacketSource(pc.handle, pc.handle.LinkType())
		for {
			select {
			case <-pc.stopCh:
				return
			case packet := <-source.Packets():
				pc.packetChan <- packet
			}
		}
	}()
}

func (pc *PacketCapture) Stop() {
	close(pc.stopCh)
	pc.wg.Wait()
	pc.handle.Close()
}

func (pc *PacketCapture) Packets() <-chan gopacket.Packet {
	return pc.packetChan
}

func (c *Capturer) Init() {
	c.setDeviceName()
	c.Running = false
}

func NewPacketCapture(device string) (*PacketCapture, error) {
	handle, err := pcap.OpenLive(device, 1600, true, pcap.BlockForever)
	if err != nil {
		return nil, fmt.Errorf("error opening device: %v", err)
	}

	return &PacketCapture{
		handle:     handle,
		stopCh:     make(chan struct{}),
		packetChan: make(chan gopacket.Packet, 1000),
	}, nil
}

func (c *Capturer) StartCapture() error {
	var err error
	c.capture, err = NewPacketCapture(c.capDevice.Name)
	if err != nil {
		return err
	}

	c.capture.Start()
	c.Running = true

	go func() {
		for packet := range c.capture.Packets() {
			//Process Packets
			fmt.Printf("Packet: %s \n", packet)
		}
	}()
	return nil
}

func (c *Capturer) EndCapture() {
	c.capture.Stop()
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
