package process

import (
	"cmp"
	"fmt"
	"math/rand/v2"
	"reflect"
	"slices"
	"testing"

	"gotest.tools/v3/assert"
)

var t = &testing.T{}

type Process struct {
	id                uint16
	arriveTime        uint16
	executionTime     uint16
	executionTimeLeft uint16
	waitTime          uint16
}

// Gen generates a slice of processes, sorted by arriveTime
func Gen(num uint16, maxArriveTime uint16, maxExecutionTime uint16) *Slice {
	assert.Assert(t, num != 0, "Cannot generate 0 processes")
	assert.Assert(t, maxExecutionTime != 0, "Cannot generate processes with zero execution time")

	var processes Slice = make([]Process, num)
	for i := range processes {
		processes[i] = Process{id: uint16(i),
			arriveTime:    uint16(rand.UintN(uint(maxArriveTime))),
			executionTime: uint16(1 + rand.UintN(uint(maxExecutionTime-1)))}
		processes[i].executionTimeLeft = processes[i].executionTime
	}

	slices.SortFunc(processes, func(a, b Process) int {
		return cmp.Compare(a.arriveTime, b.arriveTime)
	})
	return &processes
}

type Slice []Process

var processNumFields = reflect.TypeOf(Process{}).NumField()

// Records implements the Recorder interface
func (s *Slice) Records() (records [][]string) {
	assert.Assert(t, s != nil, "The slice to get records from cannot be nil")
	assert.Assert(t, len(*s) != 0, "The slice to get records from cannot be empty")

	records = make([][]string, len(*s)+1)
	for i := range records {
		records[i] = make([]string, processNumFields)
	}

	for i := range records[0] {
		records[0][i] = reflect.TypeOf(Process{}).Field(i).Name
	}

	vals := records[1:]
	for i, process := range *s {
		for j := range vals[i] {
			field := fmt.Sprint(reflect.ValueOf(process).Field(j))
			vals[i][j] = field
		}
	}
	return records
}

// Copy makes a deep copy of the passed in Slice
func (s *Slice) Copy() *Slice {
	assert.Assert(t, s != nil, "The slice to copy cannot be nil")
	assert.Assert(t, len(*s) != 0, "The slice to copy cannot be empty")

	c := make([]Process, len(*s))
	for i := range *s {
		c[i] = (*s)[i]
	}
	return (*Slice)(&c)
}

// Heap implements the container.Heap.Interface to get a process min Heap sorted by execution time left (for SJF)
type Heap []*Process

func (h *Heap) Len() int {
	assert.Assert(t, h != nil, "The underlying slice of a heap cannot be nil")
	return len(*h)
}
func (h *Heap) Less(i, j int) bool {
	assert.Assert(t, h != nil, "The underlying slice of a heap cannot be nil")
	return (*h)[i].executionTimeLeft < (*h)[j].executionTimeLeft
}
func (h *Heap) Swap(i, j int) {
	assert.Assert(t, h != nil, "The underlying slice of a heap cannot be nil")
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}
func (h *Heap) Push(x any) {
	assert.Assert(t, h != nil, "The underlying slice of a heap cannot be nil")
	*h = append(*h, x.(*Process))
}
func (h *Heap) Pop() any {
	assert.Assert(t, h != nil, "The underlying slice of a heap cannot be nil")
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
