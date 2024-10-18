package cpu

import (
	log "github.com/sirupsen/logrus"
)

type InstructionHandler func(cpu *CPU6502) uint8

var opcodeTable = [256]InstructionHandler{
	// opcode mappings
	0x00: BRK,

	0x18: CLC,
	0xd8: CLD,
	0x58: CLI,
	0xb8: CLV,

	0x38: SEC,
	0xF8: SED,
	0x78: SEI,

	0x69: ADC_immediate,
	0x65: ADC_zeropage,
	0x75: ADC_zeropage_x,
	0x6d: ADC_absolute,
	0x7d: ADC_absolute_x,
	0x79: ADC_absolute_y,
	0x61: ADC_indirect_x,
	0x71: ADC_indirect_y,

	0xA9: LDA_immediate,
	0xA5: LDA_zeropage,
	0xB5: LDA_zeropage_x,
	0xAD: LDA_absolute,
	0xBD: LDA_absolute_x,
	0xB9: LDA_absolute_y,
	0xA1: LDA_indirect_x,
	0xB1: LDA_indirect_y,

	0xA2: LDX_immediate,
	0xA6: LDX_zeropage,
	0xB6: LDX_zeropage_y,
	0xAE: LDX_absolute,
	0xBE: LDX_absolute_y,

	0xA0: LDY_immediate,
	0xA4: LDY_zeropage,
	0xB4: LDY_zeropage_x,
	0xAC: LDY_absolute,
	0xBC: LDY_absolute_x,

	0x85: STA_zeropage,
	0x95: STA_zeropage_x,
	0x8d: STA_absolute,
	0x9d: STA_absolute_x,
	0x99: STA_absolute_y,
	0x81: STA_indirect_x,
	0x91: STA_indirect_y,

	0xEA: NOP,
}

func flagToByte(flag bool) byte {
	b := byte(0)
	if flag {
		b = byte(1)
	}
	return b
}

/*
------------------------------

	Addressing Modes

------------------------------
*/
func (cpu *CPU6502) addrModeImmediate() uint16 {
	cpu.PC++
	return cpu.PC
}

func (cpu *CPU6502) addrModeZeropage() uint16 {
	cpu.PC++
	return combineBytes(cpu.Bus.Data[cpu.PC], 0x00)
}

func (cpu *CPU6502) addrModeZeropageIndexed(register byte) uint16 {
	cpu.PC++
	start := cpu.Bus.Data[cpu.PC] + register
	addr := combineBytes(start, 0x00)
	return addr
}

func (cpu *CPU6502) addrModeAbsolute() uint16 {
	cpu.PC++
	low := cpu.Bus.Data[cpu.PC]
	cpu.PC++
	high := cpu.Bus.Data[cpu.PC]
	return combineBytes(low, high)
}

func (cpu *CPU6502) addrModeAbsoluteIndexed(register byte) uint16 {
	cpu.PC++
	low := cpu.Bus.Data[cpu.PC]
	cpu.PC++
	high := cpu.Bus.Data[cpu.PC]
	start := combineBytes(low, high)
	addr := start + combineBytes(register, 0x00)
	return addr
}

// TODO: Should indirect-mode increment the PC to read the high-order byte?
// Yes, indirect mode specifies 3 bytes - 1 for the opcode, and 2 for the address
func (cpu *CPU6502) addrModeIndirect() uint16 {
	cpu.PC++
	low := cpu.Bus.Data[cpu.PC]
	cpu.PC++
	high := cpu.Bus.Data[cpu.PC]
	return combineBytes(low, high)
}

// index-indirect and indirect-indexed do not increment the PC past
// the initial address
func (cpu *CPU6502) addrModeIndirectX() uint16 {
	cpu.PC++
	val := cpu.Bus.Data[cpu.PC]
	addr := val + cpu.X
	low := cpu.Bus.Data[addr]
	high := cpu.Bus.Data[addr+1]
	return combineBytes(low, high)
}

