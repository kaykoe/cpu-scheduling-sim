package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"
)

// Recorder is an interface used to get a records slice for easy formatted output
type Recorder interface {
	Records() [][]string
}

// writeRecords writes the records of the passed in Recorder to the output io.Writer encoded as csv,
// or with tab alignment if the io.Writer is a tabwriter.Writer
func writeRecords(r Recorder, output io.Writer) {
	w := csv.NewWriter(output)
	if reflect.TypeOf(output) == reflect.TypeOf(&tabwriter.Writer{}) {
		w.Comma = '\t'
	}
	records := r.Records()
	if err := w.WriteAll(records); err != nil {
		log.Panic(err)
	}
}

func Save(r Recorder, outDir string) {
	outDirPath := "../" + outDir + "/"
	if err := os.MkdirAll(outDirPath, 0775); err != nil {
		log.Panic(err)
	}
	// here we name the output files with the name of the package they come from
	rpath := reflect.Indirect(reflect.ValueOf(r)).Type().PkgPath()
	rname := rpath[strings.LastIndexByte(rpath, '/')+1:]
	filename := outDirPath + rname

	// we write 2 files, one is  csv file for easier work with python,
	// the other one is a txt file with tab aligned columns for readability
	csvFile, err := os.Create(filename + ".csv")
	defer csvFile.Close()
	if err != nil {
		log.Panic(err)
	}
	writeRecords(r, csvFile)

	textFile, err := os.Create(filename + ".txt")
	defer textFile.Close()
	if err != nil {
		log.Panic(err)
	}
	tw := tabwriter.NewWriter(textFile, 0, 4, 3, ' ', 0)
	writeRecords(r, tw)
	if err = tw.Flush(); err != nil {
		log.Panic(err)
	}
}
