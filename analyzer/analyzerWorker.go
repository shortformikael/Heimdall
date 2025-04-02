package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type analyzerWorker struct {
	conversationMap map[string]*Conversation
	wg              *sync.WaitGroup
	filename        string
	donePath        string
}

func NewWorker(pWG *sync.WaitGroup, pFilename string, pDone string) *analyzerWorker {
	return &analyzerWorker{
		conversationMap: make(map[string]*Conversation),
		wg:              pWG,
		filename:        pFilename,
		donePath:        pDone,
	}
}

func (a *analyzerWorker) Start() {
	defer a.wg.Done()
	//fmt.Println("Worker for file", a.filename, "started...")
	handle, err := a.getFileHandle()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer handle.Close()

	a.AnalyzePacketSource(handle)
	// Serialize hashmap
	err = a.saveToJSON()
	if err != nil {
		fmt.Println("error writing to file:", err)
		return
	}
	handle.Close()
	split := strings.Split(a.filename, "/")
	dstPath := a.donePath + "/" + split[len(split)-1]
	err = os.Rename(a.filename, dstPath)
	if err != nil {
		fmt.Println("Error:", err)
	}
	//fmt.Println("Worker for file", a.filename, "Ended!")
}

func (a *analyzerWorker) saveToJSON() error {
	split := strings.Split(a.filename, "/")
	dstPath := "./entries/" + strings.Split(split[len(split)-1], ".")[0] + ".json"
	//fmt.Println(dstPath)
	file, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	var entries []Conversation
	for _, value := range a.conversationMap {
		entries = append(entries, *value)
	}

	encoder.SetIndent("", "  ")
	if err = encoder.Encode(entries); err != nil {
		return err
	}
	return nil
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
				Src:         packet.Src,
				Dst:         packet.Dst,
				Count:       1,
				Size:        packet.Length,
				Protocol:    packet.Protocol,
				Application: packet.Application,
			}
		}
	}
}

func (a *analyzerWorker) getFileHandle() (*pcap.Handle, error) {
	handle, err := pcap.OpenOffline(a.filename)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	return handle, nil
}
