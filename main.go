package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"strings"

	"./utils"
)

//parse an array of string and returns URL

func main() {
	//this will preserve exit code and all the defers same time
	os.Exit(realMain())
}

var logIt = log.New(os.Stdout, "DUMPER:", log.Lshortfile|log.Ltime|log.Lmicroseconds|log.Ldate)

func realMain() int {
	urls := make(map[string]bool, 15000)
	domains := make(map[string]bool, 15000)
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
	DomainsFile, err := os.Create(Config.DomainsFile)
	if err != nil {
		logIt.Println("Cannot write to domains-file!")
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

	bodyReader := bufio.NewReaderSize(x.Body, 1024*1024)
	URLFilefd := bufio.NewWriter(URLFile) //buffered output fast as hell
	DomainsFilefd := bufio.NewWriter(DomainsFile)
	logIt.Println("Fetching contents")

	data, err := ioutil.ReadAll(bodyReader)
	strng := strings.Split(string(data), "\n")

	logIt.Println("Start parsing")
	for _, x := range strng {
		val := strings.Split(x, ";")
		_ = utils.URLParse(val, urls, domains)
	}
	//add manually included urls to the main list
	for _, u := range Config.AdditionalSites {
		urls[u] = true
	}
	for v := range urls {
		URLFilefd.WriteString(v + "\n")
	}
	for _, v := range utils.MakeUnboundRules(domains, Config.RedirectIP) {
		DomainsFilefd.WriteString(v)
	}
	URLFilefd.Flush()
	DomainsFilefd.Flush()
	URLFile.Close()
	DomainsFile.Close()
	logIt.Println("Scan finished")
	logIt.Println("Uploading URLs to SCE")
	err = utils.UploadToCisco(Config.SCE, Config.SCE.OptionalCMDS)
	if err != nil {
		logIt.Println("Updating SCE failed!")
		logIt.Println(err)
	}
	logIt.Println("SCE update finished")
	/*	//Turned off so far, as there is no need to make ACLs
		logIt.Println("Uploading IPs to Cisco Router")
		err = utils.UploadToCisco(Config.Router, utils.MakeCiscoACL(domains, Config.ACLName))
		if err != nil {
			logIt.Println("Updating Router failed!")
			logIt.Println(err)
		}
	*/
	return 0
}
