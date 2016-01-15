package utils

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
)

var Filename string
var IsProfiling bool

type CFG struct {
	ZapretFileURL   string
	FTPURL          string
	URLfile         string
	ACLName         string
	FlavorID        int
	AdditionalSites []string
	SCE             Remote
	Router          Remote
}

type Remote struct {
	IP           string
	Login        string
	Password     string
	EnablePW     string
	Timeout      int
	OptionalCMDS []string
}

func init() {
	flag.StringVar(&Filename, "file", "./config.json", "Filename and fullpath to the config file")
	flag.BoolVar(&IsProfiling, "profile", false, "Turns on profiling")
	flag.Parse()
}

//parses json-file referenced by path into CFG

func GetCFG(path string) (CFG, error) {
	cfgfile, err := os.Open(path)
	if err != nil {
		return CFG{}, err
	}
	defer cfgfile.Close()
	decoder := json.NewDecoder(cfgfile)
	cfg := CFG{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return CFG{}, err
	}
	return cfg, nil
}

func (c *CFG) String() string {
	return "[URL: " + c.ZapretFileURL + " ; FTP: " +
		c.FTPURL + " ; FlavorID: " + strconv.Itoa(c.FlavorID) + "]"
}
