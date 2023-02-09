package ch05

import (
	"bytes"
	"context"
	"net"
	"testing"
)

func TestListenPacketUDP(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	serverAdder, err := echoServerUDP(ctx, "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	defer cancel()

	client, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = client.Close() }()

	// 클라이언트와 에코 서버강의 메시지를 전송하여 인터럽트하기.

	/*#1*/
	interloper, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	interrupt := []byte("pardon me")

	/*#2*/
	n, err := interloper.WriteTo(interrupt, client.LocalAddr())
	if err != nil {
		t.Fatal(err)
	}

	_ = interloper.Close()

	if l := len(interrupt); l != n {
		t.Fatalf("wrote %d bytes of %d", n, l)
	}

	// 여러 송신자로부터 한번에 UDP 패킷을 수신하기

	ping := []byte("ping")
	_, err = client.WriteTo(ping, serverAdder) //#1
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1024)
	n, addr, err := client.ReadFrom(buf) //#2
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal( /*#3*/ interrupt, buf[:n]) {
		t.Errorf("expected reply %q: actual reply %q", interrupt, buf[:n])
	}

	if addr.String() != interloper.LocalAddr().String() {
		t.Errorf("expected message from %q; actual sender is %q", interloper.LocalAddr(), addr)
	}

	n, addr, err = client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal( /*#4*/ ping, buf[:n]) {
		t.Errorf("excepted reply %q; actual reply %q", ping, buf[:n])
	}

	if addr.String() != serverAdder.String() { //#5
		t.Errorf("excepted message from %q; actual sender is %q", serverAdder, addr)
	}

}
