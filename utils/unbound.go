package utils

import "fmt"

func MakeUnboundRules(domains map[string]bool, redirect_ip string) []string {
	var output []string
	for domain, _ := range domains {
		output = append(output, fmt.Sprintf("local-zone: \"%s.\" redirect\n"+
				"local-data: \"%s. IN A %s\"\n", domain, domain, redirect_ip))
	}
	return output
}
