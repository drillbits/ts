package ts

import (
	"bufio"
	"bytes"
	"io"
)

const packetDefaultSize = 188

// PacketScanner is a wrapper of bufio.Scanner.
type PacketScanner struct {
	*bufio.Scanner
}

// NewPacketScanner returns a new Scanner to read from r.
func NewPacketScanner(r io.Reader) *PacketScanner {
	s := bufio.NewScanner(r)
	s.Split(splitPacket)
	return &PacketScanner{s}
}

func splitPacket(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if len(data) < packetDefaultSize {
		return 0, nil, nil
	}
	i := bytes.IndexByte(data[packetDefaultSize:len(data)], byte(SyncWord))
	if i >= 0 {
		return i + packetDefaultSize, data[0 : i+packetDefaultSize], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

// Packet returns bytes as Packet.
func (s *PacketScanner) Packet() Packet {
	return s.Bytes()
}
