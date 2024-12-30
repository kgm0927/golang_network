package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"

	"github.com/golang_network/ch07/creds/auth"
)

func init() {
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n\t%s <group names>\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func parseGroupNames(args []string) map[string]struct{} {
	groups := make(map[string]struct{})
	for _, arg := range args {
		grp, err := /*#1*/ user.LookupGroup(arg)
		if err != nil {
			log.Println(err)
			continue
		}
		groups[ /*#2*/ grp.Gid] = struct{}{}
	}
	return groups
}

func main() {
	flag.Parse()

	groups := parseGroupNames(flag.Args())
	socket := filepath.Join(os.TempDir(), "creds.sock")
	addr, err := net.ResolveUnixAddr("unix", socket)
	if err != nil {
		log.Fatal(err)
	}

	s, err := net.ListenUnix("unix", addr)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c /*#1*/, os.Interrupt)
	/*#2*/ go func() {
		<-c
		_ = s.Close()
	}()

	fmt.Printf("Listening on %s ...\n", socket)

	for {
		conn, err := /*#3*/ s.AcceptUnix()
		if err != nil {
			break
		}
		if /*#4*/ auth.Allowed(conn, groups) {
			_, err = conn.Write([]byte("Welcome\n"))
			if err == nil {
				// 여기에 연결을 처리하는 고루틴
				continue
			}
		} else {
			_, err = conn.Write([]byte("Access denied\n"))
		}

		if err != nil {
			log.Println(err)
		}

		_ = conn.Close()
	}

}
