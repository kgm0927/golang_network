package ch03

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

func TestDialContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background()) // #1
	sync := make(chan struct{})

	go func() { // #2

		defer func() { sync <- struct{}{} }()

		var d net.Dialer

		d.Control = func(network, address string, c syscall.RawConn) error {
			time.Sleep(time.Second)
			return nil
		}

		conn, err := d.DialContext(ctx, "tcp", "10.0.0.1:80")
		if err != nil {
			t.Log(err)
			return
		}

		conn.Close()
		t.Error("connection did not time out")

	}()

	cancel() // #3
	<-sync   // 종료 알림

	// #4
	if ctx.Err() != context.Canceled {
		t.Errorf("excepted canceled context: actual: %q", ctx.Err())
	}

}

/*
		콘텍스트를 취소하여 연결 중단

콘텍스트를 이용하는 또 다른 장점으로는 cancel 함수 그 자체에 있다. 목록 3-7에서 보는 것처럼 데드라인을 지정하지 않고도
필요 시에 cancel 함수를 이용하여 연결 시도를 취소할 수 있다.


연결 시도를 중단하기 위해 데드라인을 설정해서 콘텍스트를 생성하고 데드라인이 지나가기까지 기다리는 대신 context.WithCancel
함수를 사용하여 콘텍스트와 콘텍스트를 취소할 수 있는 함수를 받는다(#1). 수동으로 콘텍스트를 취소하기 때문에 클로저를 만들어서
별도로 연결 시도를 처리하기 위한 고루틴을 시작한다(#2).

다이얼러가 연결 시도를 하고 원격 노드와 핸드셰이크가 끝나면, 콘텍스트를 취소하기 위해 cancel 함수를 호출한다(#3).

그 결과 DialContext 메서드는 즉시 nil이 아닌 에러를 반환하고 고루틴을 종료한다. 콘텍스트의 Err 메서드를 이용하여 취소된 콘텍스트의
결괏값이 무엇인지 목록 3-6의 데드라인과 비교하기 바란다. 여기서 콘텍스트의 Err 메서드는 context.Canceled를 반환해야 한다(#4).

*/
