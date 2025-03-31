package analyzer

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type AnalyzerManager struct {
	Running     bool
	PacketArray []*analyzerPacket
}

func (a *AnalyzerManager) Start() {
	handle, err := pcap.OpenOffline("pcaps/capture.pcap")
	if err != nil {
		fmt.Println("Error:", err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		a.PacketArray = append(a.PacketArray, newPacket(packet))
	}
	a.PrintPackets()
	a.startAnalysis()
}

func (a *AnalyzerManager) startAnalysis() {
	//conversations
}

func (a *AnalyzerManager) PrintPackets() {
	for _, packet := range a.PacketArray {
		fmt.Println("")
		packet.Print()
	}
}

func (a *AnalyzerManager) PrintCli() {
	fmt.Println(" -> Youre in the Analyzer menu")

}
