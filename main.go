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
	if url[5:] == "https" {
		return ""
	}
	return ""
}

func main() {
	fileprof, _ := os.Create("./profile_go")
	fileout, _ := os.Create("./urls")
	pprof.StartCPUProfile(fileprof)
	defer pprof.StopCPUProfile()
	x, _ := http.Get("https://raw.githubusercontent.com/zapret-info/z-i/master/dump.csv")
	outs := bufio.NewScanner(x.Body) //scanner returns lines one by one
	cons := bufio.NewWriter(fileout)//buffered output fast as hell
	for outs.Scan() {
		// short strings contain no data, so omit them
		if val := strings.Split(outs.Text(), ";"); len(val) > 2 {
//			fmt.Fprintln(cons, val[2])
			fileout.Write([]byte (val[2] + "\n"))
		} else {
			fmt.Printf("short: %q\n", val)
		}
	}
	cons.Flush()
	fileout.Close()
}