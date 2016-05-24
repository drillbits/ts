package ts

import (
	"bytes"
	"testing"
)

func TestSyncWord(t *testing.T) {
	// 01000111 00000001 00010001 00110111
	p := Packet{0x47, 0x01, 0x11, 0x37}

	sw := p.SyncWord()
	if sw != SyncWord {
		t.Errorf("got: 0x%x, expected: 0x%x", sw, SyncWord)
	}
}

func TestTransportErrorIndicator(t *testing.T) {
	// 010000111 10000000
	p := Packet{0x47, 0x80}

	tei := p.TransportErrorIndicator()
	if !tei {
		t.Errorf("got: %t, expected: %t", tei, true)
	}

	// 010000111 00000000
	p = Packet{0x47, 0x00}

	tei = p.TransportErrorIndicator()
	if tei {
		t.Errorf("got: %t, expected: %t", tei, false)
	}
}

func TestPayloadUnitStartIndicator(t *testing.T) {
	// 010000111 01000000
	p := Packet{0x47, 0x40}

	tei := p.PayloadUnitStartIndicator()
	if !tei {
		t.Errorf("got: %t, expected: %t", tei, true)
	}

	// 010000111 00000000
	p = Packet{0x47, 0x00}

	tei = p.PayloadUnitStartIndicator()
	if tei {
		t.Errorf("got: %t, expected: %t", tei, false)
	}
}

func TestTransportPriority(t *testing.T) {
	// 010000111 00100000
	p := Packet{0x47, 0x20}

	tei := p.TransportPriority()
	if !tei {
		t.Errorf("got: %t, expected: %t", tei, true)
	}

	// 010000111 00000000
	p = Packet{0x47, 0x00}

	tei = p.TransportPriority()
	if tei {
		t.Errorf("got: %t, expected: %t", tei, false)
	}
}

func TestPID(t *testing.T) {
	// 01000111 00000001 00010001 00110111
	p := Packet{0x47, 0x01, 0x11, 0x37}

	pid := p.PID()
	if pid != PID(0x111) {
		t.Errorf("got: 0x%x, expected: 0x%x", pid, PID(0x111))
	}
}

func TestScramblingControl(t *testing.T) {
	// 01000111 00000001 00010001 00110111
	p := Packet{0x47, 0x01, 0x11, 0x37}

	sc := p.ScramblingControl()
	if sc != 0x00 {
		t.Errorf("got: 0x%x, expected: 0x%x", sc, 0x00)
	}

	// 01000111 00000001 00010001 01110111
	p = Packet{0x47, 0x01, 0x11, 0x77}

	sc = p.ScramblingControl()
	if sc != 0x01 {
		t.Errorf("got: 0x%x, expected: 0x%x", sc, 0x01)
	}

	// 01000111 00000001 00010001 10110111
	p = Packet{0x47, 0x01, 0x11, 0xb7}

	sc = p.ScramblingControl()
	if sc != 0x02 {
		t.Errorf("got: 0x%x, expected: 0x%x", sc, 0x02)
	}

	// 01000111 00000001 00010001 11110111
	p = Packet{0x47, 0x01, 0x11, 0xf7}

	sc = p.ScramblingControl()
	if sc != 0x03 {
		t.Errorf("got: 0x%x, expected: 0x%x", sc, 0x03)
	}
}

func TestAdaptationFieldFlag(t *testing.T) {
	// 01000111 00000001 00010001 00110111
	p := Packet{0x47, 0x01, 0x11, 0x37}

	aff := p.AdaptationFieldFlag()
	if !aff {
		t.Errorf("got: %t, expected: %t", aff, true)
	}

	// 01000111 00000001 00010001 00010111
	p = Packet{0x47, 0x01, 0x11, 0x17}

	aff = p.AdaptationFieldFlag()
	if aff {
		t.Errorf("got: %t, expected: %t", aff, false)
	}
}

func TestPayloadFlag(t *testing.T) {
	// 01000111 00000001 00010001 00110111
	p := Packet{0x47, 0x01, 0x11, 0x37}

	pf := p.PayloadFlag()
	if !pf {
		t.Errorf("got: %t, expected: %t", pf, true)
	}

	// 01000111 00000001 00010001 00100111
	p = Packet{0x47, 0x01, 0x11, 0x27}

	pf = p.PayloadFlag()
	if pf {
		t.Errorf("got: %t, expected: %t", pf, false)
	}
}

