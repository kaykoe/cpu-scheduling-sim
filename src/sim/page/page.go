package page

import (
	"reflect"
	"strconv"
)

type Page struct {
	id       uint16
	lastUsed uint16
}

type Slice []Page

var pageNumFields = reflect.TypeOf(Page{}).NumField()

// Records implements the Recorder interface
func (s *Slice) Records() (records [][]string) {
	if s != nil {
		records = make([][]string, len(*s)+1)
		for i := range records {
			records[i] = make([]string, pageNumFields)
		}

		for i := range records[0] {
			records[0][i] = reflect.TypeOf(Page{}).Field(i).Name
		}

		for i, page := range *s {
			for j := range records[i+1] {
				field := strconv.FormatUint(reflect.ValueOf(page).Field(j).Uint(), 10)
				records[i+1][j] = field
			}
		}
		return records
	}
	return nil
}

func (s *Slice) Copy() *Slice {
	c := make([]Page, len(*s))
	for i := range *s {
		c[i] = (*s)[i]
	}
	return (*Slice)(&c)
}

type Heap []*Page

func (p Heap) Len() int {
	return len(p)
}

func (p Heap) Less(i, j int) bool {
	return p[i].lastUsed < p[j].lastUsed
}

func (p Heap) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Heap) Push(x any) {
	p = append(p, x.(*Page))
}

func (p Heap) Pop() any {
	old := p
	n := len(old)
	x := old[n-1]
	p = old[:n-1]
	return x
}
