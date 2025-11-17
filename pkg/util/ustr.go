package util

import (
	"strconv"
	"strings"
)

func ExtractHostPort(addresses string) map[string]int {
	outer := strings.Split(addresses, ",")
	rMap := make(map[string]int, len(outer))
	for _, out := range outer {
		if len(out) == 0 {
			continue
		}

		inner := strings.Split(out, ":")
		port, _ := strconv.Atoi(inner[1])
		rMap[inner[0]] = port
	}

	return rMap
}

func ExtractAddress(addresses string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, addr := range strings.Split(addresses, ",") {
		addr = strings.TrimSpace(addr)
		if addr == "" {
			continue
		}
		if !seen[addr] {
			seen[addr] = true
			result = append(result, addr)
		}
	}

	return result
}
