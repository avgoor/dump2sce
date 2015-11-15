package main
//sce-url-database import cleartext-file ftp://ftp:ftp@46.148.208.22/pub/zi.list flavor-id 300
import (
	"net/http"
	"fmt"
	"bufio"
	"runtime/pprof"
	"strings"
	"os"
)

func url2sce(url string) string {
	return ""
}

func main() {
	fileprof, _ := os.Create("./profile_go")
	pprof.StartCPUProfile(fileprof)
	defer pprof.StopCPUProfile()
	x, _ := http.Get("https://raw.githubusercontent.com/zapret-info/z-i/master/dump.csv")
	outs := bufio.NewScanner(x.Body) //scanner returns lines one by one
	for outs.Scan() {
		// short strings contain no data, so omit them
		if val := strings.Split(outs.Text(), ";"); len(val) > 2 {
			//fmt.Printf(">>%q<<\n", val[2])
		} else {
			fmt.Printf("short: %q\n", val)
		}
	}
}