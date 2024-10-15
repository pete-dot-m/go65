package cpu

import (
	log "github.com/sirupsen/logrus"
)

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

	// buses
	AddressBus uint16
	DataBus    byte

	// ROM
	ROM [65536]byte
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
	cpu.PC = combineBytes(cpu.ROM[0xFFFC], cpu.ROM[0xFFFD])

	cpu.SR = 0x00

	cpu.AddressBus = 0x0000
	cpu.DataBus = 0x00
}

func (cpu *CPU) PHI2() {
	log.Trace("PHI2 -> entry")

	// Fetch instruction
	log.Trace("PHI2 -> getting next opcode")
	opcode := cpu.ROM[cpu.PC]

	// Not sure where the PC should be incremented yet
	//cpu.PC++

	// execute the InstructionHandler
	opcodeTable[opcode](cpu)
	log.Trace("PHI2 -> exit")
}