func (cpu *CPU6502) addrModeIndirectY() uint16 {
	cpu.PC++
	low := cpu.Bus.Data[cpu.PC]
	high := cpu.Bus.Data[cpu.PC+1]
	addr := combineBytes(low, high)
	offset := combineBytes(cpu.Y, 0x00)
	return addr + offset
}

/*------------------------------

	Opcodes

------------------------------*/

func BRK(cpu *CPU6502) uint8 {
	log.Trace("BRK -> entry")
	log.Trace("BRK -> exit")
	return 6
}

/*
ADC - Add memory to accumulator with carry
A + M + C -> A,C
SR: NVZC
*/
func ADC(cpu *CPU6502, addr uint16) {
	// TODO - Revisit me
	m := cpu.Bus.Fetch(addr)
	a := cpu.A
	r := uint16(a) + uint16(m)
	if cpu.SR.C {
		r += 1
	}
	if r == 0 {
		cpu.SR.Z = true
	}
	cpu.A = byte(r & 0xFF)
}

func ADC_immediate(cpu *CPU6502) uint8 {
	ADC(cpu, cpu.addrModeImmediate())
	return 1
}

func ADC_zeropage(cpu *CPU6502) uint8 {
	ADC(cpu, cpu.addrModeZeropage())
	return 2
}

func ADC_zeropage_x(cpu *CPU6502) uint8 {
	ADC(cpu, cpu.addrModeZeropageIndexed(cpu.X))
	return 3
}

func ADC_absolute(cpu *CPU6502) uint8 {
	ADC(cpu, cpu.addrModeAbsolute())
	return 3
}

func ADC_absolute_x(cpu *CPU6502) uint8 {
	ADC(cpu, cpu.addrModeAbsoluteIndexed(cpu.X))
	return 3
}

func ADC_absolute_y(cpu *CPU6502) uint8 {
	ADC(cpu, cpu.addrModeAbsoluteIndexed(cpu.Y))
	return 3
}

func ADC_indirect_x(cpu *CPU6502) uint8 {
	ADC(cpu, cpu.addrModeIndirectX())
	return 5
}

func ADC_indirect_y(cpu *CPU6502) uint8 {
	ADC(cpu, cpu.addrModeIndirectY())
	return 4
}

func CLC(cpu *CPU6502) uint8 {
	cpu.SR.C = false
	return 1
}

func CLD(cpu *CPU6502) uint8 {
	cpu.SR.D = false
	return 1
}

func CLI(cpu *CPU6502) uint8 {
	cpu.SR.I = false
	return 1
}

func CLV(cpu *CPU6502) uint8 {
	cpu.SR.V = false
	return 1
}

func SEC(cpu *CPU6502) uint8 {
	cpu.SR.C = true
	return 1
}

func SED(cpu *CPU6502) uint8 {
	cpu.SR.D = true
	return 1
}

func SEI(cpu *CPU6502) uint8 {
	cpu.SR.I = true
	return 1
}

func NOP(cpu *CPU6502) uint8 {
	log.Trace("NOP -> entry")
	log.Trace("NOP -> exit")
	return 1
}

/*
LDA: M->A
cZidbvN
*/
func checkZandN(cpu *CPU6502, value byte) byte {
	cpu.SR.Z = value == 0
	cpu.SR.N = (value & 0x80) != 0
	return value
}

func LDA_immediate(cpu *CPU6502) uint8 {
	log.Trace("LDA_immediate -> entry")
	m := cpu.Bus.Fetch(cpu.addrModeImmediate())
	cpu.A = checkZandN(cpu, m)
	log.Trace("LDA_immediate -> exit")
	return 1
}

func LDA_zeropage(cpu *CPU6502) uint8 {
	log.Trace("LDA_zeropage -> entry")
	m := cpu.Bus.Fetch(cpu.addrModeZeropage())
	cpu.A = checkZandN(cpu, m)
	log.Trace("LDA_zeropage -> exit")
	return 2
}

