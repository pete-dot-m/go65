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

	theCpu := cpu.CPU{}
	for i := 0; i < SIXTYFOUR_K; i++ {
		theCpu.ROM[i] = 0xEA
	}

	log.Trace("Calling Reset")
	theCpu.Reset()
	// run it
	for range 10 {
		log.Trace("Calling PHI2")
		theCpu.PHI2()
	}
}
