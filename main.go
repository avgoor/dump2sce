package main

//sce-url-database import cleartext-file ftp://ftp:ftp@46.148.208.22/pub/zi.list flavor-id 300
import (
	"./utils"
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"runtime/pprof"
)

//parse an array of string and returns URL

func main() {
	//this will preserve exit code and all the defers same time
	os.Exit(RealMain())
}

var LOG = log.New(os.Stdout, "DUMPER:", log.Lshortfile|log.Ltime|log.Lmicroseconds|log.Ldate)

func defCloser(x io.Closer) {
	x.Close()
	return
}

func RealMain() int {
	urls := make(map[string]bool, 15000)
	ips := make(map[string]bool, 2000)
	filename := utils.Filename
	Config, err := utils.GetCFG(filename)
	if err != nil {
		fmt.Println("Cannot read config from:", filename)
		return 1
	}
	LOG.Print(Config)
	if utils.IsProfiling {
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
		_ = utils.Url_parse(val, urls, ips)
	}
	for v, _ := range urls {
		URLFile_fd.WriteString(v + "\n")
	}

	for _, rule := range utils.MakeCiscoACL(ips, Config.ACLName) {
		IPFile_fd.WriteString(rule)
	}

	URLFile_fd.Flush()
	IPFile_fd.Flush()
	LOG.Println("Scan finished")
	return 0
}
