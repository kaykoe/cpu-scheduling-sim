package main

import (
	"flag"
	"log"
	"math"
	"src/sim/process"
)

var num = flag.Uint("num", 1000, "number of processes to be generated")
var maxarrivetime = flag.Uint("maxarrivetime", 2048, "maximum time at which a process 'arrives' to be scheduled")
var maxexecutiontime = flag.Uint("maxexecutiontime", 32, "maximum execution time for a generated process")

func main() {
	flag.Parse()
	if *num != 1000 {
		if *num > math.MaxUint16 {
			log.Panicf("num has to be a 16 bit unsigned integer, only values between %d and %d are allowed", 0, math.MaxUint16)
		}
	}
	if *maxarrivetime != 2048 {
		if *maxarrivetime > math.MaxUint16 {
			log.Panicf("maxarrivetime has to be a 16 bit unsigned integer, only values between %d and %d are allowed", 0, math.MaxUint16)
		}
	}
	if *maxexecutiontime != 32 {
		if *maxexecutiontime > math.MaxUint16 {
			log.Panicf("maxexecutiontime has to be a 16 bit unsigned integer, only values between %d and %d are allowed", 0, math.MaxUint16)
		}
	}

	processes := process.Gen(uint16(*num), uint16(*maxarrivetime), uint16(*maxexecutiontime))
	Save(processes, "in")

	simulationResults := process.Sim(processes, []process.Alg{process.LCFS,
		process.PreemptiveLCFS,
		process.SJF,
		process.PreemptiveSJF})

	Save(simulationResults[0], "out/LCFS")
	Save(simulationResults[1], "out/PreemptiveLCFS")
	Save(simulationResults[2], "out/SJF")
	Save(simulationResults[3], "out/PreemptiveSJF")
}
