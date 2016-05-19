package ts

// Packet represents Transport Stream(TS) packet.
type Packet []byte

// PID represents Packet Identifier, describing the payload data.
type PID uint16

const (
	// SyncWord is used to identify the start of a TS Packet.
	SyncWord = 0x47
)

// TS Header

// SyncWord returns Sync word.
func (p Packet) SyncWord() byte {
	return p[0]
}

// TransportErrorIndicator returns Transport Error Indicator(TEI).
func (p Packet) TransportErrorIndicator() bool {
	i := p[1] & 0x80 >> 7
	return i == 1
}

// PayloadUnitStartIndicator returns Payload Unit Start Indicator.
func (p Packet) PayloadUnitStartIndicator() bool {
	i := p[1] & 0x40 >> 6
	return i == 1
}

// TransportPriority returns Transport Priority.
func (p Packet) TransportPriority() bool {
	i := p[1] & 0x20 >> 5
	return i == 1
}

// PID returns Packet Identifier(PID).
func (p Packet) PID() PID {
	return PID(uint16(p[1]&0x1f)<<8 | uint16(p[2]))
}

// ScramblingControl returns Scrambling control.
func (p Packet) ScramblingControl() uint8 {
	return p[3] & 0xc0 >> 6
}

// AdaptationFieldFlag returns Adaptation field flag.
func (p Packet) AdaptationFieldFlag() bool {
	i := p[3] & 0x20 >> 5
	return i == 1
}

// PayloadFlag returns Payload flag.
func (p Packet) PayloadFlag() bool {
	i := p[3] & 0x10 >> 4
	return i == 1
}

// ContinuityCounter returns Continuity counter.
func (p Packet) ContinuityCounter() uint8 {
	return p[3] & 0x0f
}
