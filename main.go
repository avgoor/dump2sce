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
	"io"
)

//parse an array of string and returns URL
func url_parse(raw []string, urls map[string]bool, ips map[string]bool) bool {
	if len(raw) < 3 { //short strings is a problem
		return false
	}
	//raw format [ip | ip2 | ip3][host][url | url2 | url3 (not always exists)]
	if len(raw[2]) > 5 {
		if raw[2][:5] != "http:" {
			for _, v := range strings.Split(raw[0], " | ") {
				ips[v] = true
			}
			return true
		} else {
			for _, v := range strings.Split(raw[2], " | ") {
				urls[v] = true
			}
			return true
		}
	}
	return false
}

func main() {
	//this will preserve exit code and all the defers same time
	os.Exit(RealMain())
}

var LOG = log.New(os.Stdout, "DUMPER:", log.Lshortfile|log.Ltime|log.Lmicroseconds|log.Ldate)

func defCloser(x io.Closer) {
	LOG.Println("About to close:", x)
	x.Close()
	return
}

func RealMain() int {
	urls := make(map[string]bool, 15000)
	ips := make(map[string]bool, 2000)
	filename := cfgparser.Filename
	Config, err := cfgparser.GetCFG(filename)
	if err != nil {
		fmt.Println("Cannot read config from:", filename)
		return 1
	}
	LOG.Print(Config)
	if cfgparser.IsProfiling {
		fileprof, err := os.Create("./profile_go")
		if err != nil {
			LOG.Println("Cannot create ./profile_go!")
			return 2
		}
		//		pprof.StartCPUProfile(fileprof)
		defer defCloser(fileprof)
		//		defer pprof.StopCPUProfile()
		defer pprof.WriteHeapProfile(fileprof)
	}

	URLFile, err := os.Create(Config.URLfile)
	if err != nil {
		LOG.Println("Cannot write to urls-file!")
		return 8
	}
	defer defCloser(URLFile)

	IPFile, err := os.Create(Config.IPfile)
	if err != nil {
		LOG.Println("Cannot write to ips-file!")
		return 8
	}
	defer defCloser(IPFile)

	LOG.Println("Downloading...")
	x, err := http.Get(Config.ZapretFileURL)
	if (err != nil) || (x.StatusCode != 200) {
		LOG.Println("Cannot download file:", Config.ZapretFileURL)
		return 5
	} else {
		defer defCloser(x.Body)
	}
	LOG.Println("Got file")

	outs := bufio.NewScanner(x.Body)       //scanner returns lines one by one
	URLFile_fd := bufio.NewWriter(URLFile) //buffered output fast as hell
	IPFile_fd := bufio.NewWriter(IPFile)
	LOG.Println("Starting scan")
	for outs.Scan() {
		// short strings contain no data, so omit them
		val := strings.Split(outs.Text(), ";")
		_ = url_parse(val, urls, ips)
	}
	for v, _ := range urls {
		URLFile_fd.WriteString(v + "\n")
	}
	for v, _ := range ips {
		IPFile_fd.WriteString(v + "\n")
	}
	URLFile_fd.Flush()
	IPFile_fd.Flush()
	LOG.Println("Scan finished")
	return 0
}
