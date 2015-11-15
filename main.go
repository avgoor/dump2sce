package main

import (
	"net/http"
	"fmt"
	"bufio"

)

func main() {
	x, _ := http.Get("https://raw.githubusercontent.com/zapret-info/z-i/master/dump.csv")
	outs := bufio.NewScanner(x.Body)
	for outs.Scan() {
		fmt.Printf("->>> %s\n", outs.Text())
	}
}