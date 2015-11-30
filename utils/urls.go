package utils

import "strings"

func Url_parse(raw []string, urls map[string]bool, ips map[string]bool) bool {
	if len(raw) < 3 { //short strings is a problem
		return false
	}
	//raw format [ip | ip2 | ip3][host][url | url2 | url3 (not always exists)]
	if len(raw[2]) > 5 {
		if raw[2][:5] != "http:" {
			for _, v := range strings.Split(raw[0], " | ") {
				ips[v] = true
			}
			return true
		} else {
			for _, v := range strings.Split(raw[2], " | ") {
				v := normalize_url(v)
				urls[v] = true
			}
			return true
		}
	}
	return false
}

func normalize_url(src string) string {
	//takes string and escapes ":" in it
	return src[:5] + strings.Replace(src[5:], ":", "\\:", -1)
}
