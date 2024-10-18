package cpu

import (
	log "github.com/sirupsen/logrus"
)

type DataBus struct {
	Stack [256]byte
	Data  [1024 * 64]byte
}

func (bus *DataBus) Fetch(addr uint16) byte {
	return bus.Data[addr]
}

func (bus *DataBus) Store(addr uint16, val byte) {
	bus.Data[addr] = val
}

type StatusRegister struct {
	C bool
	Z bool
	I bool
	D bool
	B bool
	V bool
	N bool
}

type CPU6502 struct {
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
	SR StatusRegister

	Bus DataBus

	// internal clock cycle counter
	cycles uint8
}

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

func NewCPU() *CPU6502 {
	cpu := CPU6502{
		RDY:  true,
		IRQ:  false,
		NMI:  false,
		SYNC: false,
		RW:   true,
		BE:   true,
		RES:  false,

		A: 0x00,
		X: 0x00,
		Y: 0x00,

		SP: 0x01FF,
		PC: 0x00,

		SR: StatusRegister{},

		Bus: DataBus{},

		cycles: 0,
	}

	// TODO: eventually, need to sort out loading ROM, etc.
	for i, _ := range cpu.Bus.Data {
		cpu.Bus.Data[i] = 0x00
	}

	return &cpu
}

func (cpu *CPU6502) Reset() {
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

	cpu.SR = StatusRegister{}

	cpu.cycles = 0
}

func (cpu *CPU6502) PHI2() {
	log.Trace("PHI2 -> entry")

	if cpu.cycles > 0 {
		cpu.cycles--
		return
	}
	// Fetch instruction
	log.Trace("PHI2 -> getting next opcode")
	opcode := cpu.Bus.Data[cpu.PC]

	// execute the InstructionHandler, set cycles to the number of
	// cycles needed for the opcode
	cpu.cycles = opcodeTable[opcode](cpu)

	cpu.PC++

	log.Trace("PHI2 -> exit")
}
