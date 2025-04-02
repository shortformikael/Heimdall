package analyzer

import (
	"fmt"
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type analyzerWorker struct {
	conversationMap map[string]*Conversation
	wg              *sync.WaitGroup
	filename        string
}

func NewWorker(pWG *sync.WaitGroup, pFilename string) *analyzerWorker {
	return &analyzerWorker{
		conversationMap: make(map[string]*Conversation),
		wg:              pWG,
		filename:        pFilename,
	}
}

func (a *analyzerWorker) Start() {
	defer a.wg.Done()
	fmt.Println("Worker for file", a.filename, "started...")
	handle, err := a.getFileHandle()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer handle.Close()

	a.AnalyzePacketSource(handle)
	// Serialize hashmap
	fmt.Println("Worker for file", a.filename, "Ended!")
}

func (a *analyzerWorker) AnalyzePacketSource(handle *pcap.Handle) {
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for rawPacket := range packetSource.Packets() {
		packet := newPacket(rawPacket)
		key := fmt.Sprintf("%s -> %s | %s", packet.Src, packet.Dst, packet.Protocol)
		conv, exists := a.conversationMap[key]
		if exists {
			conv.Append(packet)
		} else {
			a.conversationMap[key] = &Conversation{
				src:         packet.Src,
				dst:         packet.Dst,
				count:       1,
				size:        packet.Length,
				protocol:    packet.Protocol,
				application: packet.Application,
			}
		}
	}
}

func (a *analyzerWorker) getFileHandle() (*pcap.Handle, error) {
	handle, err := pcap.OpenOffline("pcaps/" + a.filename)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	return handle, nil
}
