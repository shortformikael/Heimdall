package analyzer

import (
	"fmt"
	"strconv"
)

type Conversation struct {
	src         string
	dst         string
	count       int
	size        int //Bytes
	protocol    string
	application string
}

func (c *Conversation) Append(packet *analyzerPacket) {
	c.count++
	c.size += packet.Length
}

func (c *Conversation) GetKey() string { return c.src + " -> " + c.dst }

func (c *Conversation) String() string {
	size := strconv.FormatFloat((float64(c.size) / 1000), 'f', 2, 64)
	return fmt.Sprintf(
		"Key: %v | Protocol %v | Count: %v | Size %s",
		c.GetKey(),
		c.protocol,
		c.count,
		size+"KB",
	)
}
