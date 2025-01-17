package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"src/sim/page"
	"src/sim/process"
)

var (
	num_processes      = flag.Uint("num_processes", 128, "number of processes to be generated")
	max_arrive_time    = flag.Uint("max_arrive_time", 256, "maximum time at which a process 'arrives' to be scheduled")
	max_execution_time = flag.Uint("max_execution_time", 16, "maximum execution time for a generated process")
	num_pages          = flag.Uint("num_pages", 64, "amount of pages in virtual memory")
	total_refs         = flag.Uint("total_refs", 512, "amount of virtual memory accesses")
)

func main() {
	flag.Parse()
	switch {
	case *num_processes != 128 && *num_processes > math.MaxUint16:
		log.Panicf("num_processes has to be a 16 bit unsigned integer, only values between %d and %d are allowed", 0, math.MaxUint16)
	case *max_arrive_time != 256 && *max_arrive_time > math.MaxUint16:
		log.Panicf("max_arrive_time has to be a 16 bit unsigned integer, only values between %d and %d are allowed", 0, math.MaxUint16)
	case *max_execution_time != 16 && *max_execution_time > math.MaxUint16:
		log.Panicf("max_execution_time has to be a 16 bit unsigned integer, only values between %d and %d are allowed", 0, math.MaxUint16)
	case *num_pages != 64 && *num_pages > math.MaxUint16:
		log.Panicf("num_pages has to be a 16 bit unsigned integer, only values between %d and %d are allowed", 0, math.MaxUint16)
	case *total_refs != 512 && *total_refs > math.MaxUint16:
		log.Panicf("total_refs has to be a 16 bit unsigned integer, only values between %d and %d are allowed", 0, math.MaxUint16)
	}

	// generating and saving the process simulation input
	processes := process.Gen(uint16(*num_processes), uint16(*max_arrive_time), uint16(*max_execution_time))
	Save(processes, fmt.Sprint("in/",
		*num_processes, "_processes/",
		*max_arrive_time, "_max_arrive_time/",
		*max_execution_time, "_max_execution_time"))

	// running the process simulation
	processSimulationResults := process.Sim(processes,
		process.LCFS,
		process.PreemptiveLCFS,
		process.SJF,
		process.PreemptiveSJF)

	// saving process simulation results
	save_process_results := func(i int, alg string) {
		Save(processSimulationResults[i], fmt.Sprint("out/",
			*num_processes, "_processes/",
			*max_arrive_time, "_max_arrive_time/",
			*max_execution_time, "_max_execution_time/",
			alg))
	}
	save_process_results(0, "LCFS")
	save_process_results(1, "PreemptiveLCFS")
	save_process_results(2, "SJF")
	save_process_results(3, "PreemptiveSJF")

	// generating the page simulation input and the results
	referencePattern, pageSimulationResults := page.Sim(uint16(*num_pages), uint16(*total_refs),
		page.FIFO,
		page.LFU,
		page.PersistentFrequencyLFU)

	// saving page simulation input
	page.SaveReferencePattern(referencePattern, fmt.Sprint("in/",
		*num_pages, "_pages/",
		*total_refs, "_refs"))

	// saving page simulation results
	save_page_results := func(i int, alg string) {
		Save(pageSimulationResults[i], fmt.Sprint("out/",
			*num_pages, "_pages/",
			*total_refs, "_refs/",
			alg))
	}
	save_page_results(0, "FIFO")
	save_page_results(1, "LFU")
	save_page_results(2, "PersistentFrequencyLFU")
}
