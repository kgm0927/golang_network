package ch04

import (
	"io"
	"math/rand"
	"net"
	"testing"
)

func TestReadIntoBuffer(t *testing.T) {
	//#1
	payload := make([]byte, 1<<24)
	_, err := rand.Read(payload) // 랜덤한 페이로드 생성

	if err != nil {
		t.Fatal(err)
	}

	Listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := Listener.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		_, err = conn.Write(payload) // #2
		if err != nil {
			t.Error(err)
		}
	}()

	conn, err := net.Dial("tcp", Listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1<<19) // 약 512KB #3

	for {
		n, err := conn.Read(buf)
		//#4

		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			break
		}
		t.Logf("read %d bytes", n) // buf[:n]은 conn 객체에서 읽은 데이터
	}
	conn.Close()
}
