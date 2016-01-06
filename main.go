package main

import (
	"./utils"
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"strings"
)

//parse an array of string and returns URL

func main() {
	//this will preserve exit code and all the defers same time
	os.Exit(RealMain())
}

var LOG = log.New(os.Stdout, "DUMPER:", log.Lshortfile|log.Ltime|log.Lmicroseconds|log.Ldate)

func RealMain() int {
	urls := make(map[string]bool, 15000)
	ips := make(map[string]bool, 2000)
	filename := utils.Filename
	Config, err := utils.GetCFG(filename)
	if err != nil {
		fmt.Println("Cannot read config from:", filename)
		fmt.Println(err)
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
		defer fileprof.Close()
		//		defer pprof.StopCPUProfile()
		defer pprof.WriteHeapProfile(fileprof)
	}
	URLFile, err := os.Create(Config.URLfile)
	if err != nil {
		LOG.Println("Cannot write to urls-file!")
		return 8
	}
	IPFile, err := os.Create(Config.IPfile)
	if err != nil {
		LOG.Println("Cannot write to ips-file!")
		return 8
	}
	LOG.Println("Downloading...")
	x, err := http.Get(Config.ZapretFileURL)
	if (err != nil) || (x.StatusCode != 200) {
		LOG.Println("Cannot download file:", Config.ZapretFileURL)
		return 5
	} else {
		defer x.Body.Close()
	}
	LOG.Println("Got file")

	outs := bufio.NewScanner(x.Body)       //scanner returns lines one by one
	URLFile_fd := bufio.NewWriter(URLFile) //buffered output fast as hell
	IPFile_fd := bufio.NewWriter(IPFile)
	LOG.Println("Starting scan")
	for outs.Scan() {
		val := strings.Split(outs.Text(), ";")
		_ = utils.Url_parse(val, urls, ips)
	}
	//add manually included urls to the main list
	for _, u := range Config.AdditionalSites {
		urls[u] = true
	}
	for v, _ := range urls {
		URLFile_fd.WriteString(v + "\n")
	}
	for _, rule := range utils.MakeCiscoACL(ips, Config.ACLName) {
		IPFile_fd.WriteString(rule)
	}
	URLFile_fd.Flush()
	IPFile_fd.Flush()
	URLFile.Close()
	IPFile.Close()
	LOG.Println("Scan finished")
	err = utils.UploadToCisco(Config.SCE, Config.SCE.OptionalCMDS)
	if err != nil {
		LOG.Println(err)
	}
	return 0
}
