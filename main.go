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
	os.Exit(realMain())
}

var logIt = log.New(os.Stdout, "DUMPER:", log.Lshortfile|log.Ltime|log.Lmicroseconds|log.Ldate)

func realMain() int {
	urls := make(map[string]bool, 15000)
	ips := make(map[string]bool, 2000)
	filename := utils.Filename
	Config, err := utils.GetCFG(filename)
	if err != nil {
		fmt.Println("Cannot read config from:", filename)
		fmt.Println(err)
		return 1
	}
	logIt.Print(Config)
	if utils.IsProfiling {
		fileprof, err := os.Create("./profile_go")
		if err != nil {
			logIt.Println("Cannot create ./profile_go!")
			return 2
		}

		defer fileprof.Close()
		defer pprof.WriteHeapProfile(fileprof)
	}
	URLFile, err := os.Create(Config.URLfile)
	if err != nil {
		logIt.Println("Cannot write to urls-file!")
		return 8
	}
	logIt.Println("Downloading...")
	x, err := http.Get(Config.ZapretFileURL)
	if (err != nil) || (x.StatusCode != 200) {
		logIt.Println("Cannot download file:", Config.ZapretFileURL)
		return 5
	}
	defer x.Body.Close()
	logIt.Println("Got file")

	outs := bufio.NewScanner(x.Body)       //scanner returns lines one by one
	URLFilefd := bufio.NewWriter(URLFile) //buffered output fast as hell
	logIt.Println("Starting scan")
	for outs.Scan() {
		val := strings.Split(outs.Text(), ";")
		_ = utils.URLParse(val, urls, ips)
	}
	//add manually included urls to the main list
	for _, u := range Config.AdditionalSites {
		urls[u] = true
	}
	for v := range urls {
		URLFilefd.WriteString(v + "\n")
	}

	URLFilefd.Flush()
	URLFile.Close()
	logIt.Println("Scan finished")
	logIt.Println("Uploading URLs to SCE")
	err = utils.UploadToCisco(Config.SCE, Config.SCE.OptionalCMDS)
	if err != nil {
		logIt.Println("Updating SCE failed!")
		logIt.Println(err)
	}
	logIt.Println("SCE update finished")
	logIt.Println("Uploading IPs to Cisco Router")
	err = utils.UploadToCisco(Config.Router, utils.MakeCiscoACL(ips, Config.ACLName))
	if err != nil {
		logIt.Println("Updating Router failed!")
		logIt.Println(err)
	}
	return 0
}
