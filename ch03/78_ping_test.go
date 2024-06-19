package ch03

import (
	"context"
	"io"
	"net"
	"testing"
	"time"
)

func TestPingerAdvanceDeadline(t *testing.T) {

	done := make(chan struct{})
	listener, err := net.Listen("tcp", "127.0.0.1:")

	if err != nil {
		t.Fatal(err)
	}

	begin := time.Now()
	go func() {
		defer func() { close(done) }()

		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())

		defer func() {
			cancel()
			conn.Close()
		}()

		resetTimer := make(chan time.Duration, 1)
		resetTimer <- time.Second
		go Pinger(ctx, conn, resetTimer)

		err = conn.SetDeadline(time.Now().Add(5 * time.Second)) // #1

		if err != nil {
			t.Error(err)
			return
		}

		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf) // 수신

			if err != nil {
				return
			}

			t.Logf("[%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])

			resetTimer <- 0                                         // #2
			err = conn.SetDeadline(time.Now().Add(5 * time.Second)) // #3

			if err != nil {
				t.Error(err)
				return
			}

		}

	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)

	for i := 0; i < 4; i++ { // 핑을 4개 읽음 // #4
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("[%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])
	}

	_, err = conn.Write([]byte("PONG!!!")) // 핑 타이머를 초기화함. // #5

	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 4; i++ { // 핑을 네 개더 읽음	// #6
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Fatal(err)
			}
			break
		}
		t.Logf("[%s] %s", time.Since(begin).Truncate(time.Second), buf[:n])
	}

	<-done
	end := time.Since(begin).Truncate(time.Second)

	t.Logf("[%s] done", end)
	if end != 9*time.Second { // #7
		t.Fatalf("excepted EOF at 9 seconds; actual %s", end)
	}
}

/*
		하트비트를 이용하여 데드라인 늦추기

양측의 네트워크 연결에서 Pinger를 이용하여 상대 노드가 유휴 상태가 되었을 대 데드라인을 늦을 수 있다. 지금까지는 한쪽에서만 Pinger를 사용하는 예시를 살펴보았다.
네트워크 연결 중 한쪽에서 데이터를 정상적으로 수신하면 불필요한 핑 전송을 막기 위해 핑 타이머는 리셋되어야 한다.



네트워크 연결을 수신하는 리스너를 시작하고, 초마다 핑을 전송하는 Pinger를 시작한다. 그리고 데드라인의 초깃겞으로 5초를
설정한다(#1). 클라이언트 관점에서 서버가 데드라인이 지나고 서버 측의 연결을 끊기까지 io.EOF를 마지막으로 총 네 번의 핑의
수신할 수 있는 시간이다.

하지만 서버의 데드라인이 지나기 전에 클라이언트가 서버로 데이터로 계속 송신함으로써(#5) 데드라인을 늦출 수 있다.


서버가 데이터를 수신할 수 있으면 아직 네트워크 연결이 정상적이라는 것을 알 수 있다. 그러므로 Pinger를 초기화하고(#2) 연결의 데드라인을 늦춘다(#3). 소켓이 종료되는
것을 방지하기 위해서는 클라이언트는 서버로부터 네 개의 핑 메시지를 수신하고(#4) 퐁 메시지를 송신한다(#5).

이로 인해 서버가 데드라인이 지나기 전까지 5초를 확보할 수 있다. 이후 클라이언트는 네 개의 핑 메시지를 수신하고(#6) 데드라인이 지나가기를 기다린다. 서버가 연결을 끝낸
시점에서 총 9초(#7)를 기다렸다는 것을 확인할 수 있다.

*/
