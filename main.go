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

var local = flag.String("local", "default", "on local")

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

	if *local != "" {
		fmt.Println("localだよ")
		localTester()
	} else {
		solver()
	}

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

var direction = []string{"D", "R", "U", "L"}
var reverse = map[byte]int{'D': 0, 'R': 1, 'U': 2, 'L': 3}
var di = [4]int{1, 0, -1, 0}
var dj = [4]int{0, 1, 0, -1}

var h [30][30]int
var v [30][30]int

func localTester() {
	// input testcase
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
	for k := 0; k < 1000; k++ {
		route := ask(asks[k].s.i, asks[k].s.j, asks[k].t.i, asks[k].t.j)
		//	fmt.Println(route)
		dest := restore(asks[k].s, asks[k].t, route)
		tmp := math.Pow(0.998, float64(1000-(k+1))) * (asks[k].a / float64(dest))
		log.Println(tmp)
		score += tmp
	}
	score = math.Round(2312311 * score)
	log.Printf("%f\n", score)
}

func restore(start, goal Point, route string) (dest int) {
	now := start
	for i := 0; i < len(route); i++ {
		var next Point
		d := reverse[route[i]]
		next.i = now.i + di[d]
		next.j = now.j + dj[d]
		switch d {
		case 0:
			dest += h[now.i][now.j]
		case 1:
			dest += v[now.i][now.j]
		case 2:
			dest += h[next.i][next.j]
		case 3:
			dest += v[next.i][next.j]
		}
		now = next
	}
	log.Println(now, goal)
	return
}

func ask(si, sj, ti, tj int) (route string) {
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
