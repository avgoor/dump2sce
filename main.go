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
	LOG := log.New(os.Stdout, "DUMPER:", log.Lshortfile|log.Ltime|log.Lmicroseconds)
	filename := cfgparser.GetFilename()
	Config, err := cfgparser.GetCFG(filename)
	if err != nil {
		fmt.Println("Cannot read config from:", filename)
		return 1
	}
	LOG.Print(Config.String())
	if *cfgparser.IsProfiling { //flag returns pointers
		fileprof, err := os.Create("./profile_go")
		if err != nil {
			LOG.Println("Cannot create ./profile_go!")
			return 2
		}
		pprof.StartCPUProfile(fileprof)
		defer fileprof.Close()
		defer pprof.StopCPUProfile()
	}
	fileout, _ := os.Create("./urls")
	LOG.Println("Downloading...")
	x, _ := http.Get(Config.ZapretFileURL)
	LOG.Println("Got file")
	outs := bufio.NewScanner(x.Body) //scanner returns lines one by one
	cons := bufio.NewWriter(fileout) //buffered output fast as hell
	LOG.Println("Starting scan")
	for outs.Scan() {
		// short strings contain no data, so omit them
		if val := strings.Split(outs.Text(), ";"); len(val) > 2 {
			cons.WriteString(strings.Join(val, "!") + "\n")
		} else {
			LOG.Printf("Invalid string. Too short: %q\n", val)
		}
	}
	cons.Flush()
	fileout.Close()
	LOG.Println("Scan finished")
	return 0
}
