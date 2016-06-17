package ts

// PAT (Program Association Table) provides the correspondence between a program_number and the PID value of the transport stream packets which carry the program definition.
// The program_number is the numeric label associated with a program.
type PAT []byte

type ProgramAssociationSection []byte

func (p PAT) TableID() uint8 {
	return p[0]
}

func (p PAT) SectionLength() uint16 {
	return uint16(p[1]&0x0F)<<8 | uint16(p[2])
}

func (p PAT) ProgramAssociationSections() []ProgramAssociationSection {
	var sections []ProgramAssociationSection
	start := 8
	l := (int(p.SectionLength()) - 9) / 4
	for i := 0; i < l; i++ {
		sections = append(sections, ProgramAssociationSection(p[start:(start+4)]))
		start += 4
	}
	return sections
}

func (s ProgramAssociationSection) ProgramNumber() uint16 {
	return uint16(s[0])<<8 | uint16(s[1])
}

func (s ProgramAssociationSection) PID() PID {
	return PID(uint16(s[2]&0x1F)<<8 | uint16(s[3]))
}
