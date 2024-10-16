package cpu

import (
	log "github.com/sirupsen/logrus"
)

type InstructionHandler func(cpu *CPU)

var opcodeTable = [256]InstructionHandler{
	// opcode mappings
	0x00: BRK,

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

	0xEA: NOP,
}

/*
------------------------------

	Addressing Modes

------------------------------
*/
func (cpu *CPU) addrModeImmediate() uint16 {
	cpu.PC++
	return cpu.PC
}

func (cpu *CPU) addrModeZeropage() uint16 {
	cpu.PC++
	return combineBytes(cpu.Bus.Data[cpu.PC], 0x00)
}

func (cpu *CPU) addrModeZeropageIndexed(register byte) uint16 {
	cpu.PC++
	start := combineBytes(cpu.Bus.Data[cpu.PC], 0x00)
	addr := start + combineBytes(register, 0x00)
	return addr
}

func (cpu *CPU) addrModeAbsolute() uint16 {
	cpu.PC++
	low := cpu.Bus.Data[cpu.PC]
	cpu.PC++
	high := cpu.Bus.Data[cpu.PC]
	return combineBytes(low, high)
}

func (cpu *CPU) addrModeAbsoluteIndexed(register byte) uint16 {
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
func (cpu *CPU) addrModeIndirect() uint16 {
	cpu.PC++
	low := cpu.Bus.Data[cpu.PC]
	cpu.PC++
	high := cpu.Bus.Data[cpu.PC]
	return combineBytes(low, high)
}

// index-indirect and indirect-indexed do not increment the PC past
// the initial address
func (cpu *CPU) addrModeIndirectX() uint16 {
	cpu.PC++
	val := cpu.Bus.Data[cpu.PC]
	addr := val + cpu.X
	low := cpu.Bus.Data[addr]
	high := cpu.Bus.Data[addr+1]
	return combineBytes(low, high)
}

func (cpu *CPU) addrModeIndirectY() uint16 {
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

func BRK(cpu *CPU) {
	log.Trace("BRK -> entry")
	log.Trace("BRK -> exit")
}

func NOP(cpu *CPU) {
	log.Trace("NOP -> entry")
	log.Trace("NOP -> exit")
}

func LDA_immediate(cpu *CPU) {
	log.Trace("LDA_immediate -> entry")
	cpu.A = cpu.Bus.Fetch(cpu.addrModeImmediate())
	log.Trace("LDA_immediate -> exit")
}

func LDA_zeropage(cpu *CPU) {
	log.Trace("LDA_zeropage -> entry")
	cpu.A = cpu.Bus.Fetch(cpu.addrModeZeropage())
	log.Trace("LDA_zeropage -> exit")
}

func LDA_zeropage_x(cpu *CPU) {
	log.Trace("LDA_zeropage_x -> entry")
	cpu.A = cpu.Bus.Fetch(cpu.addrModeZeropageIndexed(cpu.X))
	log.Trace("LDA_zeropage_x -> exit")
}

func LDA_absolute(cpu *CPU) {
	log.Trace("LDA_absolute -> entry")
	cpu.A = cpu.Bus.Fetch(cpu.addrModeAbsolute())
	log.Trace("LDA_absolute -> exit")
}

func LDA_absolute_x(cpu *CPU) {
	log.Trace("LDA_absolute_x -> entry")
	cpu.A = cpu.Bus.Fetch(cpu.addrModeAbsoluteIndexed(cpu.X))
	log.Trace("LDA_absolute_x -> exit")
}

func LDA_absolute_y(cpu *CPU) {
	log.Trace("LDA_absolute_y -> entry")
	cpu.A = cpu.Bus.Fetch(cpu.addrModeAbsoluteIndexed(cpu.Y))
	log.Trace("LDA_absolute_y -> exit")
}

func LDA_indirect_x(cpu *CPU) {
	log.Trace("LDA_indirect_x -> entry")
	cpu.A = cpu.Bus.Fetch(cpu.addrModeIndirectX())
	log.Trace("LDA_indirect_x -> entry")
}

func LDA_indirect_y(cpu *CPU) {
	log.Trace("LDA_indirect_y -> entry")
	cpu.A = cpu.Bus.Fetch(cpu.addrModeIndirectY())
	log.Trace("LDA_indirect_y -> entry")
}

func LDX_immediate(cpu *CPU) {
	log.Trace("LDX_immediate -> entry")
	cpu.X = cpu.Bus.Fetch(cpu.addrModeImmediate())
	log.Trace("LDX_immediate -> exit")
}

func LDX_zeropage(cpu *CPU) {
	log.Trace("LDX_zeropage -> entry")
	cpu.X = cpu.Bus.Fetch(cpu.addrModeZeropage())
	log.Trace("LDX_zeropage -> exit")
}

func LDX_zeropage_y(cpu *CPU) {
	log.Trace("LDX_zeropage_y -> entry")
	cpu.X = cpu.Bus.Fetch(cpu.addrModeZeropageIndexed(cpu.Y))
	log.Trace("LDX_zeropage_y -> exit")
}

func LDX_absolute(cpu *CPU) {
	log.Trace("LDX_absolute -> entry")
	cpu.X = cpu.Bus.Fetch(cpu.addrModeAbsolute())
	log.Trace("LDX_absolute -> exit")
}

func LDX_absolute_y(cpu *CPU) {
	log.Trace("LDX_absolute_y -> entry")
	cpu.X = cpu.Bus.Fetch(cpu.addrModeAbsoluteIndexed(cpu.Y))
	log.Trace("LDX_absolute_y -> exit")
}

func LDY_immediate(cpu *CPU) {
	log.Trace("LDY_immediate -> entry")
	cpu.Y = cpu.Bus.Fetch(cpu.addrModeImmediate())
	log.Trace("LDY_immediate -> exit")
}

func LDY_zeropage(cpu *CPU) {
	log.Trace("LDY_zeropage -> entry")
	cpu.Y = cpu.Bus.Fetch(cpu.addrModeZeropage())
	log.Trace("LDY_zeropage -> exit")
}

func LDY_zeropage_x(cpu *CPU) {
	log.Trace("LDY_zeropage_x -> entry")
	cpu.Y = cpu.Bus.Fetch(cpu.addrModeZeropageIndexed(cpu.X))
	log.Trace("LDY_zeropage_y -> exit")
}

func LDY_absolute(cpu *CPU) {
	log.Trace("LDY_absolute -> entry")
	cpu.Y = cpu.Bus.Fetch(cpu.addrModeAbsolute())
	log.Trace("LDY_absolute -> exit")
}

func LDY_absolute_x(cpu *CPU) {
	log.Trace("LDX_absolute_x -> entry")
	cpu.Y = cpu.Bus.Fetch(cpu.addrModeAbsoluteIndexed(cpu.X))
	log.Trace("LDX_absolute_x -> exit")
}
