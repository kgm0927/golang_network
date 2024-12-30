package ch07

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestEchoServerUnix(t *testing.T) {

	dir, err := os.MkdirTemp("", "echo_unix") //#1 원래 ioutil.TempDir()이지만 이제 사용하지 못한다.
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if rErr := os.RemoveAll(dir); rErr != nil { //#2
			t.Error(rErr)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	socket := filepath.Join(dir, fmt.Sprintf("%d.sock", os.Getegid())) //#3
	rAddr, err := streamingEchoServer(ctx, "unix", socket)

	if err != nil {
		t.Fatal(err)
	}

	defer cancel()

	err = os.Chmod(socket, os.ModeSocket|0666) //#4
	if err != nil {
		t.Fatal(err)
	}

	conn, err := net.Dial("unix" /*#1*/, rAddr.String())
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("ping")
	for i := 0; i < 3; i++ { //#2 write 3 "ping" messages
		_, err = conn.Write(msg)
		if err != nil {
			t.Fatal(err)
		}
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf) //#3 read once from the server
	if err != nil {
		t.Fatal(err)
	}

	expected := bytes.Repeat(msg, 3) //#4
	if !bytes.Equal(expected, buf[:n]) {
		t.Fatalf("expected reply %q; actual reply %q", expected,
			buf[:n])
	}

}
