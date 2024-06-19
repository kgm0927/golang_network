package ch03

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

func TestDialContext(t *testing.T) {

	dl := time.Now().Add(5 * time.Second)                         // #1
	ctx, cancel := context.WithDeadline(context.Background(), dl) //context.WithCancel과 다르다. //#2
	defer cancel()                                                //#3

	var d net.Dialer

	d.Control = func(network, address string, c syscall.RawConn) error { //#4
		time.Sleep(5*time.Second + time.Millisecond)
		return nil
	}

	conn, err := d.DialContext(ctx, "tcp", "10.0.0.0:80")
	//#5
	if err == nil {
		conn.Close()
		t.Fatal("connection did not time out")
	}

	nErr, ok := err.(net.Error)

	if !ok {
		t.Error(err)
	} else {
		if !nErr.Timeout() {
			t.Errorf("error is not a timeout: %v", err)
		}
	}

	if ctx.Err() != context.DeadlineExceeded { //#6
		t.Errorf("expected deadline exceeded; actual: %v", ctx.Err())
	}

}

/*
				데드라인 콘텍스트를 이용하여 연결을 타임아웃 하기

표준 라이브러리의 context 패키지를 이용하여 콘텍스트를 이용하면 더욱 현대적인 방법으로 연결 시도를
타임아웃할 수 있다. 콘텍스트(context)란 비동기 프로세스에 취소 시그널을 보낼 수 있는 객체이다. 또한,
데드라인이 지나거나 타이머가 만료된 이후에도 취소 시그널을 보낼 수도 있다.


콘텍스트를 취소하기 위해 각 콘텍스트마다 초기화 시에 반환되는 cancel함수를 사용한다. cancel 함수를 사용하면
콘텍스트가 데드라인이 지나기도 전에 함수를 의도적으로 취소할 수 있는 유연함을 제공한다. cancel 함수를 코드상의 다른
부분으로 제어권을 넘길 수도 있다.





연결을 시도하기 이전에 먼저, 5초 후 데드라인이 지나는 콘텍스트를 만들기 위해 현재 시간으로부터 5초 뒤의 시간을 저장한다(#1).
이후 WithDeadline 함수를 사용하여 콘텍스트와 cancel 함수를 생성하고 위에서 생성한 데드라인을 설정한다(#2). 가능한 한 바로 콘텍스트가
가비지 컬렉션이 되도록 cancel 함수를 defer 호출하는 것이 좋다(#3).

그 후 다이얼러의 Control 함수(#4)를 오버라이딩하여 연결을 콘텍스트의 데드라인을 간신히 초과하는 정도로 지연시킨다(#5). 테스트의 끝 부분의
에러 처리(#6)는 데드라인이 콘텍스트를 제대로 취소하였는지 cancel 함수에는 문제가 없었는지 확인한다.


Dialtimeout에 관해서 여러 IP 주소로 해석되는 호스트에 다이얼 시 Go에서는 각 IP 주소 중 먼저 연결되는 주소를 기본 IP 주소로 연결 시도한다.
첫 번째로 성공한 연결 시도에 대한 연결만 유지되고 그 외의 모든 연결 시도는 취소된다. 모든 연결 시도가 실패하거나 콘텍스트의 데드라인이 지나면
net.Dialer.DialContext 에러를 반환한다.




*/
