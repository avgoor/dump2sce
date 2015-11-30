package utils

import (
	"fmt"
	//	"net"
)

func MakeCiscoACL(ips map[string]bool, aclname string) string {
	// makes cisco-like named ACL with `aclname`
	// that denies all the IPs

	all := fmt.Sprintf("ip access-list extended %s\n", aclname)
	for ip, _ := range ips {
		all += fmt.Sprintf("deny ip %s any\n", ip)
	}
	all += "end\n"
	return all
}

func UploadToCisco(ip string, listing string) bool {
	// uploads listing to the cisco-like router/switch
	// using telnet. returns true on success, otherwise false

	return true
}
