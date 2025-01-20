package process

import (
	"cmp"
	"container/heap"
	"log"
	"slices"
	"src/sim"
)

type Alg func(processes *Slice) *Slice

// Sim runs a simulation of the given processes using the strategies in the alg slice
func Sim(processes *Slice, algs ...Alg) (res []*Slice) {
	res = make([]*Slice, len(algs))

	for i, alg := range algs {
		res[i] = alg(processes.Copy())
	}
	return res
}

func PreemptiveLCFS(processes *Slice) *Slice {
	if isSorted := slices.IsSortedFunc([]Process(*processes), func(a, b Process) int {
		return cmp.Compare(a.arriveTime, b.arriveTime)
	}); !isSorted {
		log.Panic("The process scheduling algorithms have to receive a slice of Processes sorted by arriveTime")
	}
	if *processes == nil {
		log.Panic("The process slice to be simulated cannot be nil")
	}
	if len(*processes) == 0 {
		log.Panic("The process slice to be simulated cannot be empty")
	}

	var time uint16
	// processes will be pushed onto this stack as they arrive, and so they will naturally be sorted
	// from last to first, which is perfect for LCFS
	processStack := sim.NewStack[*Process](len(*processes))
	// this is going to be another view into the underlying array, and by slicing it, we are able to
	// skip iterating over processes that have already arrived before
	unvisited := *processes

	// if there are processes that have not yet arrived, or ones that are waiting, continue
	for len(unvisited) != 0 || !processStack.Empty() {
		for i := range unvisited {
			if unvisited[i].arriveTime > time {
				// if a process arrives later than now, we know that all processes that have arrived up to this point have been iterated over
				// we can remove processes up to this one from our view of the array, as they have already been pushed onto the stack
				unvisited = unvisited[i:]
				break
			}
			// if a processes has arrived up to now, we push it onto the stack for it to wait for it's turn
			processStack.Push(&unvisited[i])
			// if the last process arrives, we need to empty our view
			if i == len(unvisited)-1 {
				unvisited = make([]Process, 0)
			}
		}

		// if there is a process waiting, we give it execution time
		if !processStack.Empty() {
			proc := processStack.Pop()
			time++
			proc.executionTimeLeft--

			if proc.executionTimeLeft == 0 {
				proc.waitTime = time - proc.arriveTime - proc.executionTime
				continue
			}
			// if a process is not done executing we push it back on the stack and check if there are any
			// newer once since this is preemptive lcfs
			processStack.Push(proc)
			continue
		}
		// if there are no processes waiting we just wait for them to arrive
		time++
	}
	return processes
}
func LCFS(processes *Slice) *Slice {
	if isSorted := slices.IsSortedFunc([]Process(*processes), func(a, b Process) int {
		return cmp.Compare(a.arriveTime, b.arriveTime)
	}); !isSorted {
		log.Panic("The process scheduling algorithms have to receive a slice of Processes sorted by arriveTime")
	}
	if *processes == nil {
		log.Panic("The process slice to be simulated cannot be nil")
	}
	if len(*processes) == 0 {
		log.Panic("The process slice to be simulated cannot be empty")
	}

	var time uint16
	// processes will be pushed onto this stack as they arrive, and so they will naturally be sorted
	// from last to first, which is perfect for LCFS
	processStack := sim.NewStack[*Process](len(*processes))
	// this is going to be another view into the underlying array, and by slicing it, we are able to
	// skip iterating over processes that have already arrived before
	unvisited := *processes

	// if there are processes that have not yet arrived, or ones that are waiting, continue
	for len(unvisited) != 0 || !processStack.Empty() {
		for i := range unvisited {
			if unvisited[i].arriveTime > time {
				// if a process arrives later than now, we know that all processes that have arrived up to this point have been iterated over
				// we can remove processes up to this one from our view of the array, as they have already been pushed onto the stack
				unvisited = unvisited[i:]
				break
			}
			// if a processes has arrived up to now, we push it onto the stack for it to wait for it's turn
			processStack.Push(&unvisited[i])
			// if the last process arrives, we need to empty our view
			if i == len(unvisited)-1 {
				unvisited = make([]Process, 0)
			}
		}

		// if there is a process waiting, we give it execution time
		if !processStack.Empty() {
			proc := processStack.Pop()
			// since this is non-preemptive lcfs we execute the process all at once
			// we can calculate the wait time right away since the process will be finished
			proc.waitTime = time - proc.arriveTime
			for proc.executionTimeLeft > 0 {
				time++
				proc.executionTimeLeft--
			}
			continue
		}
		// if there are no processes waiting we just wait for them to arrive
		time++
	}
	return processes
}
func PreemptiveSJF(processes *Slice) *Slice {
	if isSorted := slices.IsSortedFunc([]Process(*processes), func(a, b Process) int {
		return cmp.Compare(a.arriveTime, b.arriveTime)
	}); !isSorted {
		log.Panic("The process scheduling algorithms have to receive a slice of Processes sorted by arriveTime")
	}
	if *processes == nil {
		log.Panic("The process slice to be simulated cannot be nil")
	}
	if len(*processes) == 0 {
		log.Panic("The process slice to be simulated cannot be empty")
	}

	var time uint16
	// processes will be pushed onto this heap as they arrive, and they will get sorted by shortest job
	// I'm using a heap here because the time complexity for heapify is way better than sort
	// and we only need to know what the shortest job is
	processHeap := new(Heap)
	*processHeap = make([]*Process, 0, len(*processes))
	// this is going to be another view into the underlying array, and by slicing it, we are able to
	// skip iterating over processes that have already arrived before
	unvisited := *processes

	// if there are processes that have not yet arrived, or ones that are waiting, continue
	for len(unvisited) != 0 || len(*processHeap) != 0 {
		for i := range unvisited {
			if unvisited[i].arriveTime > time {
				// if a process arrives later than now, we know that all processes that have arrived up to this point have been iterated over
				// we can remove processes up to this one from our view of the array, as they have already been pushed onto the stack
				unvisited = unvisited[i:]
				break
			}
			// if a processes has arrived up to now, we push it onto the stack for it to wait for it's turn
			heap.Push(processHeap, &unvisited[i])
			// if the last process arrives, we need to empty our view
			if i == len(unvisited)-1 {
				unvisited = make([]Process, 0)
			}
		}

		// if there is a process waiting, we give it execution time
		if len(*processHeap) != 0 {
			proc := heap.Pop(processHeap).(*Process)
			time++
			proc.executionTimeLeft--

			if proc.executionTimeLeft == 0 {
				proc.waitTime = time - proc.arriveTime - proc.executionTime
				continue
			}
			// if the process is not done executing we push it back onto the Heap
			// which also fixes the ordering if it got broken
			// we will then see if this is still the shortest job
			heap.Push(processHeap, proc)
			continue
		}
		// if there are no processes waiting we just wait for them to arrive
		time++
	}

	return processes
}
func SJF(processes *Slice) *Slice {
	if isSorted := slices.IsSortedFunc([]Process(*processes), func(a, b Process) int {
		return cmp.Compare(a.arriveTime, b.arriveTime)
	}); !isSorted {
		log.Panic("The process scheduling algorithms have to receive a slice of Processes sorted by arriveTime")
	}
	if *processes == nil {
		log.Panic("The process slice to be simulated cannot be nil")
	}
	if len(*processes) == 0 {
		log.Panic("The process slice to be simulated cannot be empty")
	}

	var time uint16
	// processes will be pushed onto this heap as they arrive, and they will get sorted by shortest job
	// I'm using a heap here because the time complexity for heapify is way better than sort
	// and we only need to know what the shortest job is
	processHeap := new(Heap)
	*processHeap = make([]*Process, 0, len(*processes))
	// this is going to be another view into the underlying array, and by slicing it, we are able to
	// skip iterating over processes that have already arrived before
	unvisited := *processes

	// if there are processes that have not yet arrived, or ones that are waiting, continue
	for len(unvisited) != 0 || len(*processHeap) != 0 {
		for i := range unvisited {
			if unvisited[i].arriveTime > time {
				// if a process arrives later than now, we know that all processes that have arrived up to this point have been iterated over
				// we can remove processes up to this one from our view of the array, as they have already been pushed onto the stack
				unvisited = unvisited[i:]
				break
			}
			// if a processes has arrived up to now, we push it onto the stack for it to wait for it's turn
			heap.Push(processHeap, &unvisited[i])
			// if the last process arrives, we need to empty our view
			if i == len(unvisited)-1 {
				unvisited = make([]Process, 0)
			}
		}

		// if there is a process waiting, we give it execution time
		if len(*processHeap) != 0 {
			proc := heap.Pop(processHeap).(*Process)
			// we can calculate the wait time right away, since the process will be executed all the way
			proc.waitTime = time - proc.arriveTime
			for proc.executionTimeLeft > 0 {
				time++
				proc.executionTimeLeft--
			}
			continue
		}
		// if there are no processes waiting we just wait for them to arrive
		time++
	}
	return processes
}
