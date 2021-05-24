package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile)
	P, _ = os.Getwd()
	//testRun()
	parallelRun()
}

var P string = ""

func testn(n int) {
	sumScore := 0
	for i := 0; i < n; i++ {
		fmt.Print("case=", i)
		score, _ := run(i)
		fmt.Printf(" score=%d \n", score)
		sumScore += score
	}
	fmt.Println("ALL SCORE = ", sumScore)
}

func testRun() {
	score, n := run(0)
	log.Printf("score=%d loop=%d\n", score, n)
	// vscore := vis(inputpaths[50], out)
}

// 	./tools/target/release/tester tools/in/0000.txt ./solver > out.txt
func run(seed int) (int, int) {
	tester := P + "/tools/target/release/tester"
	solver := P + "/solver"
	inFile := fmt.Sprintf("%s/tools/in/%s.txt", P, fmt.Sprintf("%04d", seed))
	outFile := fmt.Sprintf("%s/out/%s.out", P, fmt.Sprintf("%04d", seed))
	cmdStr := tester + " " + inFile + " " + solver + " > " + outFile
	//cmdStr := exe + " < " + inFile + " > " + outFile
	cmds := []string{"sh", "-c", cmdStr}
	cmd := exec.Command(cmds[0], cmds[1:]...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		log.Println(cmds)
		log.Fatal(err)
	}
	cmd.Wait()
	score := parseScore(stderr.String())
	if score == 0 {
		log.Println(stderr.String())
	}
	//loop := parseLoop(stderr.String())
	return score, 0
}

type Date struct {
	seed  int
	score int
	time  int
}

func parallelRun() {
	CORE := 4
	maxSeed := 100
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, CORE-1)
	datas := make([]Date, 0)
	sumScore := 0
	for seed := 0; seed < maxSeed; seed++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(seed int) {
			startTime := time.Now()
			var d Date
			d.score, d.time = run(seed)
			elapsed := time.Since(startTime)
			d.seed = seed
			mu.Lock()
			datas = append(datas, d)
			// fmt.Print(".")
			fmt.Printf("seed=%d score=%d time=%vms\n", d.seed, d.score, elapsed)
			sumScore += d.score
			mu.Unlock()
			wg.Done()
			<-sem
		}(seed)
	}
	fmt.Printf("SCORE=%d\n", sumScore)
}

func parseScore(s string) int {
	ms := `Score = ([0-9]+)`
	re := regexp.MustCompile(ms)
	ma := re.FindString(s)
	score, err := strconv.Atoi(strings.Replace(ma, "Score = ", "", -1))
	if err != nil {
		log.Println(score)
		log.Println(ma)
	}
	return score
}

func parseTime(s string) int {
	ms := `time=([0-9]+)`
	re := regexp.MustCompile(ms)
	ma := re.FindString(s)
	n, err := strconv.Atoi(strings.Replace(ma, "loop=", "", -1))
	if err != nil {
		log.Println(n)
	}
	return n
}

func vis(input string, output string) (score int) {
	vispath := P + "/tools/target/release/vis"
	cmdStr := vispath + " " + input + " " + output
	cmds := []string{"sh", "-c", cmdStr}
	var out []byte
	var err error
	out, err = exec.Command(cmds[0], cmds[1:]...).Output()
	if err != nil {
		log.Fatal(err)
	}
	outs := strings.Split(string(out), "\n")
	score, err = strconv.Atoi(outs[0])
	if err != nil {
		panic(err)
	}
	return score
}
