package main

import (
	"github.com/pete-dot-m/go65/cpu"
	log "github.com/sirupsen/logrus"
)

const (
	SIXTYFOUR_K = 64 * 1024
)

func init() {
	log.SetLevel(log.TraceLevel)
}

func main() {
	log.Trace("Started app")

	cpu := cpu.NewCPU()
	for i := 0; i < SIXTYFOUR_K; i++ {
		cpu.Bus.Data[i] = 0x00
	}

	log.Trace("Calling Reset")
	cpu.Reset()
	// run it
	for range 10 {
		log.Trace("Calling PHI2")
		cpu.PHI2()
	}
}
