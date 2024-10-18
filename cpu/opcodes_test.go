package cpu

import (
	"testing"
)

// test the various addrMode functions to
// ensure the PC is incremented by the correct
// amount, and optionally that the correct address
// is returned
func Test_CanCallAddrModes(t *testing.T) {
	cpu := NewCPU()

	var tests = []struct {
		title  string
		f      func() uint16
		offset uint16
	}{
		{"addrModeImmediate", cpu.addrModeImmediate, 1},
		{"addrModeAbsolute", cpu.addrModeAbsolute, 2},
		{"addrModeZeropage", cpu.addrModeZeropage, 1},
		{"addrModeIndirect", cpu.addrModeIndirect, 2},
		{"addrModeIndirectX", cpu.addrModeIndirectX, 1},
		{"addrModeIndirectY", cpu.addrModeIndirectY, 1},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			pcStart := cpu.PC
			_ = tt.f()
			got := cpu.PC - pcStart
			want := tt.offset
			if got != want {
				t.Errorf("%s: got 0x%4.4X, want 0x%4.4X", tt.title, got, want)
			}
		})
	}
}

func Test_CanCallAddrModesWithRegister(t *testing.T) {
	cpu := NewCPU()

	var tests = []struct {
		title   string
		f       func(reg byte) uint16
		offset  uint16
		address uint16
	}{
		{"addrModeAbsoluteIndexed", cpu.addrModeAbsoluteIndexed, 2, 0},
		{"addrModeZeropageIndexed", cpu.addrModeZeropageIndexed, 1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			pcStart := cpu.PC
			_ = tt.f(cpu.X)
			got := cpu.PC - pcStart
			want := tt.offset
			if got != want {
				t.Errorf("%s: got 0x%4.4X, want 0x%4.4X", tt.title, got, want)
			}
		})
	}
}

func Test_AddrModeImmediateReturnsIncrementedPC(t *testing.T) {
	cpu := NewCPU()
	cpu.PC = 0x0d24
	addr := cpu.addrModeImmediate()
	if addr != cpu.PC {
		t.Errorf("addrModeImmediate returned %4.4X, wanted %4.4X", addr, cpu.PC+1)
	}
}

func Test_AddrModeZeropageReturnsCorrectAddress(t *testing.T) {
	cpu := NewCPU()

	cpu.PC = 0x000F
	cpu.Bus.Data[0x0010] = 0xA0

	got := cpu.addrModeZeropage()
	want := uint16(0x00A0)

	if got != want {
		t.Errorf("got %4.4X, want %4.4X", got, want)
	}
}

func Test_AddrModeZeropageIndexedReturnsCorrectAddress(t *testing.T) {
	cpu := NewCPU()

	cpu.PC = 0x000F
	cpu.X = 0x10
	cpu.Bus.Data[0x0010] = 0xA0

	got := cpu.addrModeZeropageIndexed(cpu.X)
	want := uint16(0xA0 + 0x10)

	if got != want {
		t.Errorf("got %4.4X, want %4.4X", got, want)
	}

	cpu.PC = uint16(0x000F)
	cpu.X = byte(0x01)
	cpu.Bus.Data[0x0010] = byte(0xFF)

	got = cpu.addrModeZeropageIndexed(cpu.X)
	want = uint16(0x0000)
	if got != want {
		t.Errorf("got 0x%4.4X, want 0x%4.4X", got, want)
	}
}

func Test_AddrModeAbsoluteReturnsCorrectAddress(t *testing.T) {
	cpu := NewCPU()

	cpu.PC = 0x001F
	cpu.Bus.Data[0x0020] = 0xEF
	cpu.Bus.Data[0x0021] = 0xBE

	got := cpu.addrModeAbsolute()
	want := uint16(0xBEEF)

	if got != want {
		t.Errorf("got %4.4X, want %4.4X", got, want)
	}
}

