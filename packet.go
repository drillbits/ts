package ts

// Packet represents Transport Stream(TS) packet.
type Packet []byte

// PID represents Packet Identifier, describing the payload data.
type PID uint16

// AdaptationField represents Extended TS header.
type AdaptationField []byte

const (
	// SyncByte is used to identify the start of a TS Packet.
	SyncByte = 0x47
)

// TS Header

// SyncByte returns Sync byte.
func (p Packet) SyncByte() byte {
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

// TransportScramblingControl returns Scrambling control.
func (p Packet) TransportScramblingControl() uint8 {
	return p[3] & 0xc0 >> 6
}

// AdaptationFieldControl returns Adaptation field control.
func (p Packet) AdaptationFieldControl() byte {
	return p[3] & 0x30 >> 4
}

// AdaptationFieldFlag returns Adaptation field flag.
func (p Packet) AdaptationFieldFlag() bool {
	ctrl := p.AdaptationFieldControl()
	return ctrl == 0x02 || ctrl == 0x03
}

// PayloadFlag returns Payload flag.
func (p Packet) PayloadFlag() bool {
	ctrl := p.AdaptationFieldControl()
	return ctrl == 0x01 || ctrl == 0x03
}

// ContinuityCounter returns Continuity counter.
func (p Packet) ContinuityCounter() uint8 {
	return p[3] & 0x0F
}

// AdaptationFieldLength returns number of bytes in the adaptation field immediately following this byte.
func (p Packet) AdaptationFieldLength() int {
	if !p.AdaptationFieldFlag() {
		return 0
	}
	return int(p[4])
}

// AdaptationField returns Adaptation field.
func (p Packet) AdaptationField() AdaptationField {
	afLen := p.AdaptationFieldLength()
	if afLen == 0 {
		return nil
	}
	return AdaptationField(p[4 : 5+afLen])
}

// PayloadData returns Payload data bytes.
func (p Packet) PayloadData() []byte {
	if !p.PayloadFlag() {
		return nil
	}
	start := 4
	if p.AdaptationFieldFlag() {
		start += p.AdaptationFieldLength()
	}
	return p[start:len(p)]
}

// Length returns number of bytes in the adaptation field immediately following this byte.
func (af AdaptationField) Length() int {
	return int(af[0])
}

// DiscontinuityIndicator returns Discontinuity indicator.
func (af AdaptationField) DiscontinuityIndicator() bool {
	i := af[1] & 0x80 >> 7
	return i == 1
}

// RandomAccessIndicator returns Random Access indicator.
func (af AdaptationField) RandomAccessIndicator() bool {
	i := af[1] & 0x40 >> 6
	return i == 1
}

// ElementaryStreamPriorityIndicator returns Elementary stream priority indicator.
func (af AdaptationField) ElementaryStreamPriorityIndicator() bool {
	i := af[1] & 0x20 >> 5
	return i == 1
}

// PCRFlag returns PCR flag.
func (af AdaptationField) PCRFlag() bool {
	i := af[1] & 0x10 >> 4
	return i == 1
}

// OPCRFlag returns OPCR flag.
func (af AdaptationField) OPCRFlag() bool {
	i := af[1] & 0x08 >> 3
	return i == 1
}

// SplicingPointFlag returns Splicing point flag.
func (af AdaptationField) SplicingPointFlag() bool {
	i := af[1] & 0x04 >> 2
	return i == 1
}

// TransportPrivateDataFlag returns Transport private data flag.
func (af AdaptationField) TransportPrivateDataFlag() bool {
	i := af[1] & 0x02 >> 1
	return i == 1
}

// AdaptationFieldExtensionFlag returns Adaptation field extension flag.
func (af AdaptationField) AdaptationFieldExtensionFlag() bool {
	i := af[1] & 0x01
	return i == 1
}

// PCR returns Program clock reference.
func (af AdaptationField) PCR() []byte {
	if !af.PCRFlag() {
		return nil
	}
	return af[2:8]
}

// OPCR returns Original program clock reference.
func (af AdaptationField) OPCR() []byte {
	if !af.OPCRFlag() {
		return nil
	}
	start := 2
	if af.PCRFlag() {
		start += 6
	}
	return af[start : start+6]
}

// SpliceCountdown indicates how many TS packets from this one a splicing point occurs.
func (af AdaptationField) SpliceCountdown() byte {
	if !af.SplicingPointFlag() {
		return 0
	}
	start := 2
	if af.PCRFlag() {
		start += 6
	}
	if af.OPCRFlag() {
		start += 6
	}
	return af[start]
}

// TransportPrivateDataLength returns number of bytes in the transport private data immediately following this byte.
func (af AdaptationField) TransportPrivateDataLength() int {
	if !af.TransportPrivateDataFlag() {
		return 0
	}
	start := 2
	if af.PCRFlag() {
		start += 6
	}
	if af.OPCRFlag() {
		start += 6
	}
	if af.SplicingPointFlag() {
		start++
	}
	return int(af[start])
}

// TransportPrivateData returns private data.
func (af AdaptationField) TransportPrivateData() []byte {
	if !af.TransportPrivateDataFlag() {
		return nil
	}
	start := 2
	if af.PCRFlag() {
		start += 6
	}
	if af.OPCRFlag() {
		start += 6
	}
	if af.SplicingPointFlag() {
		start++
	}
	// Transport private data length
	start++
	return af[start : start+af.TransportPrivateDataLength()]
}

// AdaptationExtension returns Adaptation field extension.
func (af AdaptationField) AdaptationExtension() []byte {
	if !af.AdaptationFieldExtensionFlag() {
		return nil
	}
	start := 2
	if af.PCRFlag() {
		start += 6
	}
	if af.OPCRFlag() {
		start += 6
	}
	if af.SplicingPointFlag() {
		start++
	}
	if af.TransportPrivateDataFlag() {
		start++ // Transport private data length
		start += af.TransportPrivateDataLength()
	}
	return af[start : start+af.AdaptationExtensionLength()]
}

// AdaptationExtensionLength returns number of bytes in the adaptation extension field immediately following this byte.
func (af AdaptationField) AdaptationExtensionLength() int {
	if !af.AdaptationFieldExtensionFlag() {
		return 0
	}
	start := 2
	if af.PCRFlag() {
		start += 6
	}
	if af.OPCRFlag() {
		start += 6
	}
	if af.SplicingPointFlag() {
		start++
	}
	if af.TransportPrivateDataFlag() {
		start++ // Transport private data length
		start += af.TransportPrivateDataLength()
	}
	return int(af[start])
}

// StuffingBytes returns array of 0xFF.
func (af AdaptationField) StuffingBytes() []byte {
	start := 2
	if af.PCRFlag() {
		start += 6
	}
	if af.OPCRFlag() {
		start += 6
	}
	if af.SplicingPointFlag() {
		start++
	}
	if af.TransportPrivateDataFlag() {
		start++ // Transport private data length
		start += af.TransportPrivateDataLength()
	}
	if af.AdaptationFieldExtensionFlag() {
		start++ // Adaptation extension length
		start += af.AdaptationExtensionLength()
	}
	return af[start:len(af)]
}
