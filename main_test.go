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

func Test_url_parsing_and_normalisation(t *testing.T) {
	_originalHTTP := []string{"127.0.0.1", "example.com", "http://example.com/?newsid=1:1"}
	_expectedHTTP := "example.com/?newsid=1\\:1"
	_originalHTTPS := []string{"127.0.0.2", "example.org", "https://example.org/?newsid=1:1"}
	_expectedHTTPS := "127.0.0.2"
	_invalidString := []string{"tooshort"}

	urls := make(map[string]bool)
	ips := make(map[string]bool)

	if !utils.Url_parse(_originalHTTP, urls, ips) {
		t.Error("error parsing", _originalHTTP)
	}
	if _url, ok := urls[_expectedHTTP]; !ok {
		t.Error("error parsing http, expected", _expectedHTTP, "got", _url)
	}
	if !utils.Url_parse(_originalHTTPS, urls, ips) {
		t.Error("error parsing", _originalHTTPS)
	}
	if _ip, ok := ips[_expectedHTTPS]; !ok {
		t.Error("error parsing https, expected", _expectedHTTPS, "got", _ip)
	}
	if utils.Url_parse(_invalidString, urls, ips) {
		t.Error("Url_parse should fail on", _invalidString, "but it succeed!")
	}
}
