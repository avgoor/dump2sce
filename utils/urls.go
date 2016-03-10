package utils

import (
	"net"
	"net/url"
	"strings"

	"golang.org/x/net/idna"
)

// URLParse is an exported function that fills urls/ips maps and return bool
func URLParse(raw []string, urls map[string]bool, ips map[string]bool) bool {

	u_parsed := []*url.URL{}

	/*
	 * raw is an array of substrings from the original string that was
	 * splitted by semicolons. The first one [0] is an array of IPs (one or more,
	 * divided with pipe), the third one is an array of URLs (zero(!) or more,
	 * divided by pipe). However, the Split function might return an array of
	 * one substring if it was unable to find a separator in the original
	 * string. Therefore, we must firstly ensure that our array consists
	 * from at least 3 elements.
	 * [date] [(urls)] [host(s)] [ip(s)]
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

	if len(raw[1]) < 4 {
		goto HaveNonHTTP
	}

	for _, tmp := range strings.Split(raw[1], ",") {
		_url, err := url.Parse(tmp)
		if err != nil {
			// as "" is a valid URL too (sic!)
			// this code does nothing useful
			return false //goto host check
		}
		u_parsed = append(u_parsed, _url)
	}

	for _, _url := range u_parsed {
		/* It is much better to treat a URL having "/" URI
		   or too long URI as a domain-only, along with a pure non-http URL */
		rq := _url.RequestURI()
		not_ok := ((_url.Scheme != "http") || (rq == "/") || (len(rq) > 200))
		not_ok = not_ok || (strings.ContainsAny(rq, ":*"))
		//TODO: check if our URL is already in the non-http database
		if not_ok {
			goto HaveNonHTTP
		}
	}
	/* If we get here this means all the checks above are ok */
	for _, u := range u_parsed {
		host, _, err := net.SplitHostPort(u.Host)
		if err != nil {
			host = u.Host
		}
		_t, _ := idna.ToASCII(host)
		_t = strings.TrimSuffix(_t, ".")
		_t = _t + u.RequestURI()
		urls[_t] = true
	}

	return true

HaveNonHTTP:
	/*
	   We only get here if some of URLs in the array have non-http scheme.
	   So the domain name should be used instead.
	*/
	for _, v := range strings.Split(raw[3], ",") {
		ips[v] = true
	}
	return true

	/*Never get here*/
	return false
}

func normalizeURL(src string) string {
	//takes string, throws away http:// and escapes ":" in it
	return strings.Replace(src[7:], ":", "\\:", -1)
}
