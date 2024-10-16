package cpu

import (
	log "github.com/sirupsen/logrus"
)

type DataBus struct {
	Data [1024 * 64]byte
}

func (bus *DataBus) Fetch(addr uint16) byte {
	return bus.Data[addr]
}

func (bus *DataBus) Store(addr uint16, val byte) {
	bus.Data[addr] = val
}

type CPU struct {
	// Control signals
	RDY  bool
	IRQ  bool
	NMI  bool
	SYNC bool
	RW   bool
	BE   bool
	RES  bool

	// Registers
	A byte // accumulator
	X byte // index register
	Y byte // index register

	// Stack Pointer, Program Counter
	SP uint16 // could be a byte since it only uses the low side (01ff-0100)
	PC uint16

	// Status Register
	// 7 - - - - - - 0
	// N V   B D I Z C
	SR byte

	Bus DataBus
}

type StatusRegister int

const (
	C StatusRegister = iota
	Z
	I
	D
	B
	_
	V
	N
)

func combineBytes(low byte, high byte) uint16 {
	return uint16(high)<<8 | uint16(low)
}

func extractBit(val byte, n int) bool {
	mask := byte(1) << n
	return (val & mask) != 0
}

func setBit(val byte, bit int, status bool) {
	// TODO - implement me
}

func (cpu *CPU) Reset() {
	cpu.RDY = true
	cpu.RES = false

	cpu.A = 0x00
	cpu.X = 0x00
	cpu.Y = 0x00

	cpu.SP = 0x01ff

	// load the reset vector into the program counter
	// FFFC (low)
	// FFFD (high)
	cpu.PC = combineBytes(cpu.Bus.Data[0xFFFC], cpu.Bus.Data[0xFFFD])

	cpu.SR = 0x00
}

func (cpu *CPU) PHI2() {
	log.Trace("PHI2 -> entry")

	// Fetch instruction
	log.Trace("PHI2 -> getting next opcode")
	opcode := cpu.Bus.Data[cpu.PC]

	// Not sure where the PC should be incremented yet
	//cpu.PC++

	// execute the InstructionHandler
	opcodeTable[opcode](cpu)
	log.Trace("PHI2 -> exit")
}
