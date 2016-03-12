package utils

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
)

//Filename is exported
var Filename string

//IsProfiling is exported
var IsProfiling bool

type config struct {
	ZapretFileURL   string
	FTPURL          string
	URLfile         string
	DomainsFile     string
	RedirectIP	string
	ACLName         string
	FlavorID        int
	AdditionalSites []string
	SCE             remote
	Router          remote
}

type remote struct {
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

//GetCFG parses json-file referenced by path into CFG
func GetCFG(path string) (config, error) {
	cfgfile, err := os.Open(path)
	if err != nil {
		return config{}, err
	}
	defer cfgfile.Close()
	decoder := json.NewDecoder(cfgfile)
	cfg := config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return config{}, err
	}
	return cfg, nil
}

func (c *config) String() string {
	return "[URL: " + c.ZapretFileURL + " ; FTP: " +
		c.FTPURL + " ; FlavorID: " + strconv.Itoa(c.FlavorID) + "]"
}
