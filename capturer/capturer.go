package capturer

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

type CaptureManager struct {
	capDevice *pcap.Interface
	capture   *PacketCapture
	Running   bool
	SigCh     chan string
	CapCh     chan string
	drawCh    chan string
	track     string
	track1    string
	track2    string
	track3    string
}

func (c *CaptureManager) Init(pDrawCh chan string) {
	c.setDeviceName()
	c.Running = false
	c.SigCh = make(chan string)
	c.CapCh = make(chan string)
	c.drawCh = pDrawCh
	c.track = ""
	c.track1 = ""
	c.track2 = ""
	c.track3 = ""
}

func (c *CaptureManager) StartAutomation() {
	c.Running = true
	c.SigCh = make(chan string)
	go c.actionListener()

	go c.StartCapture()
}

func (c *CaptureManager) actionListener() {
	c.updateTrackers("", "", " - Action Listener Started")
	for {
		select {
		case <-c.SigCh:
			c.updateTrackers("", "", " - Action Listener Stopped")
			return
		case <-c.CapCh:
			go c.StartCapture()
		}
	}
}

func (c *CaptureManager) EndAutomation() {
	c.Running = false
	c.capture.Stop()
	close(c.SigCh)
}

func (c *CaptureManager) StartCapture() error {
	var err error
	c.capture, err = NewPacketCapture(c.capDevice.Name)
	if err != nil {
		return err
	}

	c.Running = true
	c.capture.Start()

	filename := c.getFilename()

	go c.WritePacketToFile(filename, c.capture.packetChan)
	c.updateTrackers((" - Now Capturing: " + filename), "", "")
	return nil
}

func (c *CaptureManager) WritePacketToFile(filename string, packetChan <-chan gopacket.Packet) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := pcapgo.NewWriter(f)
	err = w.WriteFileHeader(128, layers.LinkTypeEthernet)
	if err != nil {
		return err
	}

	for packet := range packetChan {
		err = w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		if err != nil {
			return err
		}
	}

	//Move packet to pcaps/done folder
	f.Close()
	dstPath := strings.Split(filename, "/")[0] + "/done/" + strings.Split(filename, "/")[1]
	c.updateTrackers("", " - "+(strings.Split(filename, "/")[1]+" moved to "+strings.Split(filename, "/")[0]+"/done"), "")
	err = os.Rename(filename, dstPath)
	if err != nil {
		return err
	}

	if c.Running {
		c.CapCh <- ""
	}
	return nil
}
func (c *CaptureManager) EndCapture() {
	c.Running = false
	c.capture.Stop()
}

func (c *CaptureManager) PrintCli() {
	// fmt.Println(" -> You're Within the Capture Menu")
	if c.Running {
		fmt.Println("Capture Running...")
		fmt.Println(c.getFilename())
	} else {
		fmt.Println("")
		fmt.Println("")
	}
	fmt.Println(c.track)
	c.PrintTargetDevice()
}

func (c *CaptureManager) updateTrackers(t1 string, t2 string, t3 string) {
	if t1 != "" {
		c.track1 = t1
	}
	if t2 != "" {
		c.track2 = t2
	}
	if t3 != "" {
		c.track3 = t3
	}
	c.drawCh <- ""
}

func (c *CaptureManager) PrintAutomation() {
	if c.Running {
		fmt.Println(" =====> CAPTURE RUNNING...")
	} else {
		fmt.Println(" =====> not capturing")
	}
	// What captures are running?
	fmt.Println(c.track1)
	fmt.Println(c.track2)
	// What device is capturing?
	fmt.Println(c.track3)
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
			strings.Contains(device.Name, "eth") ||
			strings.Contains(device.Name, "wl") ||
			strings.Contains(device.Description, "Wi-Fi") {
			c.capDevice = &device
			return
		}
	}
}

func (c *CaptureManager) getFilename() string {
	format := "2006-01-02_15-04-05"
	timeNow := time.Now().Format(format)
	return fmt.Sprintf("pcaps/capture_%s.pcap", timeNow)
}
