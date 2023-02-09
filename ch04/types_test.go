package ch04

import (
	"bytes"
	"encoding/binary"
	"net"
	"reflect"
	"testing"
)

func TestPayloads(t *testing.T) {

	b1 := Binary("Clear is better than clever") //#1
	b2 := Binary("Don`t panic.")

	s1 := String("Errors are values")    //#2
	payloads := []Payload{&b1, &s1, &b2} //#3

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		for _, p := range payloads {
			_, err = p.WriteTo(conn) //#4
			if err != nil {
				t.Error(err)
				break
			}
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String()) //#1
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	for i := 0; i < len(payloads); i++ {

		actual, err := decode(conn) //#2

		if err != nil {
			t.Fatal(err)
		}

		if expected := payloads[i]; !reflect.DeepEqual(expected, actual) { //#3
			t.Errorf("value mismatch: %v !=%v", expected, actual)
			continue
		}

		t.Logf("[%T] %[1]q", actual) //#4
	}

}

func TestMaxPayloadSize(t *testing.T) {
	buf := new(bytes.Buffer)
	err := buf.WriteByte(BinaryType)
	if err != nil {
		t.Fatal(err)
	}

	err = binary.Write(buf, binary.BigEndian, uint32(1<<30)) // 1기가 #1
	if err != nil {
		t.Fatal(err)
	}

	var b Binary

	_, err = b.ReadFrom(buf)

	if err != ErrMaxPayloadSize { //#2
		t.Fatalf("expected ErrMaxPayloadSize; actual: %v", err)
	}
}
