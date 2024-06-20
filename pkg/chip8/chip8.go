package chip8

import "math/rand"

//44:20
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
	// clear screen
	CLS uint16 = 0x00E0
	// return
	RET uint16 = 0x00EE
	// jump to routine
	SYS uint16 = 0x0
	// jump to location
	JP uint16 = 0x1
	// Call subroutine
	CALL uint16 = 0x2
	// skip next instruction if Vx = kk
	SE_VX_BYTE uint16 = 0x3
	// skip next instruction if Vx != kk
	SNE uint16 = 0x4
	// skip next instruction if Vx = V
	SE_VX_VY uint16 = 0x5
	//Set Vx = kk
	LD_VX_BYTE = 0x6
	//Set Vx = Vx + kk
	ADD_VX_BYTE = 0x7
	// skip next instruction if Vx != Vy
	SNE_VX_VY uint16 = 0x9
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
	c.opcode = c.getOpcode()

	switch instruction := c.opcode >> 12; instruction {
	case SYS:
		switch op := c.opcode; op {
		case CLS:
			c.clearScreen()
		case RET:
			c.popStack()
		}
		c.incrementPC()
	case JP:
		var addr = c.opcode & AddressBitMask
		c.programCounter = addr
	case CALL:
		c.callSubroutine()
	case SE_VX_BYTE:
		c.skipVxEqualByte()
	case SNE:
		c.skipVxNotEqualByte()
	case SE_VX_VY:
		c.skipVxEqualVy()
	case LD_VX_BYTE:
		c.loadVxBytes()
	case ADD_VX_BYTE:
		c.addVxBytes()
	case 0x8:
		var x = (c.opcode & 0x0F00) >> 8
		var y = (c.opcode & 0x00F0) >> 4
		switch mode := c.opcode & 0x00F; mode {
		case 0:
			c.registers[x] = c.registers[y]
		case 1:
			c.registers[x] |= c.registers[y]
		case 2:
			c.registers[x] &= c.registers[y]
		case 3:
			c.registers[x] ^= c.registers[y]
		case 4:
			var sum uint16 = uint16(c.registers[x])
			sum += uint16(c.registers[y])
			if sum > 255 {
				c.registers[0xF] = 1
			} else {
				c.registers[0xF] = 0
			}
			c.registers[x] = uint8(sum & 0x00FF)
		case 5:
			if c.registers[x] > c.registers[y] {
				c.registers[0xF] = 1
			} else {

				c.registers[0xF] = 0
			}
			c.registers[x] -= c.registers[y]
		case 6:
			if c.registers[x]&1 != 0 {
				c.registers[0xF] = 1
			} else {
				c.registers[0xF] = 0
			}
			c.registers[x] = c.registers[x] >> 1
		case 7:
			if c.registers[y] > c.registers[x] {
				c.registers[0xF] = 1
			} else {

				c.registers[0xF] = 0
			}
			c.registers[x] = c.registers[y] - c.registers[x]
		case 14:
			if c.registers[x]&0x80 != 0 {
				c.registers[0xF] = 1
			} else {
				c.registers[0xF] = 0
			}
			c.registers[x] = c.registers[x] >> 1

		}
		c.incrementPC()

	case SNE_VX_VY:
		c.skipVxNotEqualVy()
	case 0xA:
		c.index = c.opcode & 0x0FFF
		c.incrementPC()
	case 0xB:
		c.programCounter = (c.opcode & 0x0FFF) + uint16(c.registers[0])
	case 0xC:
		var x = (c.opcode & 0x0F00) >> 8
		var kk = c.opcode & 0x00FF
		c.registers[x] = uint8(uint16(rand.Intn(100)) & kk)
		c.incrementPC()
	case 0xD:
		c.registers[0xF] = 0
		var xx = (c.opcode & 0x0F00) >> 8
		var yy = (c.opcode & 0x00F0) >> 4
		var nn = c.opcode & 0x000F

		var regX = c.registers[xx]
		var regY = c.registers[yy]

		var y uint8 = 0
		for uint16(y) < nn {
			var pixel = c.memory[c.index+uint16(y)]
			var x uint8 = 0
			for x < 8 {
				const msb uint8 = 0x080
				if pixel&(msb>>x) != 0 {
					var tX = (regX + x) % 64
					var tY = (regY + y) % 32
					var idx = tX + tY*64
					c.graphics[idx] ^= 1
					if c.graphics[idx] == 0 {
						c.registers[0xF] = 1

					}
				}

				x += 1
			}
			y += 1
		}
		c.incrementPC()
	case 0xE:
		var x = (c.opcode & 0x0F00) >> 8
		var kk = c.opcode & 0x00FF

		if kk == 0x9E {
			if c.keys[c.registers[x]] == 1 {
				c.incrementPC()
			}
		} else if kk == 0xA1 {
			if c.keys[c.registers[x]] != 1 {
				c.incrementPC()
			}
		}
		c.incrementPC()
	case 0xF:
		var x = (c.opcode & 0x0F00) >> 8
		var kk = c.opcode & 0x00FF
		if kk == 0x07 {
			c.registers[x] = c.delayTimer
		} else if kk == 0x0A {
			var key_pressed = false
			for i := 0; i < len(c.keys); i++ {
				if c.keys[i] != 0 {
					c.registers[x] = uint8(i)
					key_pressed = true
					break
				}
			}
			if !key_pressed {
				return
			}
		} else if kk == 0x15 {
			c.delayTimer = c.registers[x]
		} else if kk == 0x18 {
			c.delayTimer = c.registers[x]
		} else if kk == 0x1E {
			c.index += uint16(c.registers[x])
		} else if kk == 0x29 {
			c.index = uint16(c.registers[x]) * 0x5
		} else if kk == 0x33 {
			c.memory[c.index] = c.registers[x] / 100
			c.memory[c.index+1] = (c.registers[x] / 10) % 10
			c.memory[c.index+2] = c.registers[x] % 10
		} else if kk == 0x55 {
			var i uint16 = 0
			for ; i <= x; i++ {
				c.memory[c.index+i] = c.registers[i]
			}
		} else if kk == 0x65 {
			var i uint16 = 0
			for ; i <= x; i++ {
				c.registers[i] = c.memory[c.index+i]
			}
		}

		c.incrementPC()
	}
}

