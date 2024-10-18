package cpu

import (
	"testing"
)

const (
	X_init  = 0x00
	Y_init  = 0x00
	A_init  = 0x00
	SP_init = 0x01FF
	PC_init = 0x0000
)

func Test_CanResetCPU(t *testing.T) {

	dut := NewCPU()
	dut.X = 0xDE
	dut.Y = 0xAD
	dut.A = 0xBE
	dut.SP = 0xEFDE
	dut.PC = 0xADBE

	dut.Reset()
	if dut.X != X_init {
		t.Errorf("X: got %2.2X, want %2.2X", dut.X, X_init)
	}
	if dut.Y != Y_init {
		t.Errorf("Y: got %2.2X, want %2.2X", dut.Y, Y_init)
	}
	if dut.A != A_init {
		t.Errorf("A: got %2.2X, want %2.2X", dut.A, A_init)
	}
	if dut.SP != SP_init {
		t.Errorf("SP: got %4.4X, want %4.4X", dut.SP, SP_init)
	}
	if dut.PC != PC_init {
		t.Errorf("PC: got %4.4X, want %4.4X", dut.PC, PC_init)
	}

}

func Test_CombineBytes(t *testing.T) {

	var tests = []struct {
		low  byte
		high byte
		want uint16
	}{
		{0xEF, 0xBE, 0xBEEF},
		{0xAB, 0x00, 0x00AB},
		{0x00, 0x02, 0x0200},
	}

	for _, tt := range tests {
		t.Run("bytes", func(t *testing.T) {
			got := combineBytes(tt.low, tt.high)
			if got != tt.want {
				t.Errorf("got 0x%4.4X, want 0x%4.4X", got, tt.want)
			}
		})
	}
}
