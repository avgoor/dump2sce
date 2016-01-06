package utils

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"
)

func MakeCiscoACL(ips map[string]bool, aclname string) []string {
	// makes cisco-like named ACL with `aclname`
	// that denies all the IPs. returns slice of strings
	all := make([]string, 0)
	all = append(all, fmt.Sprintf("conf t\n"))
	all = append(all, fmt.Sprintf("ip access-list extended %s\n", aclname))
	for ip, _ := range ips {
		all = append(all, fmt.Sprintf("deny ip %s any\n", ip))
	}
	all = append(all, "end\n")
	all = append(all, "end\n")
	return all
}

func read_until(r io.Reader, bu []byte) ([]byte, error) {

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

func UploadToCisco(cfg Remote, listing []string) error {
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
		t, e := read_until(conn, []byte(event.expect))
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
		t, e := read_until(conn, []byte(")#"))
		fmt.Printf("%s", t)
		if e != nil {
			return e
		}
	}
	// cleanup
	conn.Write([]byte("end\n"))
	t, _ := read_until(conn, []byte("#"))
	fmt.Printf("%s", t)
	conn.Write([]byte("exit\n"))
	t, _ = read_until(conn, []byte("\n"))
	fmt.Printf("%s", t)
	conn.Close()
	return nil
}