func LDA_zeropage_x(cpu *CPU6502) uint8 {
	log.Trace("LDA_zeropage_x -> entry")
	m := cpu.Bus.Fetch(cpu.addrModeZeropageIndexed(cpu.X))
	cpu.A = checkZandN(cpu, m)
	log.Trace("LDA_zeropage_x -> exit")
	return 3
}

func LDA_absolute(cpu *CPU6502) uint8 {
	log.Trace("LDA_absolute -> entry")
	m := cpu.Bus.Fetch(cpu.addrModeAbsolute())
	cpu.A = checkZandN(cpu, m)
	log.Trace("LDA_absolute -> exit")
	return 3
}

func LDA_absolute_x(cpu *CPU6502) uint8 {
	log.Trace("LDA_absolute_x -> entry")
	m := cpu.Bus.Fetch(cpu.addrModeAbsoluteIndexed(cpu.X))
	cpu.A = checkZandN(cpu, m)
	log.Trace("LDA_absolute_x -> exit")
	return 3
}

func LDA_absolute_y(cpu *CPU6502) uint8 {
	log.Trace("LDA_absolute_y -> entry")
	m := cpu.Bus.Fetch(cpu.addrModeAbsoluteIndexed(cpu.Y))
	cpu.A = checkZandN(cpu, m)
	log.Trace("LDA_absolute_y -> exit")
	return 3
}

func LDA_indirect_x(cpu *CPU6502) uint8 {
	log.Trace("LDA_indirect_x -> entry")
	m := cpu.Bus.Fetch(cpu.addrModeIndirectX())
	cpu.A = checkZandN(cpu, m)
	log.Trace("LDA_indirect_x -> entry")
	return 5
}

func LDA_indirect_y(cpu *CPU6502) uint8 {
	log.Trace("LDA_indirect_y -> entry")
	m := cpu.Bus.Fetch(cpu.addrModeIndirectY())
	cpu.A = checkZandN(cpu, m)
	log.Trace("LDA_indirect_y -> entry")
	return 4
}

/*
LDX: M->X
cZidbvN
*/
func LDX_immediate(cpu *CPU6502) uint8 {
	log.Trace("LDX_immediate -> entry")
	m := cpu.Bus.Fetch(cpu.addrModeImmediate())
	cpu.X = checkZandN(cpu, m)
	log.Trace("LDX_immediate -> exit")
	return 1
}

func LDX_zeropage(cpu *CPU6502) uint8 {
	log.Trace("LDX_zeropage -> entry")
	m := cpu.Bus.Fetch(cpu.addrModeZeropage())
	cpu.X = checkZandN(cpu, m)
	log.Trace("LDX_zeropage -> exit")
	return 2
}

func LDX_zeropage_y(cpu *CPU6502) uint8 {
	log.Trace("LDX_zeropage_y -> entry")
	m := cpu.Bus.Fetch(cpu.addrModeZeropageIndexed(cpu.Y))
	cpu.X = checkZandN(cpu, m)
	log.Trace("LDX_zeropage_y -> exit")
	return 3
}

func LDX_absolute(cpu *CPU6502) uint8 {
	log.Trace("LDX_absolute -> entry")
	m := cpu.Bus.Fetch(cpu.addrModeAbsolute())
	cpu.X = checkZandN(cpu, m)
	log.Trace("LDX_absolute -> exit")
	return 3
}

func LDX_absolute_y(cpu *CPU6502) uint8 {
	log.Trace("LDX_absolute_y -> entry")
	m := cpu.Bus.Fetch(cpu.addrModeAbsoluteIndexed(cpu.Y))
	cpu.X = checkZandN(cpu, m)
	log.Trace("LDX_absolute_y -> exit")
	return 3
}

