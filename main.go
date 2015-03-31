package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
)

var (
	fl_header    = flag.Bool("h", false, "Output header")
	fl_npoints   = flag.Int("n", 0, "Number of points")
	fl_delta     = flag.Float64("d", 1, "Distance between points")
	fl_clength   = flag.Float64("w", 0, "Correlation length")
	fl_amplitude = flag.Float64("a", 1, "RMS amplitude")
	fl_kern      = flag.Float64("k", 5, "Kernel cut-off in correlation lengths")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	if *fl_header {
		fmt.Println("# position (m) \tE (J)")
	}

	N := *fl_npoints
	delta := *fl_delta
	width := *fl_clength
	C := int(*fl_kern * width / delta)

	if N <= 0 || delta <= 0 || width < 0 {
		flag.Usage()
		os.Exit(1)
	}

	// make input noise
	in := make([]float64, N)
	for i := range in {
		in[i] = rand.Float64() - 0.5
	}

	out := make([]float64, N)
	if width != 0 {
		// make kernel
		klen := 2*C + 1
		kern := make([]float64, klen)
		for d := -C; d <= C; d++ {
			x := (2 * float64(d) * delta / width)
			kern[C-d] = math.Exp(-x * x)
		}

		// convolution: out = in * kern
		for i := range out {
			for d := -C; d <= C; d++ {
				j := i + d
				for j < 0 {
					j += N
				}
				for j >= N {
					j -= N
				}
				in := in[j]
				k := kern[C-d]
				out[i] += in * k
			}
		}
	} else {
		out = in
	}

	// determine stddev
	sumSq := 0.0
	for _, x := range out {
		sumSq += x * x
	}
	stddev := math.Sqrt(sumSq / float64(N))

	// normalize to stddev
	scale := *fl_amplitude / stddev
	for i := range out {
		out[i] *= scale
	}

	for i := range out {
		fmt.Println(float64(i)*delta, "\t", out[i])
	}

}
