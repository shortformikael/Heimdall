package analyzer

import (
	"fmt"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type analyzerPacket struct {
	Src         string
	Dst         string
	Application string
	Protocol    string
	Timestamp   time.Time
	Length      int // Bytes
}

func newPacket(packet gopacket.Packet) *analyzerPacket {
	r := &analyzerPacket{}

	transportLayer := packet.TransportLayer()
	networkLayer := packet.NetworkLayer()

	if networkLayer != nil && transportLayer != nil {
		r.Src = networkLayer.NetworkFlow().Src().String()
		r.Dst = networkLayer.NetworkFlow().Dst().String()
		r.Protocol = getProtocolName(packet)
		r.Application = getApplicationProtocol(packet)
	} else {
		r.Src = "Unkown"
		r.Dst = "Unkown"
		r.Application = "Unkown"
		r.Protocol = "Unkown"
	}

	r.Timestamp = packet.Metadata().Timestamp
	r.Length = packet.Metadata().CaptureLength

	return r
}

func getProtocolName(packet gopacket.Packet) string {
	// Check for different protocol layers
	// HTTP MISSING
	switch {
	case packet.Layer(layers.LayerTypeTCP) != nil:
		return "TCP"
	case packet.Layer(layers.LayerTypeUDP) != nil:
		return "UDP"
	case packet.Layer(layers.LayerTypeICMPv4) != nil:
		return "ICMPv4"
	case packet.Layer(layers.LayerTypeICMPv6) != nil:
		return "ICMPv6"
	case packet.Layer(layers.LayerTypeTLS) != nil:
		return "TLS"
	default:
		return "Unknown"
	}
}

func getApplicationProtocol(packet gopacket.Packet) string {
	switch {
	case packet.Layer(layers.LayerTypeDNS) != nil:
		return "DNS"
	case packet.Layer(layers.LayerTypeTLS) != nil:
		return "TLS/SSL"
	case packet.Layer(layers.LayerTypeDHCPv4) != nil:
		return "DHCP"
	default:
		// Try to guess from ports if no specific layer
		if tcp := packet.Layer(layers.LayerTypeTCP); tcp != nil {
			tcp, _ := tcp.(*layers.TCP)
			switch {
			case tcp.SrcPort == 80 || tcp.DstPort == 80:
				return "HTTP (port 80)"
			case tcp.SrcPort == 443 || tcp.DstPort == 443:
				return "HTTPS (port 443)"
			case tcp.SrcPort == 22 || tcp.DstPort == 22:
				return "SSH (port 22)"
			case tcp.SrcPort == 21 || tcp.DstPort == 21:
				return "FTP (port 21)"
			}
		}
		return "Unknown Application"
	}
}

func (ap *analyzerPacket) Print() {
	fmt.Printf("  * IP: %v -> %v\n", ap.Src, ap.Dst)
	fmt.Printf("  * Protocol: %v | %v\n", ap.Protocol, ap.Application)
	fmt.Printf("  * Length: %d \n", ap.Length)
	fmt.Printf("  * Timestamp: %v \n", ap.Timestamp.String())
}
