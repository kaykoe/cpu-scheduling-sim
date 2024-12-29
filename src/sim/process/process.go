package process

import (
	"cmp"
	"log"
	"math/rand/v2"
	"reflect"
	"slices"
	"strconv"
)

type Process struct {
	id                uint16
	arriveTime        uint16
	executionTime     uint16
	executionTimeLeft uint16
	waitTime          uint16
}

// Gen generates a slice of processes, sorted by arriveTime
func Gen(num uint16, maxArriveTime uint16, maxExecutionTime uint16) *Slice {
	if num == 0 {
		log.Panic("process.Gen called with num = 0")
	}
	if maxExecutionTime == 0 {
		log.Panic("process.Gen called with maxExecutionTime = 0")
	}
	var p Slice = make([]Process, num)
	for i := range p {
		p[i] = Process{id: uint16(i),
			arriveTime:    uint16(rand.UintN(uint(maxArriveTime))),
			executionTime: uint16(1 + rand.UintN(uint(maxExecutionTime-1)))} // NormFloat64() for normal distribution
		p[i].executionTimeLeft = p[i].executionTime
	}
	slices.SortFunc(p, func(a, b Process) int {
		return cmp.Compare(a.arriveTime, b.arriveTime)
	})
	return &p
}

type Slice []Process

var processNumFields = reflect.TypeOf(Process{}).NumField()

// Records implements the Recorder interface
func (s *Slice) Records() (records [][]string) {
	if s != nil {
		records = make([][]string, len(*s)+1)
		for i := range records {
			records[i] = make([]string, processNumFields)
		}

		for i := range records[0] {
			records[0][i] = reflect.TypeOf(Process{}).Field(i).Name
		}

		for i, process := range *s {
			for j := range records[i+1] {
				field := strconv.FormatUint(reflect.ValueOf(process).Field(j).Uint(), 10)
				records[i+1][j] = field
			}
		}
		return records
	}
	return nil
}

// Copy makes a deep copy of the passed in Slice
func (s *Slice) Copy() *Slice {
	if s != nil {
		c := make([]Process, len(*s))
		for i := range *s {
			c[i] = (*s)[i]
		}
		return (*Slice)(&c)
	}
	return nil
}

// Heap implements the container.Heap.Interface to get a process min Heap sorted by execution time left (for SJF)
type Heap []*Process

func (p *Heap) Len() int {
	return len(*p)
}
func (p *Heap) Less(i, j int) bool {
	return (*p)[i].executionTimeLeft < (*p)[j].executionTimeLeft
}
func (p *Heap) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}
func (p *Heap) Push(x any) {
	*p = append(*p, x.(*Process))
}
func (p *Heap) Pop() any {
	old := *p
	n := len(old)
	x := old[n-1]
	*p = old[:n-1]
	return x
}
