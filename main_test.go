package main_test

import (
	"./utils"
	"os"
	"testing"
)

func TestConfigParser(t *testing.T) {
	FakeJSON := []byte(`{
  "ZapretFileURL": "https://raw.githubusercontent.com/zapret-info/z-i/master/dump.csv",
  "FTPURL": "ftp://fake.ftp/fake.list",
  "FlavorID": 31,
  "URLFile": "./zapret.url",
  "IPFile": "./zapret.ips",
  "ACLName": "zapret"
}`)
	filename := "_fake_config.json"
	f, _ := os.Create(filename)
	defer func() {
		f.Close()
		os.Remove(filename)
	}()
	f.Write(FakeJSON)
	f.Sync()
	FakeCFG, _ := utils.GetCFG(filename)
	if FakeCFG.FTPURL != "ftp://fake.ftp/fake.list" {
		t.Error("Expected: ftp://fake.ftp/fake.list, got", FakeCFG.FTPURL)
	}
}
