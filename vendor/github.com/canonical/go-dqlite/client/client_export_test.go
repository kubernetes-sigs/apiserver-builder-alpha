package client

import (
	"github.com/canonical/go-dqlite/internal/protocol"
)

func (c *Client) Protocol() *protocol.Protocol {
	return c.protocol
}
