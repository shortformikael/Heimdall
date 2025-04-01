package analyzer

import (
	"fmt"
	"os"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type AnalyzerManager struct {
	Running         bool
	conversationMap map[string]*Conversation
	drawCh          *chan string
}

func (a *AnalyzerManager) Init(drawCh *chan string) {
	a.conversationMap = make(map[string]*Conversation)
	a.drawCh = drawCh
}

func (a *AnalyzerManager) Start() {

	handle, err := a.getFileHandle()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer handle.Close()

	a.AnalyzePacketSource(handle)
	a.printAnalysis()
}

func (a *AnalyzerManager) getFileHandle() (*pcap.Handle, error) {
	handle, err := pcap.OpenOffline("pcaps/capture.pcap")
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	return handle, nil
}

func (a *AnalyzerManager) getAvailableFiles() []string {
	r := []string{}

	dirPath := "./pcaps/"

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
	}

	for _, entry := range entries {
		r = append(r, entry.Name())
	}

	return r
}

func (a *AnalyzerManager) AnalyzePacketSource(handle *pcap.Handle) {
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

func (a *AnalyzerManager) printAnalysis() {
	//conversations
	for _, packet := range a.conversationMap {
		fmt.Println(packet.String())
	}
}

func (a *AnalyzerManager) PrintCli() {
	fmt.Println("=== Analysis ===")
	fmt.Println(" -> Youre in the Analyzer menu")
	a.PrintAvailableFiles()
}

func (a *AnalyzerManager) PrintAvailableFiles() {
	for _, files := range a.getAvailableFiles() {
		fmt.Println("-", files)
	}
}