func Test_AddrModeAbsoluteIndexedReturnsCorrectAddress(t *testing.T) {
	cpu := NewCPU()

	cpu.PC = 0x001F
	cpu.X = byte(0x01)
	cpu.Bus.Data[0x0020] = 0xEF
	cpu.Bus.Data[0x0021] = 0xBE

	got := cpu.addrModeAbsoluteIndexed(cpu.X)
	want := uint16(0xBEF0)

	if got != want {
		t.Errorf("got %4.4X, want %4.4X", got, want)
	}
}

func Test_AddrModeIndirectReturnsCorrectAddress(t *testing.T) {
	cpu := NewCPU()

	cpu.PC = 0x001F
	cpu.Bus.Data[0x0020] = 0xEF
	cpu.Bus.Data[0x0021] = 0xBE

	got := cpu.addrModeAbsolute()
	want := uint16(0xBEEF)

	if got != want {
		t.Errorf("got %4.4X, want %4.4X", got, want)
	}
}

func Test_AddrModeIndirectXReturnsCorrectAddress(t *testing.T) {
	cpu := NewCPU()

	cpu.PC = uint16(0x001F)
	cpu.X = byte(0x01)
	cpu.Bus.Data[0x0020] = byte(0xEF)
	cpu.Bus.Data[0x00F0] = byte(0xAD)
	cpu.Bus.Data[0x00F1] = byte(0xDE)

	got := cpu.addrModeIndirectX()
	want := uint16(0xDEAD)

	if got != want {
		t.Errorf("got %4.4X, want %4.4X", got, want)
	}
}

func Test_AddrModeIndirectYReturnsCorrectAddress(t *testing.T) {
	cpu := NewCPU()

	cpu.PC = uint16(0x001F)
	cpu.Y = byte(0x01)
	cpu.Bus.Data[0x0020] = byte(0xAD)
	cpu.Bus.Data[0x0021] = byte(0xDE)

	got := cpu.addrModeIndirectY()
	want := uint16(0xDEAE)

	if got != want {
		t.Errorf("got %4.4X, want %4.4X", got, want)
	}
}

func Test_LDAImmediateLoadsAndSetsFlags(t *testing.T) {
	var tests = []struct {
		title string
		m     byte
		wantZ bool
		wantN bool
	}{
		{"LDA 0", 0, true, false},
		{"LDA 0xAB", 0xAB, false, true},
		{"LDA 0x10", 0x10, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			cpu := NewCPU()
			cpu.PC = uint16(0x001F)
			cpu.Bus.Data[0x0020] = tt.m

			LDA_immediate(cpu)

			got := cpu.A
			gotZ := cpu.SR.Z
			gotN := cpu.SR.N

			if got != tt.m {
				t.Errorf("cpu.A: got 0x%2.2X, want 0x%2.2X", got, tt.m)
			}
			if gotZ != tt.wantZ {
				t.Errorf("cpu.SR.Z: got %v, want %v", gotZ, tt.wantZ)
			}
			if gotN != tt.wantN {
				t.Errorf("cpu.SR.N: got %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func Test_LDAZeropageLoadsAndSetsFlags(t *testing.T) {
	var tests = []struct {
		title string
		m     byte
		wantZ bool
		wantN bool
	}{
		{"LDA 0", 0, true, false},
		{"LDA 0xAB", 0xAB, false, true},
		{"LDA 0x10", 0x10, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			cpu := NewCPU()
			cpu.PC = uint16(0x0001)
			cpu.Bus.Data[0x0002] = 0x11
			cpu.Bus.Data[0x0011] = tt.m

			LDA_zeropage(cpu)

			got := cpu.A
			gotZ := cpu.SR.Z
			gotN := cpu.SR.N

			if got != tt.m {
				t.Errorf("cpu.A: got 0x%2.2X, want 0x%2.2X", got, tt.m)
			}
			if gotZ != tt.wantZ {
				t.Errorf("cpu.SR.Z: got %v, want %v", gotZ, tt.wantZ)
			}
			if gotN != tt.wantN {
				t.Errorf("cpu.SR.N: got %v, want %v", gotN, tt.wantN)
			}
		})
	}
}
