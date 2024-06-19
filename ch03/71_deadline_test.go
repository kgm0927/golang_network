package ch03

import (
	"io"
	"net"
	"testing"
	"time"
)

func TestDeadline(t *testing.T) {
	sync := make(chan struct{})

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := listener.Accept()

		if err != nil {
			t.Log(err)
			return
		}

		defer func() {
			conn.Close()
			close(sync)
		}()

		err = conn.SetDeadline(time.Now().Add(5 * time.Second)) // #1

		if err != nil {
			t.Error(err)
			return
		}
		buf := make([]byte, 1)
		_, err = conn.Read(buf)

		nErr, ok := err.(net.Error)

		if !ok || !nErr.Timeout() { // #2
			t.Errorf("excepted timeout error; actual: %v", err)
		}

		sync <- struct{}{}

		err = conn.SetDeadline(time.Now().Add(5 * time.Second)) // #3

		if err != nil {
			t.Error(err)
			return
		}

		_, err = conn.Read(buf)

		if err != nil {
			t.Error(err)
		}

	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	<-sync
	_, err = conn.Write([]byte("1"))

	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 1)
	_, err = conn.Read(buf)
	if err != io.EOF { // #4
		t.Errorf("expected server termination; actual: %v", err)
	}

}

/*

		데드라인 구현하기

Go의 네트워크 연결 객체는 읽기와 쓰기 동작에 대해 모두 데드라인(deadline)을 포함한다. 데드라인은 아무런 패킷도 오고
가지 않은 채로 네트워크 연결이 얼마나 유휴 상태로 지속할 수 있는지를 제어한다.


Read 메서드에 대한 데드라인은 연결 객체 내의 SetReadDeadline 메서드를 사용하고 제어하며, Write 메서드에 대한 데드라인은
SetWriteDeadline 메서드를 사용하여 제어하고, 또는 Read와 Write의 데드라인을 동시에 SetDeadline 메서드를 사용하여 제어한다.

연결상의 읽기 데드라인이 지나게 되면 현재 블로킹되어 있는 동작과 앞으로의 네트워크 연결상의 Read 메서드는 곧바로 타임아웃
에러를 반환한다. Write메서드 또한 그러하다.


Go의 네트워크 연결은 기본적으로 읽기와 쓰기 동작에 대해 데드라인을 설정하지 않는다. 즉, 네트워크 연결이 유휴 상태로 끊기지
안하고 아주 오랜 시간 존재할 수 있다.

서버가 클라이언트의 TCP 연결을 수락하고 나서 연결상의 읽기에 대한 데드라인을 설정한다(#1).


테스트상에서 클라이언트가 데이터를 전송하지 않기 때문에 Read 메서드는 데드라인이 지날 때까지 블로킹 될 것이다. Read 메서드는
타임아웃으로 설정한(#2) 5초 뒤에 에러를 반환한다.

이후에 발생하는 모든 읽기 시도는 즉시 또 다른 타임아웃 에러를 반환할 것이다. 한편 데드라인을 좀 더 뒤로 설정하여(#3) 다시 읽기가
정상적으로 동작하게 할 수도 있다.

이후에는 Read 메서드가 성공한다. 서버는 네트워크 연결을 종료하며, 클라이언트와의 연결 종료 절차를 시작한다. 현재 Read 메서드에서
블로킹된 클라이언트는 네트워크 연결이 종료됨에 따라 io.EOF를 반환받는다(#4).



정해진 시간 동안 원격 노드에서 아무런 데이터 송신이 없는 경우 원격 노드와의 연결이 끊어졌는데 FIN 패킷을 받지 못한 경우이거나, 유휴
상태로 판단할 수 있다.
*/
