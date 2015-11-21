package cfgparser

import (
	"encoding/json"
	"flag"
	"os"
	"strconv"
)

var filename = flag.String("file", "./config.json", "Filename and fullpath to the config file")
var IsProfiling = flag.Bool("profile", false, "Turns on profiling")

type CFG struct {
	ZapretFileURL string
	FTPURL        string
	FlavorID      int
}

func init(){
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

func GetFilename() string {
	return *filename
}
