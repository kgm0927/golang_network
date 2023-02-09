package ch04

import (
	"bufio"
	"net"
	"reflect"
	"testing"
)

// #1
const payload = "The bigger the interface, the weaker the abstraction."

func TestScanner(t *testing.T) {

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		c, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer c.Close()
		_, err = c.Write([]byte(payload))
		if err != nil {
			t.Error(err)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	scanner := bufio.NewScanner(conn) // #1
	scanner.Split(bufio.ScanWords)

	var words []string

	for scanner.Scan() { // #2
		words = append(words, scanner.Text() /*#3*/)
	}

	err = scanner.Err()
	if err != nil {
		t.Error(err)
	}

	expected := []string{"The", "bigger", "the", "interface,", "the", "weaker", "the", "abstraction."}

	if !reflect.DeepEqual(words, expected) {
		t.Fatal("inaccurate scanned word list")
	}
	t.Logf("Scanned words: %#v", words)
}
