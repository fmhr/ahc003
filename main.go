package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
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

var local = flag.Bool("local", false, "if local")

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
	if *local {
		localTester()
	} else {
		log.Println("not local")
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

func (p *Point) move(m byte) {
	p.i = p.i + di[direction[m]]
	p.j = p.j + dj[direction[m]]
}

type Ask struct {
	s Point
	t Point
	a float64
	e float64
}

var direction = map[byte]int{'D': 0, 'R': 1, 'U': 2, 'L': 3}
var di = [4]int{1, 0, -1, 0}
var dj = [4]int{0, 1, 0, -1}

var h [30][30]int
var v [30][30]int

func localTester() {
	// input testcase
	for i := 0; i < 30; i++ {
		for j := 0; j < 29; j++ {
			h[i][j] = nextInt()
		}
	}
	for i := 0; i < 29; i++ {
		for j := 0; j < 30; j++ {
			v[i][j] = nextInt()
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
		path := query(asks[k].s.i, asks[k].s.j, asks[k].t.i, asks[k].t.j)
		fmt.Println(path)
		b := compute_path_length(asks[k].s, asks[k].t, path)
		score = score*0.998 + (asks[k].a)/float64(b)
	}
	score = math.Round(2312311.0 * score)
	log.Printf("%f\n", score)
}

func compute_path_length(start, goal Point, route string) (dest int) {
	now := start
	for i := 0; i < len(route); i++ {
		var next Point
		d := direction[route[i]]
		next.i = now.i + di[d]
		next.j = now.j + dj[d]
		switch d {
		case 0: // D
			dest += v[now.i][now.j]
		case 1: // R
			dest += h[now.i][now.j]
		case 2: // U
			dest += v[next.i][next.j]
		case 3: // L
			dest += h[next.i][next.j]
		}
		now = next
	}
	return
}

// 直線的に動く暫定
func query(si, sj, ti, tj int) (route string) {
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

func randomSolver(q *QueryRecord, pr *PathRecord) []byte {
	si := q.start.i
	sj := q.start.j
	ti := q.stop.i
	tj := q.stop.j
	var now Point
	now.i = si
	now.j = sj
	rb := make([]byte, absInt(si-ti)+absInt(sj-tj))
	cnt := 0
	for !(now.i == ti && now.j == tj) {
		r := ""
		if now.i < ti {
			r += "D"
		} else if now.i > ti {
			r += "U"
		}
		if now.j < tj {
			r += "R"
		} else if now.j > tj {
			r += "L"
		}
		if len(r) == 0 {
			panic("Errorrrrr")
		}
		for i := 0; i < len(r); i++ {
			pr.AddAppeared(now, r[i])
		}
		rb[cnt] = r[rand.Intn(len(r))]
		pr.AddAppeared(now, rb[cnt])
		now.move(rb[cnt])
		cnt++
	}
	return rb
}

func sampleUCB(p Path) float64 {
	v := 0.0
	//log.Println(math.Sqrt(math.Log(float64(p.numOfAppeared)) / float64(2*p.numOfSelected)))
	//v = float64(p.SampleAverage)
	v = float64(p.SampleAverage) - math.Sqrt(math.Log(float64(p.numOfAppeared))/float64(2*p.numOfSelected))
	return v
}

func greedySolver(q *QueryRecord, pr *PathRecord) []byte {
	si := q.start.i
	sj := q.start.j
	ti := q.stop.i
	tj := q.stop.j
	var now Point
	now.i = si
	now.j = sj
	rb := make([]byte, absInt(si-ti)+absInt(sj-tj))
	cnt := 0
	for !(now.i == ti && now.j == tj) {
		r := ""
		if now.i < ti {
			r += "D"
		} else if now.i > ti {
			r += "U"
		}
		if now.j < tj {
			r += "R"
		} else if now.j > tj {
			r += "L"
		}
		if len(r) == 0 {
			panic("Errorrrrr")
		}
		for i := 0; i < len(r); i++ {
			pr.AddAppeared(now, r[i])
		}
		//
		ps := make([]Path, len(r))
		nouse := -1
		for i := 0; i < len(r); i++ {
			y, x := getIj(now, r[i])
			ps[i] = pr.getPath(y, x, r[i])
			ps[i].index = i
			if ps[i].numOfSelected == 0 {
				nouse = i
			}
		}
		if nouse != -1 {
			rb[cnt] = r[nouse]
		} else {
			sort.Slice(ps, func(i, j int) bool {
				return sampleUCB(ps[i]) < sampleUCB(ps[j])
			})
			rb[cnt] = r[ps[0].index]
		}
		//rb[cnt] = r[rand.Intn(len(r))]
		pr.AddSelected(now, rb[cnt])
		now.move(rb[cnt])
		cnt++
	}
	return rb
}

// worchal floyd ----------------------------------------------
const inf int = 2 << 29

type Graph struct {
	cost [900][900]int
	size int
}

var next [900][900]int
var g Graph

func toindex(i, j int) int {
	return i*30 + j
}
func fromindex(k int) (int, int) {
	return k / 30, k % 30
}
func buildGraph(pr PathRecord) {
	for i := 0; i < 900; i++ {
		for j := 0; j < 900; j++ {
			g.cost[i][j] = inf
		}
	}
	for i := 0; i < 30; i++ {
		for j := 0; j < 29; j++ {
			a := toindex(i, j)
			b := toindex(i, j+1)
			g.cost[a][b] = pr.h[i][j].SampleAverage
			g.cost[b][a] = pr.h[i][j].SampleAverage
		}
	}
	for i := 0; i < 29; i++ {
		for j := 0; j < 30; j++ {
			a := toindex(i, j)
			b := toindex(i+1, j)
			g.cost[a][b] = pr.h[i][j].SampleAverage
			g.cost[b][a] = pr.h[i][j].SampleAverage
		}
	}
}

func warchalFloyd() {
	for i := 0; i < 900; i++ {
		for j := 0; j < 900; j++ {
			next[i][j] = j
		}
	}
	for k := 0; k < 900; k++ {
		for i := 0; i < 900; i++ {
			for j := 0; j < 900; j++ {
				if g.cost[i][j] > g.cost[i][k]+g.cost[k][j] {
					g.cost[i][j] = g.cost[i][k] + g.cost[k][j]
					next[i][j] = next[i][k]
				}
			}
		}
	}
}

func routeRestor(start, stop int) []int {
	route := make([]int, 0)
	for cur := start; cur != stop; cur = next[cur][stop] {
		route = append(route, cur)
	}
	route = append(route, stop)
	return route
}

func toMoves(route []int) (move string) {
	log.Println(route)
	for i := 0; i < len(route)-1; i++ {
		switch route[i+1] - route[i] {
		case 1:
			move += "R"
		case 30:
			move += "D"
		case -1:
			move += "L"
		case -30:
			move += "U"
		}
	}
	log.Println(move)
	return
}

/// ------------------------------------------------------
type QueryRecord struct {
	start  Point
	stop   Point
	move   []byte
	result int
}

type Path struct {
	numOfAppeared int
	numOfSelected int
	SampleAverage int
	index         int
}

type PathRecord struct {
	h    [30][30]Path // y,i方向
	v    [30][30]Path // x,j方向
	time int
}

func (pr PathRecord) getPath(i, j int, move byte) Path {
	if move == 'D' || move == 'U' {
		return pr.h[i][j]
	} else {
		return pr.v[i][j]
	}
}

func getIj(now Point, move byte) (int, int) {
	var i, j int
	if move == 'D' || move == 'R' {
		i = now.i
		j = now.j
	} else if move == 'U' || move == 'L' {
		i = now.i + di[direction[move]]
		j = now.j + dj[direction[move]]
	}
	return i, j
}

func (pr *PathRecord) AddAppeared(now Point, move byte) {
	i, j := getIj(now, move)
	if move == 'D' || move == 'U' {
		pr.h[i][j].numOfAppeared++
	} else if move == 'R' || move == 'L' {
		pr.v[i][j].numOfAppeared++
	}
}

func (pr *PathRecord) AddSelected(now Point, move byte) {
	i, j := getIj(now, move)
	if move == 'D' || move == 'U' {
		pr.h[i][j].numOfSelected++
	} else if move == 'R' || move == 'L' {
		pr.v[i][j].numOfSelected++
	}
}

func (pr *PathRecord) AddAverage(now Point, move byte, dis int) {
	i, j := getIj(now, move)
	if move == 'D' || move == 'U' {
		if pr.h[i][j].numOfSelected == 1 {
			pr.h[i][j].SampleAverage = dis
		} else {
			pr.h[i][j].SampleAverage = pr.h[i][j].SampleAverage*(pr.h[i][j].numOfSelected-1) + dis
			pr.h[i][j].SampleAverage = pr.h[i][j].SampleAverage / pr.h[i][j].numOfSelected
		}
	} else if move == 'R' || move == 'L' {
		if pr.v[i][j].numOfSelected == 1 {
			pr.v[i][j].SampleAverage = dis
		} else {
			pr.v[i][j].SampleAverage = pr.v[i][j].SampleAverage*(pr.v[i][j].numOfSelected-1) + dis
			pr.v[i][j].SampleAverage = pr.v[i][j].SampleAverage / pr.v[i][j].numOfSelected
		}
	}
}

func (pr *PathRecord) ReflectResult(q QueryRecord) {
	log.Println(q)
	now := q.start
	average := q.result / len(q.move)
	for i := 0; i < len(q.move); i++ {
		pr.AddAverage(now, q.move[i], average)
		now.move(q.move[i])
	}
}

func solver() {
	var pr PathRecord
	var last QueryRecord
	for i := 0; i < 1000; i++ {
		var q QueryRecord
		q.start.i = nextInt()
		q.start.j = nextInt()
		q.stop.i = nextInt()
		q.stop.j = nextInt()
		// route := query(si, sj, ti, tj)
		//q.move = randomSolver(&q, &pr)
		q.move = greedySolver(&q, &pr)
		fmt.Println(string(q.move))
		q.result = nextInt()
		pr.ReflectResult(q)
	}
	// buildGraph(pr)
	// warchalFloyd()
	// for i := 1000; i < 1000; i++ {
	// 	var q QueryRecord
	// 	q.start.i = nextInt()
	// 	q.start.j = nextInt()
	// 	q.stop.i = nextInt()
	// 	q.stop.j = nextInt()
	// 	s := toindex(q.start.i, q.start.j)
	// 	t := toindex(q.stop.i, q.stop.j)
	// 	log.Println(s, t)
	// 	log.Println(g.cost[s][t], g.cost[t][s])
	// 	path := routeRestor(s, t)
	// 	log.Println(path)
	// 	fmt.Println(toMoves(path))
	// 	q.move = []byte(toMoves(path))
	// 	q.result = nextInt()
	// 	//pr.ReflectResult(q)
	//
	// }
	log.Println(last)
	// buildGraph(pr)
	// warchalFloyd()
	// s := toindex(last.start.i, last.start.j)
	// t := toindex(last.stop.i, last.stop.j)
	// path := routeRestor(s, t)
	// log.Println(path)
	// log.Println(toMoves(path))
}
