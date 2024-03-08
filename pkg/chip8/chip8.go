package chip8

var fontset = [...]uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

const (
	ClearScreen uint16 = 0x00E0
	Return      uint16 = 0x00EE
)

const AddressBitMask uint16 = 0x0FFF

type Chip8 struct {
	//0xXXXX
	opcode         uint16
	memory         [4096]uint8
	graphics       [64 * 32]uint8
	registers      [16]uint8
	index          uint16
	programCounter uint16

	delayTimer uint8
	soundTimer uint8

	stack [16]uint16
	sp    uint16

	keys [16]uint8
}

func NewChip8() *Chip8 {
	machine := new(Chip8)
	machine.programCounter = 0x200

	for idx, char := range fontset {
		machine.memory[idx] = char
	}

	return machine
}

func (c *Chip8) incrementPC() {
	// every instruction is two bytes but
	// can only read one byte at time via memory
	// hence we increment pc by 2
	c.programCounter += 2
}

func (c *Chip8) Cycle() {
	c.opcode = uint16(c.memory[c.programCounter])<<8 | uint16(c.memory[c.programCounter+1])

	//0xX000 most significat byte
	switch first := c.opcode >> 12; first {
	case 0x0:
		switch op := c.opcode; op {
		case ClearScreen:
			for idx := range c.graphics {
				c.graphics[idx] = 0

			}

		case Return:
			
		}
		c.incrementPC()
	case 0x1:
		c.programCounter = c.opcode & AddressBitMask
	case 0x2:
		c.stack[c.sp] = c.programCounter
		c.sp += 1 
		c.programCounter = c.opcode & AddressBitMask
	}
}


//TimeStamp 25:11
