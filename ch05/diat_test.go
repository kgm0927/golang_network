package ch05

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"
)

func TestDailUDP(t *testing.T) {
	// 에코 서버와 클라이언트 생성하기
	ctx, cancel := context.WithCancel(context.Background())

	/*#1*/
	serverAddr, err := echoServerUDP(ctx, "127.0.0.1:")
	if err != nil {
		t.Fatal()
	}
	defer cancel()

	client, err := /*#2*/ net.Dial("udp", serverAddr.String())

	if err != nil {
		t.Fatal(err)
	}

	defer func() { _ = client.Close() }()

	// 클라이언트 인터럽트

	interloper, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	interrupt := []byte("pardon me")

	/*#1*/
	n, err := interloper.WriteTo(interrupt, client.LocalAddr())
	if err != nil {
		t.Fatal(err)
	}

	_ = interloper.Close()

	if l := len(interrupt); l != n {
		t.Fatalf("wrote %d bytes of %d", n, l)
	}

	ping := []byte("ping")
	_, err = client.Write(ping) //#1
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1024)

	n, err = client.Read(buf) //#2
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(ping, buf[:n]) {
		t.Errorf("expected reply %q; actual reply %q", ping, buf[:n])
	}

	err = client.SetDeadline(time.Now().Add(time.Second)) //#3

	_, err = client.Read(buf) //#4
	if err != nil {
		t.Fatal("unexpected packet")
	}
}
