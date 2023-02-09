package ch04

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	count    = flag.Int("c", 3, "number of pings: <=0 means forever")
	interval = flag.Duration("i", time.Second, "Interval between pings")
	timeout  = flag.Duration("W", 5*time.Second, "time to wait for a reply")
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] host:port\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
}
