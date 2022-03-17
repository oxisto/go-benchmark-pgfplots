package gobenchmarkpgfplots

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Benchmark map[string]*Dataset

type Dataset struct {
	ID      string
	Results []*Result
}

type Result struct {
	X int64
	Y float64
}

// Convert converts go benchmark output of input into data files for pgfplots.
// One can specificy a seperator in sep, which specifies how benchmark parameters
// are encoded into the benchmark name.
//
// The assumption is that the benchmark name is structured like BenchmarkAssessEvidence/1/2
// in which case 1 is taken as the x value and 2 is taken as the identifier for the dataset.
func Convert(input []byte, sep string, prec time.Duration) (results map[string]*Benchmark, err error) {
	var (
		x  int64
		y  float64
		ns int64
		id string
		s  string
	)

	results = make(map[string]*Benchmark)

	// Split lines
	lines := strings.Split(string(input), "\n")

	for _, line := range lines {
		// Line must start with Benchmark
		if strings.Index(line, "Benchmark") != 0 {
			continue
		}

		columns := strings.Split(line, "\t")

		// Must be 5 columns (we assume -benchmem for now)
		if len(columns) != 5 {
			continue
		}

		name, _, _ := strings.Cut(columns[0], "/")

		s, _, _ = strings.Cut(columns[0], "-")

		params := strings.Split(s, sep)[1:]

		// For now, we only support exactly 2 parameters
		if len(params) != 2 {
			continue
		}

		log.Printf("%s\n", line)

		// The first parameter is our x value
		x, err = strconv.ParseInt(params[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse int: %v", err)
		}

		// The second parameter identifies our dataset
		id = params[1]

		// The third column is our y value
		ns, err = strconv.ParseInt(value(columns[2]), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse int: %v", err)
		}

		y = float64(ns) / float64(prec)

		// Lookup the benchmark by name
		b, ok := results[name]
		if !ok {
			b = &Benchmark{}
			results[name] = b
		}

		ds, ok := (*b)[id]
		if !ok {
			ds = &Dataset{
				ID: id,
			}
			(*b)[id] = ds
		}

		res := Result{
			X: x,
			Y: y,
		}

		ds.Results = append(ds.Results, &res)
	}

	return
}

func Serialize(results map[string]*Benchmark) (err error) {
	for name, b := range results {
		log.Printf("Writing benchmark %s...\n", name)

		for _, ds := range *b {
			fn := fmt.Sprintf("data-%s-%s.dat", name, ds.ID)

			log.Printf("Writing dataset %s...\n", fn)

			f, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, 0600)
			defer f.Close()

			if err != nil {
				return fmt.Errorf("could not open file: %v", err)
			}

			for _, res := range ds.Results {
				fmt.Fprintf(f, "%v %v\n", res.X, res.Y)
			}
		}
	}

	return
}

func value(s string) string {
	s, _, _ = strings.Cut(s, " ns/op")

	return strings.TrimLeft(s, " ")
}
