package utils

import (
	"strings"
)

func Url_parse(raw []string, urls map[string]bool, ips map[string]bool) bool {
	/*
	 * raw is an array of substrings from the original string that was
	 * splitted by semicolons. The first one [0] is an array of IPs (one or more,
	 * divided with pipe), the third one is an array of URLs (zero(!) or more,
	 * divided by pipe). However, the Split function might return an array of
	 * one substring if it was unable to find a separator in the original
	 * string. Therefore, we must firstly ensure that our array consists
	 * from at least 3 elements.
	 * [ip (| ip2 | ip3 -- optional)] [host] [url (| url2 | url3 optional)]
	 */
	if len(raw) < 3 {
		return false
	}
	/*
	 * We assume here that after splitting the original string into substrings,
	 * we get an array of substrings that explicitly has [0], [1] and [2]
	 * elements. The previous check exits from the procedure otherwise we would get
	 * here and catch a panic. Now we check if there is enough a URL-substring length
	 * to operate on, if not - just return raw[1] which is a domain name of the resource.
	 */
	if len(raw[2]) < 5 {
		urls[raw[1]] = true
		return true
	}
	/*
	 * As Cisco SCE is unable to block/redirect non-http requests (https, for example),
	 * the only decision is to block non-http with ip rules. So here we are checking
	 * whether the URLs substring contains valid plain-http URLs and then including them
	 * to the urls-list, otherwise we should avoid all the URLs in a string and go with IPs.
	 */

	urls_temp := strings.Split(raw[2], " | ")

	for _, _url := range urls_temp {
		if !strings.Contains(_url, "http://") {
			goto Have_non_http
		}
	}
	for _, _url := range urls_temp {
		_url := normalize_url(_url)
		urls[_url] = true
	}
	return true

Have_non_http:
	// TODO: don't trust IPs from list, make internal resolving (awful)
	/* We only get here if some of URLs in the array have non-http scheme */
	for _, v := range strings.Split(raw[0], " | ") {
		ips[v] = true
	}
	return true

	/*Never get here*/
	return false
}

func normalize_url(src string) string {
	//takes string, throws away http:// and escapes ":" in it
	return strings.Replace(src[7:], ":", "\\:", -1)
}
