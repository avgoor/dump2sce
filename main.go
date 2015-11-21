package main

//sce-url-database import cleartext-file ftp://ftp:ftp@46.148.208.22/pub/zi.list flavor-id 300
import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	"./cfgparser"
)

func url2sce(url string) string {
	if url[5:] == "https" {
		return ""
	}
	return ""
}

func main() {
	os.Exit(RealMain())
}

func RealMain() int {
	filename := cfgparser.GetFilename()
	Config, err := cfgparser.GetCFG(filename)
	if err != nil {
		fmt.Println("Cannot read config from:", filename)
		return 1
	}

	LOG := log.New(os.Stdout, "DUMPER:", log.Lshortfile|log.Ltime)
	LOG.Print(Config.String())
	fileprof, _ := os.Create("./profile_go")
	fileout, _ := os.Create("./urls")
	pprof.StartCPUProfile(fileprof)
	defer pprof.StopCPUProfile()
	LOG.Printf("Downloading: %s\n", time.Now())
	x, _ := http.Get(Config.ZapretFileURL)
	LOG.Printf("Getted: %s\n", time.Now())
	outs := bufio.NewScanner(x.Body) //scanner returns lines one by one
	cons := bufio.NewWriter(fileout) //buffered output fast as hell
	LOG.Printf("Before scan: %s\n", time.Now())
	for outs.Scan() {
		// short strings contain no data, so omit them
		if val := strings.Split(outs.Text(), ";"); len(val) > 2 {
			cons.WriteString(strings.Join(val, "!") + "\n")
		} else {
			LOG.Printf("short: %q\n", val)
		}
	}
	cons.Flush()
	fileout.Close()
	LOG.Printf("After scan: %s\n", time.Now())
	return 0
}
