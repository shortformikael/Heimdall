package analyzer

import (
	"fmt"
	"strconv"
)

type Conversation struct {
	Src         string `json:"src"`
	Dst         string `json:"dst"`
	Count       int    `json:"count"`
	Size        int    `json:"size"` //Bytes
	Protocol    string `json:"protocol"`
	Application string `json:"application"`
}

func (c *Conversation) Append(packet *analyzerPacket) {
	c.Count++
	c.Size += packet.Length
}

func (c *Conversation) GetKey() string { return c.Src + " -> " + c.Dst }

func (c *Conversation) String() string {
	size := strconv.FormatFloat((float64(c.Size) / 1000), 'f', 2, 64)
	return fmt.Sprintf(
		"Key: %v | Protocol %v | Count: %v | Size %s",
		c.GetKey(),
		c.Protocol,
		c.Count,
		size+"KB",
	)
}
