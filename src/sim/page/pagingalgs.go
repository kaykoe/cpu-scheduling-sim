package page

import (
	"container/heap"
	"log"
	"maps"
	"math"
	"src/sim"
)

type Alg func(referencePattern []uint16) *Slice

// Sim runs a simulation of a randomly generated reference pattern using the strategies in the alg slice
func Sim(numPages, referencePatternLen uint16, algs ...Alg) (referencePattern []uint16, res []*Slice) {
	if numPages == 0 {
		log.Panic("The number of pages to simulate cannot be zero")
	}
	if referencePatternLen == 0 {
		log.Panic("The page reference pattern length must be greater than zero")
	}
	if algs == nil {
		log.Panic("The alg slice cannot be nil")
	}
	if len(algs) == 0 {
		log.Panic("The number of algorithms to simulate cannot be zero")
	}

	referencePattern = gen(numPages, referencePatternLen)
	res = make([]*Slice, len(algs))
	for i, alg := range algs {
		res[i] = alg(referencePattern)
	}
	return referencePattern, res
}

func FIFO(referencePattern []uint16) *Slice {
	if referencePattern == nil {
		log.Panic("The reference pattern slice cannot be nil")
	}
	if len(referencePattern) == 0 {
		log.Panic("The reference pattern has to contain something")
	}
	if len(referencePattern) > math.MaxUint16 {
		log.Panicf("The reference pattern length cannot exceed %d", math.MaxUint16)
	}

	// the page table stores whether a page is in memory
	pageTable := make(map[uint16]bool)
	memory := make(map[uint16]*Page)
	swap := make(map[uint16]*Page)
	const FRAME_SIZE = 16
	// the delete queue stores the indices of pages that are in memory in the order they were referenced
	// this is perfect for implementing fifo
	deleteQueue := sim.NewQueue[uint16](FRAME_SIZE)

	for _, page := range referencePattern {
		pageTable[page] = false
	}
	for page := range maps.Keys(pageTable) {
		swap[page] = &Page{
			page,
			0,
			make([]uint16, 0, len(referencePattern)),
			make([]uint16, 0, len(referencePattern))}
	}

	for i, page := range referencePattern {
		if inMemory := pageTable[page]; !inMemory {
			// if there is no space left in memory we need to move a page to swap
			if len(memory) == FRAME_SIZE {
				victimPage := deleteQueue.Pop()
				swap[victimPage] = memory[victimPage]
				delete(memory, victimPage)
				swap[victimPage].swappedOutAt = append(swap[victimPage].swappedOutAt, uint16(i))
				pageTable[victimPage] = false
			}

			memory[page] = swap[page]
			delete(swap, page)
			pageTable[page] = true
			memory[page].pageFaultAt = append(memory[page].pageFaultAt, uint16(i))
			deleteQueue.Push(page)
			continue
		}
	}
	res := Slice(make([]Page, 0, len(pageTable)))
	for page := range maps.Values(swap) {
		res = append(res, *page)
	}
	for page := range maps.Values(memory) {
		res = append(res, *page)
	}
	return &res
}

func LFU(referencePattern []uint16) *Slice {
	if referencePattern == nil {
		log.Panic("The reference pattern slice cannot be nil")
	}
	if len(referencePattern) == 0 {
		log.Panic("The reference pattern has to contain something")
	}
	if len(referencePattern) > math.MaxUint16 {
		log.Panicf("The reference pattern length cannot exceed %d", math.MaxUint16)
	}

	// the page table stores whether a page is in memory
	pageTable := make(map[uint16]bool)
	memory := make(map[uint16]*Page)
	swap := make(map[uint16]*Page)
	const FRAME_SIZE = 16
	// this heap stores pages that are in memory, sorted by the amount of times they have been used since getting swapped into memory
	// a heap is used because heapify has better time complexity than sort, and we only need the smallest element
	deleteHeap := new(Heap)
	*deleteHeap = make([]*Page, 0, FRAME_SIZE)

	for _, page := range referencePattern {
		pageTable[page] = false
	}
	for page := range maps.Keys(pageTable) {
		swap[page] = &Page{
			page,
			0,
			make([]uint16, 0, len(referencePattern)),
			make([]uint16, 0, len(referencePattern))}
	}

	for i, page := range referencePattern {
		if inMemory := pageTable[page]; !inMemory {
			// if there is no space left in memory we need to move a page to swap
			if deleteHeap.Len() == FRAME_SIZE {
				// here we heapify the heap so that it takes into account uses that did not require swapping
				heap.Init(deleteHeap)
				victimPage := heap.Pop(deleteHeap).(*Page)
				delete(memory, victimPage.id)
				swap[victimPage.id] = victimPage
				victimPage.swappedOutAt = append(victimPage.swappedOutAt, uint16(i))
				// whenever a page is swapped out it's counter is reset, so that the algorithm responds to locality changes better
				victimPage.timesUsed = 0
				pageTable[victimPage.id] = false
			}

			memory[page] = swap[page]
			heap.Push(deleteHeap, memory[page])
			delete(swap, page)
			pageTable[page] = true
			memory[page].pageFaultAt = append(memory[page].pageFaultAt, uint16(i))
		}
		// if the page was already in memory we just increment the amount of times it was used
		memory[page].timesUsed++
	}
	res := Slice(make([]Page, 0, len(pageTable)))
	for page := range maps.Values(swap) {
		res = append(res, *page)
	}
	for page := range maps.Values(memory) {
		res = append(res, *page)
	}
	return &res
}
