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

func TestEchoServerUnixDatagram(t *testing.T) {

	dir, err := os.MkdirTemp("", "echo_unixgram")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if rErr := os.RemoveAll(dir); rErr != nil {
			t.Error(rErr)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	/*#1*/ sSockect := filepath.Join(dir, fmt.Sprintf("s%d.sock", os.Getegid()))

	serverAddr, err := datagramEchoServer(ctx, "unixgram", sSockect)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()//#2

	err = os.Chmod(sSockect, os.ModeSocket|0622)//#3
	if err != nil {
		t.Fatal(err)
	}

	/*#1*/
	cSockect := filepath.Join(dir, fmt.Sprintf("c%d.sock", os.Getegid()))
	client, err := net.ListenPacket("unixgram", cSockect)
	if err != nil {
		t.Fatal(err)
	}

	/*#2*/
	defer func() { _ = client.Close() }()

	err = os.Chmod(cSockect, os.ModeSocket|0622) //#3
	if err != nil {
		t.Fatal(err)
	}

	msg := []byte("ping")
	for i := 0; i < 3; i++ { //ping 메시지 3번 읽기
		_, err = client.WriteTo(msg, serverAddr)//#1
		if err != nil {
			t.Fatal(err)
		}
	}

	buf := make([]byte, 1024)
	for i := 0; i < 3; i++ { // "ping" 메시지 3번 읽기
		n, addr, err := client.ReadFrom(buf)//#2
		if err != nil {
			t.Fatal(err)
		}

		if addr.String() != serverAddr.String() {
			t.Fatalf("received reply from %q instead of %q", addr, serverAddr)
		}

		if !bytes.Equal(msg, buf[:n]) {
			t.Fatalf("expected reply %q; actual reply %q", msg, buf[:n])
		}

	}
}
