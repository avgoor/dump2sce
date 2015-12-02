package utils

import (
	"fmt"
	"net"
)

func MakeCiscoACL(ips map[string]bool, aclname string) []string {
	// makes cisco-like named ACL with `aclname`
	// that denies all the IPs
	all := make([]string, 0)
	all = append(all, fmt.Sprintf("ip access-list extended %s\n", aclname))
	for ip, _ := range ips {
		all = append(all, fmt.Sprintf("deny ip %s any\n", ip))
	}
	all = append(all, "end\n")
	return all
}

func UploadToCisco(ip string, listing []string) (bool, error) {
	// uploads listing to the cisco-like router/switch
	// using telnet. returns true on success, otherwise false

	_, err := net.Dial("tcp", ip + ":23")
	if err != nil {
		return false, err
	}

	return true, nil
}
