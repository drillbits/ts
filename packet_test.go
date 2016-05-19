package ts

import (
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
