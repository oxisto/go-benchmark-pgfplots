package gobenchmarkpgfplots

import (
	"reflect"
	"testing"
	"time"
)

func TestConvert(t *testing.T) {
	type args struct {
		input []byte
		sep   string
		prec  time.Duration
	}
	tests := []struct {
		name        string
		args        args
		wantResults map[string]*Benchmark
		wantErr     bool
	}{
		{
			name: "benchmark with 2 parameters",
			args: args{
				input: []byte(`
BenchmarkTest/1/2-8         	     3	   4000000000 ns/op	  5 B/op	    6 allocs/op
`),
				sep:  "/",
				prec: time.Second,
			},
			wantResults: map[string]*Benchmark{
				"BenchmarkTest": {
					"2": {
						ID: "2",
						Results: []*Result{
							{
								X: 1,
								Y: 4,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResults, err := Convert(tt.args.input, tt.args.sep, tt.args.prec)
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotResults, tt.wantResults) {
				t.Errorf("Convert() = %v, want %v", gotResults, tt.wantResults)
			}
		})
	}
}
