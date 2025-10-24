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
