package capturer

import (
	"fmt"
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

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
