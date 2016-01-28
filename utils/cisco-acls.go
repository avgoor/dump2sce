package utils

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"
)
// MakeCiscoACL returns a slice of strings that contains
// rules in cisco's format
func MakeCiscoACL(ips map[string]bool, aclname string) []string {
	// makes cisco-like named ACL with `aclname`
	// that denies all the IPs. returns slice of strings
	var all []string
	all = append(all, fmt.Sprintf("conf t\n"))
	// at first delete this acl, as we don't want to keep old IPs there
	// and don't want to make some change tracking algo, which can precisely
	// remove or insert rule. (it's a way more complicated task)
	all = append(all, fmt.Sprintf("no ip access-list extended %s\n", aclname))
	all = append(all, fmt.Sprintf("ip access-list extended %s\n", aclname))
	for ip := range ips {
		all = append(all, fmt.Sprintf("deny ip host %s any\n", ip))
	}
	// permit any other traffic, otherwise it will be a bad day for sysadmin
	all = append(all, "permit ip any any\n")
	all = append(all, "exit\n")
	return all
}

func readUntil(r io.Reader, bu []byte) ([]byte, error) {

	br := []byte{}
	b := make([]byte, 1)
	for {
		_, e := r.Read(b)
		if e != nil {
			return br, e
		}
		br = append(br, b[0])
		if bytes.HasSuffix(br, bu) {
			return br, nil
		}
	}
}
// UploadToCisco is a naive telnet implementation that logs into
// a device, sends arbitrary commands (like acl) and logs out.
func UploadToCisco(cfg remote, listing []string) error {
	// TODO: make a new, options-aware telnet handling
	// uploads listing to the cisco-like router/switch
	// using telnet. returns true on success, otherwise false
	fmt.Println("Reaching", cfg)
	conn, err := net.Dial("tcp4", cfg.IP+":23")
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(time.Duration(cfg.Timeout) * time.Second))
	fmt.Println("Dialed. Reading...", cfg.IP)

	/*
		Interaction structure:

		1) Log in (send login & pw)
		2) Check whether we logged in
		3) command loop (send commands checking reaction)
		4) exit & cleanup

	*/

	events := []struct {
		expect string
		send   string
	}{
		{"ame:", cfg.Login},
		{"ord:", cfg.Password},
		{">", "enable"},
		{"ord:", cfg.EnablePW},
		{"#", ""},
	}

	// login part
	for _, event := range events {
		t, e := readUntil(conn, []byte(event.expect))
		fmt.Printf("%s", t)
		if e != nil {
			return e
		}
		_, e = conn.Write([]byte(event.send + "\n"))
		if e != nil {
			return e
		}
	}
	// cmd loop
	for _, cmd := range listing {
		_, e := conn.Write([]byte(cmd))
		if e != nil {
			return e
		}
		t, e := readUntil(conn, []byte(")#"))
		fmt.Printf("%s", t)
		if e != nil {
			return e
		}
	}
	// cleanup
	conn.Write([]byte("end\n"))
	t, _ := readUntil(conn, []byte("#"))
	fmt.Printf("%s", t)
	conn.Write([]byte("exit\n"))
	t, _ = readUntil(conn, []byte("\n"))
	fmt.Printf("%s", t)
	conn.Close()
	return nil
}
