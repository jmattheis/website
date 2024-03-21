package util

import "strconv"

type Port struct {
	Addr string
	I    int
	S    string
}

var PortOffset = 0

func PortOf(port int) Port {
	if port <= 1000 {
		port += PortOffset
	}

	return Port{
		Addr: ":" + strconv.Itoa(port),
		I:    port,
		S:    strconv.Itoa(port),
	}
}
