package page

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"reflect"
)

type Page struct {
	id           uint16
	timesUsed    uint16
	pageFaultAt  []uint16
	swappedOutAt []uint16
}

// gen generates a random pattern of referencing a given amount of pages a given number of times
func gen(numPages, len uint16) (referencePattern []uint16) {
	if len == 0 {
		log.Panic("The page reference pattern length must be greater than zero")
	}
	if numPages == 0 {
		log.Panic("The number of pages to simulate must be larger than zero")
	}

	referencePattern = make([]uint16, len)
	for i := range referencePattern {
		referencePattern[i] = uint16(rand.UintN(uint(numPages)))
	}
	return referencePattern
}

// SaveReferencePattern saves the reference pattern to an output file in a .csv format
func SaveReferencePattern(referencePattern []uint16) {
	if referencePattern == nil {
		log.Panic("The reference pattern slice cannot be nil")
	}
	if len(referencePattern) == 0 {
		log.Panic("The reference pattern has to contain something")
	}

	record := make([]string, len(referencePattern))
	for i, index := range referencePattern {
		record[i] = fmt.Sprint(index)
	}

	outDirPath := "../in/"
	if err := os.MkdirAll(outDirPath, 0775); err != nil {
		log.Panic(err)
	}
	filename := outDirPath + "pageReferencePattern"

	csvFile, err := os.Create(filename + ".csv")
	defer csvFile.Close()
	if err != nil {
		log.Panic(err)
	}
	w := csv.NewWriter(csvFile)
	if err := w.Write(record); err != nil {
		log.Panic(err)
	}
	w.Flush()
}

type Slice []Page

var pageNumFields = reflect.TypeOf(Page{}).NumField()

// Records implements the Recorder interface
func (pages *Slice) Records() (records [][]string) {
	if pages == nil {
		log.Panic("The page slice for extracting records cannot be nil")
	}
	if len(*pages) == 0 {
		log.Panic("The page slice for extracting records cannot be empty")
	}

	records = make([][]string, len(*pages)+1)
	for i := range records {
		records[i] = make([]string, pageNumFields)
	}

	for i := range records[0] {
		records[0][i] = reflect.TypeOf(Page{}).Field(i).Name
	}

	vals := records[1:]
	for i, page := range *pages {
		for j := range vals[i] {
			field := fmt.Sprint(reflect.ValueOf(page).Field(j))
			vals[i][j] = field
		}
	}
	return records
}

// Heap implements the container.Heap.Interface to get a page min Heap sorted by least times used (for LFU)
type Heap []*Page

func (h *Heap) Len() int {
	if h == nil {
		log.Panic("The underlying slice of a heap cannot be nil")
	}

	return len(*h)
}

func (h *Heap) Less(i, j int) bool {
	if h == nil {
		log.Panic("The underlying slice of a heap cannot be nil")
	}

	return (*h)[i].timesUsed < (*h)[j].timesUsed
}

func (h *Heap) Swap(i, j int) {
	if h == nil {
		log.Panic("The underlying slice of a heap cannot be nil")
	}

	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *Heap) Push(x any) {
	if h == nil {
		log.Panic("The underlying slice of a heap cannot be nil")
	}

	*h = append(*h, x.(*Page))
}

func (h *Heap) Pop() any {
	if h == nil {
		log.Panic("The underlying slice of a heap cannot be nil")
	}

	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
