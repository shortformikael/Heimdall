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
	maxBytes   int
	totalBytes int
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
				pc.totalBytes += len(packet.Data())
				if pc.totalBytes > pc.maxBytes {
					//Limit Reached
					return
				}
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
	// 96 or 128 for snaplen
	handle, err := pcap.OpenLive(device, 96, true, pcap.BlockForever)
	if err != nil {
		return nil, fmt.Errorf("error opening device: %v", err)
	}

	//Set Filter
	err = handle.SetBPFFilter("tcp or udp or icmp")
	if err != nil {
		return nil, fmt.Errorf("error applying filter: %v", err)
	}

	return &PacketCapture{
		handle:     handle,
		stopCh:     make(chan struct{}),
		packetChan: make(chan gopacket.Packet, 1000),
		maxBytes:   10 * (1000 * 1024), // 10MB
		totalBytes: 0,
	}, nil
}
