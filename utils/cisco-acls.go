package utils

import (
	"fmt"
	"net"
	"time"
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
	conn, err := net.Dial("tcp4", ip+":23")
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	fmt.Println("Dialed. Reading...", ip)

	p := make([]byte, 32)
	cnt, err := conn.Read(p)
	if err != nil {
		return err
	}
	fmt.Println(p, cnt)
	_, err = conn.Write([]byte{'\xfd', '\x18', '\xff', '\xfa', '\x18', '\x01', '\xfd', '\xf0'})
	if err != nil {
		return err
	}
	cnt, err = conn.Write([]byte("123123123123"))
	if err != nil {
		return err
	}
	fmt.Println(cnt)
	cnt, err = conn.Read(p)
	if err != nil {
		return err
	}
	fmt.Println(p, cnt)
	return nil
}
