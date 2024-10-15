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
	0xEA: NOP,
}

// opcodes
func BRK(cpu *CPU) {
	log.Trace("BRK -> entry")
	log.Trace("BRK -> exit")
}

func NOP(cpu *CPU) {
	log.Trace("NOP -> entry")
	log.Trace("NOP -> exit")
}

func loadRegisterImmediate(target *byte, cpu *CPU) {
	// load the next byte into target
	cpu.PC += 1
	*target = cpu.ROM[cpu.PC]
}

func loadRegisterZeropage(target *byte, cpu *CPU) {
	// load the value at the zero-page location
	// specified by the next byte (low)
	cpu.PC += 1
	location := combineBytes(cpu.ROM[cpu.PC], 0x00)
	*target = cpu.ROM[location]
}

func loadRegisterZeropageX(target *byte, cpu *CPU) {
	// load the value at the zero-page location
	// specified by the next byte (low)
	cpu.PC += 1
	location := combineBytes(cpu.ROM[cpu.PC], 0x00)
	*target = cpu.ROM[location]
}

func loadRegisterAbsolute(target *byte, cpu *CPU) {
	// load the value at the 16-bit location
	// specified by the next 2 bytes (low, high)
	cpu.PC += 1
	low := cpu.ROM[cpu.PC]
	cpu.PC += 1
	high := cpu.ROM[cpu.PC]
	location := combineBytes(low, high)
	*target = cpu.ROM[location]
}

func loadRegisterAbsoluteX(target *byte, cpu *CPU) {
	// load the value at the 16-bit location
	// specified by the next 2 bytes (low, high)
	cpu.PC += 1
	low := cpu.ROM[cpu.PC]
	cpu.PC += 1
	high := cpu.ROM[cpu.PC]
	location := combineBytes(low, high)
	*target = cpu.ROM[location]
}

func LDA_immediate(cpu *CPU) {
	log.Trace("LDA_immediate -> entry")
	loadRegisterImmediate(&cpu.A, cpu)
	log.Trace("LDA_immediate -> exit")
}

func LDA_zeropage(cpu *CPU) {
	log.Trace("LDA_zeropage -> entry")
	loadRegisterZeropage(&cpu.A, cpu)
	log.Trace("LDA_zeropage -> exit")
}

func LDA_zeropage_x(cpu *CPU) {
	log.Trace("LDA_zeropage_x -> entry")
	loadRegisterZeropageX(&cpu.A, cpu)
	log.Trace("LDA_zeropage_x -> exit")
}

func LDA_absolute(cpu *CPU) {
	log.Trace("LDA_absolute -> entry")
	loadRegisterAbsolute(&cpu.A, cpu)
	log.Trace("LDA_absolute -> exit")
}

func LDA_absolute_x(cpu *CPU) {
	log.Trace("LDA_absolute_x -> entry")
	loadRegisterAbsoluteX(&cpu.A, cpu)
	log.Trace("LDA_absolute_x -> exit")
}
