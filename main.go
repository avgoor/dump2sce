package main

import (
	"net/http"
	"fmt"
	"bufio"
	"runtime/pprof"
	"strings"
	"os"
)

func main() {
	fileprof, _ := os.Create("./profile_go")
	pprof.StartCPUProfile(fileprof)
	defer pprof.StopCPUProfile()
	x, _ := http.Get("https://raw.githubusercontent.com/zapret-info/z-i/master/dump.csv")
	outs := bufio.NewScanner(x.Body)
	for outs.Scan() {
		if val := strings.Split(outs.Text(), ";"); len(val) > 2 {
			//fmt.Printf(">>%q<<\n", val[2])
		} else {
			fmt.Printf("short: %q\n", val)
		}
	}
}