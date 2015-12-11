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
	_original_http := []string{"127.0.0.1", "example.com", "http://example.com/?newsid=1:1"}
	_expected_http := "example.com/?newsid=1\\:1"
	_original_https := []string{"127.0.0.2", "example.org", "https://example.org/?newsid=1:1"}
	_expected_https := "127.0.0.2"
	_invalid_string := []string{"tooshort"}

	urls := make(map[string]bool)
	ips := make(map[string]bool)

	if !utils.Url_parse(_original_http, urls, ips) {
		t.Error("error parsing", _original_http)
	}
	if _url, ok := urls[_expected_http]; !ok {
		t.Error("error parsing http, expected", _expected_http, "got", _url)
	}
	if !utils.Url_parse(_original_https, urls, ips) {
		t.Error("error parsing", _original_https)
	}
	if _ip, ok := ips[_expected_https]; !ok {
		t.Error("error parsing https, expected", _expected_https, "got", _ip)
	}
	if utils.Url_parse(_invalid_string, urls, ips) {
		t.Error("Url_parse should fail on", _invalid_string, "but it succeed!")
	}
}
