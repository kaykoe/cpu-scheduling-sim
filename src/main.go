package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"src/sim/page"
	"src/sim/process"
	"strings"
)

var (
	sim_processes      = flag.Bool("sim-processes", false, "run the process simulation")
	sim_pages          = flag.Bool("sim-pages", false, "run the page simulation")
	num_processes      = flag.Uint("num-processes", 128, "number of processes to be generated")
	max_arrive_time    = flag.Uint("max-arrive-time", 256, "maximum time at which a process 'arrives' to be scheduled")
	max_execution_time = flag.Uint("max-execution-time", 16, "maximum execution time for a generated process")
	num_pages          = flag.Uint("num-pages", 64, "amount of pages in virtual memory")
	total_refs         = flag.Uint("total-refs", 512, "amount of virtual memory accesses")
)

func main() {
	flag.Parse()
	switch {
	case *num_processes != 128 && *num_processes > math.MaxUint16:
		log.Panicf("num-processes has to be a 16 bit unsigned integer, only values between %d and %d are allowed", 0, math.MaxUint16)
	case *max_arrive_time != 256 && *max_arrive_time > math.MaxUint16:
		log.Panicf("max-arrive-time has to be a 16 bit unsigned integer, only values between %d and %d are allowed", 0, math.MaxUint16)
	case *max_execution_time != 16 && *max_execution_time > math.MaxUint16:
		log.Panicf("max-execution-time has to be a 16 bit unsigned integer, only values between %d and %d are allowed", 0, math.MaxUint16)
	case *num_pages != 64 && *num_pages > math.MaxUint16:
		log.Panicf("num-pages has to be a 16 bit unsigned integer, only values between %d and %d are allowed", 0, math.MaxUint16)
	case *total_refs != 512 && *total_refs > math.MaxUint16:
		log.Panicf("total-refs has to be a 16 bit unsigned integer, only values between %d and %d are allowed", 0, math.MaxUint16)
	case !*sim_processes && !*sim_pages:
		log.Panic("you must specify at least one simulation to run, run the program with -h for help")
	}

	if *sim_processes {
		log.Printf("Running process simulation with the following parameters:"+
			"\nnum-processes: %d"+
			"\nmax-arrive-time: %d"+
			"\nmax-execution-time: %d\n\n",
			*num_processes, *max_arrive_time, *max_execution_time)

		log.Println("Generating process simulation input...")
		processes := process.Gen(uint16(*num_processes), uint16(*max_arrive_time), uint16(*max_execution_time))
		log.Print("Processes generated successfully\n\n")

		processInputDirectory := fmt.Sprint("in/",
			*num_processes, "-processes/",
			*max_arrive_time, "-max-arrive-time/",
			*max_execution_time, "-max-execution-time")

		log.Println("Saving process simulation input...")
		Save(processes, processInputDirectory)
		log.Print("Process simulation input saved to : ../", processInputDirectory, "\n\n")

		log.Println("Running process simulation...")
		processSimulationResults := process.Sim(processes,
			process.LCFS,
			process.PreemptiveLCFS,
			process.SJF,
			process.PreemptiveSJF)
		log.Print("Process simulation completed successfully\n\n")

		log.Println("Saving process simulation results...")
		save_process_results := func(i int, alg string) {
			Save(processSimulationResults[i], fmt.Sprint("out/",
				*num_processes, "-processes/",
				*max_arrive_time, "-max-arrive-time/",
				*max_execution_time, "-max-execution-time/",
				alg))
		}
		save_process_results(0, "LCFS")
		save_process_results(1, "PreemptiveLCFS")
		save_process_results(2, "SJF")
		save_process_results(3, "PreemptiveSJF")
		log.Print("Process simulation results saved to : ../out/",
			*num_processes, "-processes/",
			*max_arrive_time, "-max-arrive-time/",
			*max_execution_time, "-max-execution-time/",
			"{algorithmName}\n\n",
			strings.Repeat("-", 80),
			"\n\n")
	}

	if *sim_pages {
		log.Printf("Running page simulation with the following parameters:"+
			"\nnum-pages: %d"+
			"\ntotal-refs: %d\n\n",
			*num_pages, *total_refs)

		log.Println("Generating page simulation input and running simulation...")
		referencePattern, pageSimulationResults := page.Sim(uint16(*num_pages), uint16(*total_refs),
			page.FIFO,
			page.LFU,
			page.PersistentFrequencyLFU)
		log.Print("Page simulation completed successfully\n\n")

		log.Println("Saving page simulation input...")
		page.SaveReferencePattern(referencePattern, fmt.Sprint("in/",
			*num_pages, "-pages/",
			*total_refs, "-refs"))
		log.Print("Page simulation input saved to : ../in/",
			*num_pages, "-pages/",
			*total_refs, "-refs\n\n")

		log.Println("Saving page simulation results...")
		save_page_results := func(i int, alg string) {
			Save(pageSimulationResults[i], fmt.Sprint("out/",
				*num_pages, "-pages/",
				*total_refs, "-refs/",
				alg))
		}
		save_page_results(0, "FIFO")
		save_page_results(1, "LFU")
		save_page_results(2, "PersistentFrequencyLFU")
		log.Print("Page simulation results saved to : ../out/",
			*num_pages, "-pages/",
			*total_refs, "-refs/",
			"{algorithmName}\n\n",
			strings.Repeat("-", 80),
			"\n\n")
	}
}
