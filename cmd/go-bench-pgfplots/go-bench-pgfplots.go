package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"

	gobenchmarkpgfplots "github.com/oxisto/go-benchmark-pgfplots"
)

var file = flag.String("file", "bench.txt", "The benchmark file")
var sep = flag.String("sep", "/", "The seperator in the benchmark name")

func main() {
	f, err := os.OpenFile(*file, os.O_RDONLY, 0600)
	if err != nil {
		log.Printf("Could not open file: %v", err)
		return
	}

	doc, err := ioutil.ReadAll(f)
	if err != nil {
		log.Printf("Could not read file: %v", err)
		return
	}

	res, err := gobenchmarkpgfplots.Convert(doc, *sep, time.Millisecond)
	if err != nil {
		log.Printf("Could not convert: %v", err)
		return
	}

	err = gobenchmarkpgfplots.Serialize(res)
	if err != nil {
		log.Printf("Could not save results: %v", err)
		return
	}
}