/*
LDY: M->Y
cZidbvN
*/
func LDY_immediate(cpu *CPU6502) uint8 {
	log.Trace("LDY_immediate -> entry")
	m := cpu.Bus.Fetch(cpu.addrModeImmediate())
	cpu.Y = checkZandN(cpu, m)
	log.Trace("LDY_immediate -> exit")
	return 1
}

func LDY_zeropage(cpu *CPU6502) uint8 {
	log.Trace("LDY_zeropage -> entry")
	m := cpu.Bus.Fetch(cpu.addrModeZeropage())
	cpu.Y = checkZandN(cpu, m)
	log.Trace("LDY_zeropage -> exit")
	return 2
}

func LDY_zeropage_x(cpu *CPU6502) uint8 {
	log.Trace("LDY_zeropage_x -> entry")
	m := cpu.Bus.Fetch(cpu.addrModeZeropageIndexed(cpu.X))
	cpu.Y = checkZandN(cpu, m)
	log.Trace("LDY_zeropage_y -> exit")
	return 3
}

func LDY_absolute(cpu *CPU6502) uint8 {
	log.Trace("LDY_absolute -> entry")
	m := cpu.Bus.Fetch(cpu.addrModeAbsolute())
	cpu.Y = checkZandN(cpu, m)
	log.Trace("LDY_absolute -> exit")
	return 3
}

func LDY_absolute_x(cpu *CPU6502) uint8 {
	log.Trace("LDX_absolute_x -> entry")
	m := cpu.Bus.Fetch(cpu.addrModeAbsoluteIndexed(cpu.X))
	cpu.Y = checkZandN(cpu, m)
	log.Trace("LDX_absolute_x -> exit")
	return 3
}

func STA_zeropage(cpu *CPU6502) uint8 {
	cpu.Bus.Store(cpu.addrModeZeropage(), cpu.A)
	return 2
}

func STA_zeropage_x(cpu *CPU6502) uint8 {
	cpu.Bus.Store(cpu.addrModeZeropageIndexed(cpu.X), cpu.A)
	return 3
}

func STA_absolute(cpu *CPU6502) uint8 {
	cpu.Bus.Store(cpu.addrModeAbsolute(), cpu.A)
	return 3
}

func STA_absolute_x(cpu *CPU6502) uint8 {
	cpu.Bus.Store(cpu.addrModeAbsoluteIndexed(cpu.X), cpu.A)
	return 4
}

func STA_absolute_y(cpu *CPU6502) uint8 {
	cpu.Bus.Store(cpu.addrModeAbsoluteIndexed(cpu.Y), cpu.A)
	return 4
}

func STA_indirect_x(cpu *CPU6502) uint8 {
	cpu.Bus.Store(cpu.addrModeIndirectX(), cpu.A)
	return 5
}

func STA_indirect_y(cpu *CPU6502) uint8 {
	cpu.Bus.Store(cpu.addrModeIndirectY(), cpu.A)
	return 5
}

func STX_zeropage(cpu *CPU6502) uint8 {
	cpu.Bus.Store(cpu.addrModeZeropage(), cpu.X)
	return 2
}

func STX_zeropage_y(cpu *CPU6502) uint8 {
	cpu.Bus.Store(cpu.addrModeZeropageIndexed(cpu.Y), cpu.X)
	return 3
}

func STX_absolute(cpu *CPU6502) uint8 {
	cpu.Bus.Store(cpu.addrModeAbsolute(), cpu.X)
	return 3
}

func STY_zeropage(cpu *CPU6502) uint8 {
	cpu.Bus.Store(cpu.addrModeZeropage(), cpu.Y)
	return 2
}

func STY_zeropage_x(cpu *CPU6502) uint8 {
	cpu.Bus.Store(cpu.addrModeZeropageIndexed(cpu.X), cpu.Y)
	return 3
}

func STY_absolute(cpu *CPU6502) uint8 {
	cpu.Bus.Store(cpu.addrModeAbsolute(), cpu.Y)
	return 3
}