func (c *Chip8) getOpcode() uint16 {
	return uint16(c.memory[c.programCounter])<<8 | uint16(c.memory[c.programCounter+1])
}
func (c *Chip8) clearScreen() {
	for idx := range c.graphics {
		c.graphics[idx] = 0
	}
}
func (c *Chip8) popStack() {
	c.sp -= 1
	c.programCounter = c.stack[c.sp]
}
func (c *Chip8) callSubroutine() {
	c.stack[c.sp] = c.programCounter
	c.sp += 1
	c.programCounter = c.opcode & AddressBitMask
}
func (c *Chip8) skipVxEqualByte() {
	var x = (c.opcode & 0x0F00) >> 8
	var r = uint16(c.registers[x])
	if r == (c.opcode & 0x00FF) {
		c.incrementPC()
	}
	c.incrementPC()
}
func (c *Chip8) skipVxNotEqualByte() {
	var x = (c.opcode & 0x0F00) >> 8
	var r = uint16(c.registers[x])
	if r != (c.opcode & 0x00FF) {
		c.incrementPC()
	}
	c.incrementPC()
}
func (c *Chip8) skipVxEqualVy() {
	var x = (c.opcode & 0x0F00) >> 8
	var y = (c.opcode & 0x00F0) >> 4
	var rX = uint16(c.registers[x])
	var rY = uint16(c.registers[y])
	if rX == rY {
		c.incrementPC()
	}
	c.incrementPC()
}
func (c *Chip8) skipVxNotEqualVy() {
	var x = (c.opcode & 0x0F00) >> 8
	var y = (c.opcode & 0x00F0) >> 4
	var rX = uint16(c.registers[x])
	var rY = uint16(c.registers[y])
	if rX != rY {
		c.incrementPC()
	}
	c.incrementPC()
}
func (c *Chip8) loadVxBytes() {
	var x = (c.opcode & 0x0F00) >> 8
	c.registers[x] = uint8(c.opcode & 0x00FF)
	c.incrementPC()
}
func (c *Chip8) addVxBytes() {
	var x = (c.opcode & 0x0F00) >> 8
	c.registers[x] += uint8(c.opcode & 0x00FF)
	c.incrementPC()
}
