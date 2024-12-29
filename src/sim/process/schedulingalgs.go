package process

import (
	"container/heap"
	"log"
	"src/sim"
)

type Alg func(p *Slice) *Slice

// Sim runs a simulation of the given processes using the strategies in the alg slice
func Sim(p *Slice, a []Alg) (res []*Slice) {
	if len(a) == 0 {
		log.Panic("empty algorithm slice")
	}

	res = make([]*Slice, len(a))

	for i, alg := range a {
		res[i] = alg(p.Copy())
	}
	return res
}

func PreemptiveLCFS(p *Slice) *Slice {
	var time uint16
	s := sim.NewStack[*Process]()
	unvisited := 0

	for unvisited != len(*p) || !s.Empty() {
		for i := unvisited; i < len(*p); i++ {
			if (*p)[i].arriveTime > time {
				unvisited = i
				break
			}
			s.Push(&(*p)[i])
			if i == len(*p)-1 {
				unvisited = len(*p)
			}
		}

		if !s.Empty() {
			proc := s.Pop()
			time++
			proc.executionTimeLeft--
			if proc.executionTimeLeft == 0 {
				proc.waitTime = time - proc.arriveTime - proc.executionTime
				continue
			}
			s.Push(proc)
			continue
		}
		time++
	}
	return p
}
func LCFS(p *Slice) *Slice {
	var time uint16
	s := sim.NewStack[*Process]()
	unvisited := 0

	for unvisited != len(*p) || !s.Empty() {
		for i := unvisited; i < len(*p); i++ {
			if (*p)[i].arriveTime > time {
				unvisited = i
				break
			}
			s.Push(&(*p)[i])
			if i == len(*p)-1 {
				unvisited = len(*p)
			}
		}

		if !s.Empty() {
			proc := s.Pop()
			proc.waitTime = time - proc.arriveTime
			for proc.executionTimeLeft > 0 {
				time++
				proc.executionTimeLeft--
			}
			continue
		}
		time++
	}
	return p
}
func PreemptiveSJF(p *Slice) *Slice {
	var time uint16
	unvisited := 0
	h := new(Heap)
	*h = make([]*Process, 0, len(*p))

	for unvisited != len(*p) || len(*h) != 0 {
		for i := unvisited; i < len(*p); i++ {
			if (*p)[i].arriveTime > time {
				unvisited = i
				break
			}
			heap.Push(h, &(*p)[i])
			if i == len(*p)-1 {
				unvisited = len(*p)
			}
		}

		if len(*h) != 0 {
			proc := heap.Pop(h).(*Process)
			time++
			proc.executionTimeLeft--
			if proc.executionTimeLeft == 0 {
				proc.waitTime = time - proc.arriveTime - proc.executionTime
				continue
			}
			heap.Push(h, proc)
			continue
		}
		time++
	}

	return p
}
func SJF(p *Slice) *Slice {
	var time uint16
	unvisited := 0
	h := new(Heap)
	*h = make([]*Process, 0, len(*p))

	for unvisited != len(*p) || len(*h) != 0 {
		for i := unvisited; i < len(*p); i++ {
			if (*p)[i].arriveTime > time {
				unvisited = i
				break
			}
			heap.Push(h, &(*p)[i])
			if i == len(*p)-1 {
				unvisited = len(*p)
			}
		}

		if len(*h) != 0 {
			proc := heap.Pop(h).(*Process)
			proc.waitTime = time - proc.arriveTime
			for proc.executionTimeLeft > 0 {
				time++
				proc.executionTimeLeft--
			}
			continue
		}
		time++
	}
	return p
}
