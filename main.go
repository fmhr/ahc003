package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
)

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

var sc = bufio.NewScanner(os.Stdin)
var buff []byte

func nextInt() int {
	sc.Scan()
	i, err := strconv.Atoi(sc.Text())
	if err != nil {
		panic(err)
	}
	return i
}

func nextFloat64() float64 {
	sc.Scan()
	f, err := strconv.ParseFloat(sc.Text(), 64)
	if err != nil {
		panic(err)
	}
	return f
}

func init() {
	sc.Split(bufio.ScanWords)
	sc.Buffer(buff, bufio.MaxScanTokenSize*1024)
	log.SetFlags(log.Lshortfile)
}

// https://golang.org/pkg/runtime/pprof/
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// ... rest of the program ...
	// localTester()
	solver()

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

type Point struct {
	i, j int
}

type Ask struct {
	s Point
	t Point
	a float64
	e float64
}

func localTester() {
	// input testcase
	var h [30][30]int
	var v [30][30]int
	for i := 0; i < 30; i++ {
		for j := 0; j < 30; j++ {
			if i == j {
				h[i][j] = 0
			} else {
				h[i][j] = nextInt()
			}
		}
	}
	for i := 0; i < 30; i++ {
		for j := 0; j < 30; j++ {
			if i == j {
				v[i][j] = 0
			} else {
				v[i][j] = nextInt()
			}
		}
	}
	asks := make([]Ask, 1000)
	for i := 0; i < 1000; i++ {
		asks[i].s.i = nextInt()
		asks[i].s.j = nextInt()
		asks[i].t.i = nextInt()
		asks[i].t.j = nextInt()
		asks[i].a = nextFloat64()
		asks[i].e = nextFloat64()
	}
	var score float64
	for i := 0; i < 1000; i++ {
		route := ask(asks[i].s.i, asks[i].s.j, asks[i].t.i, asks[i].t.j)
		fmt.Println(route)
		log.Println(route)
		dest := restore(asks[i].s, asks[i].t, route)
		score += math.Pow(0.998, float64(1000-i)) * asks[i].a / float64(dest)
	}
	score = math.Round(2312311 * score)
	log.Println(score)
}

func restore(start, goal Point, route string) (dest int) {
	return
}

func ask(si, sj, ti, tj int) (route string) {
	log.Println(si, sj, ti, tj)
	if si-ti < 0 {
		route += strings.Repeat("D", ti-si)
	} else {
		route += strings.Repeat("U", si-ti)
	}
	if sj-tj < 0 {
		route += strings.Repeat("R", tj-sj)
	} else {
		route += strings.Repeat("L", sj-tj)
	}
	return route
}

func solver() {
	for i := 0; i < 1000; i++ {
		si := nextInt()
		sj := nextInt()
		ti := nextInt()
		tj := nextInt()
		route := ask(si, sj, ti, tj)
		fmt.Println(route)
		_ = nextInt()
	}
}
