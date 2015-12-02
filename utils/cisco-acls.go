package utils

import (
	"fmt"
	"net"
	//	"bufio"
)

func MakeCiscoACL(ips map[string]bool, aclname string) []string {
	// makes cisco-like named ACL with `aclname`
	// that denies all the IPs. returns slice of strings
	all := make([]string, 0)
	all = append(all, fmt.Sprintf("ip access-list extended %s\n", aclname))
	for ip, _ := range ips {
		all = append(all, fmt.Sprintf("deny ip %s any\n", ip))
	}
	all = append(all, "end\n")
	return all
}

func UploadToCisco(ip string, listing []string) error {
	// uploads listing to the cisco-like router/switch
	// using telnet. returns true on success, otherwise false
	fmt.Println("Reaching", ip)
	conn, err := net.Dial("tcp4", "127.0.0.1:23")
	if err != nil {
		return err
	}
	fmt.Println("Dialed. Reading...", ip)
	strn := []byte{}
	for {
		count, err := conn.Read(strn)
		if err != nil {
			return err
		}
		fmt.Println(strn, count)
	}

	return nil
}