func TestContinuityCounter(t *testing.T) {
	// 01000111 00000001 00010001 00110111
	p := Packet{0x47, 0x01, 0x11, 0x37}

	cc := p.ContinuityCounter()
	if cc != 7 {
		t.Errorf("got: %d, expected: %d", cc, 7)
	}
}

func TestAdaptationFieldLength(t *testing.T) {
	// 01000111 00000001 00010001 00110111 00000001 00100000 11100000 01110010
	p := Packet{0x47, 0x01, 0x11, 0x37, 0x01, 0x20, 0xE0, 0x72}

	afLen := p.AdaptationFieldLength()
	if afLen != 1 {
		t.Errorf("got: %d, expected: %d", afLen, 1)
	}
}

func TestAdaptationField(t *testing.T) {
	// 01000111 00000001 00010001 00110111
	p := Packet{0x47, 0x01, 0x11, 0x37,
		// 10110111 00010000 01111010 00110100 00001111 00010100 01111110 01111000 11111111...
		0xB7, 0x10, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
	expected := AdaptationField{
		// 10110111 00010000 01111010 00110100 00001111 00010100 01111110 01111000 11111111...
		0xB7, 0x10, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

	af := p.AdaptationField()
	t.Logf("PID: 0x%X", p.PID())
	if bytes.Compare(af, expected) != 0 {
		t.Errorf("got: 0x%0X, expected: 0x%X", af, expected)
	}
}

func TestAdaptationFieldAdaptationFieldLength(t *testing.T) {
	// 01000111 00000001 00010001 00110111
	p := Packet{0x47, 0x01, 0x11, 0x37,
		// 10110111 00010000 01111010 00110100 00001111 00010100 01111110 01111000 11111111...
		0xB7, 0x10, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

	af := p.AdaptationField()
	afLen := af.Length()
	if afLen != 183 {
		t.Errorf("got: %d, expected: %d", afLen, 183)
	}
}

func TestDiscontinuityIndicator(t *testing.T) {
	// 10110111 01111111 01111010 00110100 00001111 00010100 01111110 01111000 11111111
	af := AdaptationField{0xB7, 0x7F, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF}

	di := af.DiscontinuityIndicator()
	if di {
		t.Errorf("got: %t, expected: %t", di, false)
	}

	// 10110111 10000000 01111010 00110100 00001111 00010100 01111110 01111000 11111111
	af = AdaptationField{0xB7, 0x80, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF}

	di = af.DiscontinuityIndicator()
	if !di {
		t.Errorf("got: %t, expected: %t", di, true)
	}
}

func TestRandomAccessIndicator(t *testing.T) {
	// 10110111 10111111 01111010 00110100 00001111 00010100 01111110 01111000 11111111
	af := AdaptationField{0xB7, 0xBF, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF}

	di := af.RandomAccessIndicator()
	if di {
		t.Errorf("got: %t, expected: %t", di, false)
	}

	// 10110111 01000000 01111010 00110100 00001111 00010100 01111110 01111000 11111111
	af = AdaptationField{0xB7, 0x40, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF}

	di = af.RandomAccessIndicator()
	if !di {
		t.Errorf("got: %t, expected: %t", di, true)
	}
}

func TestElementaryStreamPriorityIndicator(t *testing.T) {
	// 10110111 11011111 01111010 00110100 00001111 00010100 01111110 01111000 11111111
	af := AdaptationField{0xB7, 0xDF, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF}

	di := af.ElementaryStreamPriorityIndicator()
	if di {
		t.Errorf("got: %t, expected: %t", di, false)
	}

	// 10110111 00100000 01111010 00110100 00001111 00010100 01111110 01111000 11111111
	af = AdaptationField{0xB7, 0x20, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF}

	di = af.ElementaryStreamPriorityIndicator()
	if !di {
		t.Errorf("got: %t, expected: %t", di, true)
	}
}

func TestPCRFlag(t *testing.T) {
	// 10110111 11101111 01111010 00110100 00001111 00010100 01111110 01111000 11111111
	af := AdaptationField{0xB7, 0xEF, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF}

	di := af.PCRFlag()
	if di {
		t.Errorf("got: %t, expected: %t", di, false)
	}

	// 10110111 00010000 01111010 00110100 00001111 00010100 01111110 01111000 11111111
	af = AdaptationField{0xB7, 0x10, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF}

	di = af.PCRFlag()
	if !di {
		t.Errorf("got: %t, expected: %t", di, true)
	}
}

func TestOPCRFlag(t *testing.T) {
	// 10110111 11110111 01111010 00110100 00001111 00010100 01111110 01111000 11111111
	af := AdaptationField{0xB7, 0xF7, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF}

	di := af.OPCRFlag()
	if di {
		t.Errorf("got: %t, expected: %t", di, false)
	}

	// 10110111 00001000 01111010 00110100 00001111 00010100 01111110 01111000 11111111
	af = AdaptationField{0xB7, 0x08, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF}

	di = af.OPCRFlag()
	if !di {
		t.Errorf("got: %t, expected: %t", di, true)
	}
}

func TestSplicingPointFlag(t *testing.T) {
	// 10110111 11111011 01111010 00110100 00001111 00010100 01111110 01111000 11111111
	af := AdaptationField{0xB7, 0xFB, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF}

	di := af.SplicingPointFlag()
	if di {
		t.Errorf("got: %t, expected: %t", di, false)
	}

	// 10110111 00000100 01111010 00110100 00001111 00010100 01111110 01111000 11111111
	af = AdaptationField{0xB7, 0x04, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF}

	di = af.SplicingPointFlag()
	if !di {
		t.Errorf("got: %t, expected: %t", di, true)
	}
}

func TestTransportPrivateDataFlag(t *testing.T) {
	// 10110111 11111101 01111010 00110100 00001111 00010100 01111110 01111000 11111111
	af := AdaptationField{0xB7, 0xFD, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF}

	di := af.TransportPrivateDataFlag()
	if di {
		t.Errorf("got: %t, expected: %t", di, false)
	}

	// 10110111 00000010 01111010 00110100 00001111 00010100 01111110 01111000 11111111
	af = AdaptationField{0xB7, 0x02, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF}

	di = af.TransportPrivateDataFlag()
	if !di {
		t.Errorf("got: %t, expected: %t", di, true)
	}
}

func TestAdaptationFieldExtensionFlag(t *testing.T) {
	// 10110111 11111110 01111010 00110100 00001111 00010100 01111110 01111000 11111111
	af := AdaptationField{0xB7, 0xFE, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF}

	di := af.AdaptationFieldExtensionFlag()
	if di {
		t.Errorf("got: %t, expected: %t", di, false)
	}

	// 10110111 00000001 01111010 00110100 00001111 00010100 01111110 01111000 11111111
	af = AdaptationField{0xB7, 0x01, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF}

	di = af.AdaptationFieldExtensionFlag()
	if !di {
		t.Errorf("got: %t, expected: %t", di, true)
	}
}

// TODO
// func TestPCR(t *testing.T) {}
// func TestOPCR(t *testing.T) {}
// func TestSpliceCountdown(t *testing.T) {}
// func TestTransportPrivateDataLength(t *testing.T) {}
// func TestTransportPrivateData(t *testing.T) {}
// func TestAdaptationExtension(t *testing.T) {}
// func TestAdaptationExtensionLength(t *testing.T) {}
// func TestStuffingBytes(t *testing.T) {}

func TestAdaptationFieldSample(t *testing.T) {
	// 10110111 00010000 01111010 00110100 00001111 00010100 01111110 01111000 11111111...
	af := AdaptationField{
		0xB7, 0x10, 0x7A, 0x34, 0x0F, 0x14, 0x7E, 0x78, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

	t.Logf("Length: %d", af.Length())
	t.Logf("DiscontinuityIndicator: %t", af.DiscontinuityIndicator())
	t.Logf("RandomAccessIndicator: %t", af.RandomAccessIndicator())
	t.Logf("ElementaryStreamPriorityIndicator: %t", af.ElementaryStreamPriorityIndicator())
	t.Logf("PCRFlag: %t", af.PCRFlag())
	t.Logf("OPCRFlag: %t", af.OPCRFlag())
	t.Logf("SplicingPointFlag: %t", af.SplicingPointFlag())
	t.Logf("TransportPrivateDataFlag: %t", af.TransportPrivateDataFlag())
	t.Logf("AdaptationFieldExtensionFlag: %t", af.AdaptationFieldExtensionFlag())
	t.Logf("PCR: 0x%X", af.PCR())
	t.Logf("OPCR: 0x%X", af.OPCR())
	t.Logf("SpliceCountdown: 0x%X", af.SpliceCountdown())
	t.Logf("TransportPrivateDataLength: %d", af.TransportPrivateDataLength())
	t.Logf("TransportPrivateData: 0x%X", af.TransportPrivateData())
	t.Logf("AdaptationExtensionLength: %d", af.AdaptationExtensionLength())
	t.Logf("AdaptationExtension: 0x%X", af.AdaptationExtension())
	t.Logf("StuffingBytes: 0x%X", af.StuffingBytes())
}
