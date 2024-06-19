package ch03

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"
)

func TestDialContextCancelFanout(t *testing.T) {

	// #1
	ctx, cancel := context.WithDeadline(
		context.Background(), time.Now().Add(10*time.Second),
	)

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	// #2
	go func() {
		// 하나의 연결만을 수락합니다.
		conn, err := listener.Accept()
		if err == nil {
			conn.Close()
		}
	}()

	// #3
	dial := func(ctx context.Context, address string, response chan int, id int, wg *sync.WaitGroup) {

		defer wg.Done()

		var d net.Dialer
		c, err := d.DialContext(ctx, "tcp", address)

		if err != nil {
			return
		}
		c.Close()

		select {

		case <-ctx.Done():
		case response <- id:
		}
	}

	res := make(chan int)
	var wg sync.WaitGroup

	// #4
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go dial(ctx, listener.Addr().String(), res, i+1, &wg)
	}

	// #5
	response := <-res
	cancel()
	wg.Wait()
	close(res)

	// #6
	if ctx.Err() != context.Canceled {
		t.Errorf("expected canceled context; actual: %s", ctx.Err())
	}

	t.Logf("dialer %d retrieved the resource", response)

}

/*
			다중 다이얼러 취소


하나의 콘텍스트를 여러 개의 DialContext 함수 호출에 넘겨서 해당 콘텍스트의 cancel 함수를 호출함으로써 동시에 여러 개
다이얼(dial) 요청을 취소할 수도 있다.

여러 개의 서버에서 TCP를 통해 단 하나의 리소스만을 받아 올 필요가 있다고 가정하자. 각 서버에 비동기적으로 연결을 요청하고, 동일한
콘텍스트를 각 다이얼러(dialer)에 연결한다. 한 서버로부터 응답이 왔다면 다른 응답은 필요가 없으니 나머지 다이얼러들은 콘텍스트를
취소하여 연결 시도를 중단할 수 있다.




context.WithDeadline 함수(#1)를 사용하여 콘텍스트를 생성 시 콘텍스트의 Err 메서드에서는 잠재적으로 context.Canceled, context.DeadlineExceeded
또는 nil 이라는 총 세 개 중 하나의 값을 반환한다. 테스트상에서 cancel 함수로 다이얼러에 대한 연결을 중단하였으므로 Err 메서드는 context.Canceledconst을
반환할 것이다.

먼저 리스너를 생성한다. 리스너는 하나의 연결을 수락하고 성공적으로 핸드셰이크(#2)를 마치면 연결을 종료한다. 이후에 다이얼러를 생성한다. 여러 개의
다이얼러를 실행하기 때문에 다이얼링을 위한 코드를 추상화하여 별도의 함수로 분리한다(#3). 이 익명함후는 DialContext 함수를 사용하여 매개변수로 주어진
주소로 연결을 시도한다. 연결이 성공하면 아직 콘텍스트를 취소하지 않았다고 생각하고 다이얼러의 ID를 응답 채널에 전송한다. for 루프(#4)를 사용하고 별도의
고루틴을 호출하여 여러 개의 다이얼러를 생성한다.

다른 다이얼러가 먼저 연결되어 DialContext 함수의 다이얼링이 블로킹 된 경우 이를 해제하기 위해 cancel 함수를 호출하거나 데드라인을 통해 콘텍스트를 취소한다.
WaitGroup을 이용하여 콘텍스트를 취소하여 for(#4)에서 생성한 모든 다이얼 고루틴을 정상적으로 종료한다.


고루틴이 정상적으로 작동하면 한 연결 시도는 다른 연결 시도보다 먼저 성공적으로 리스너에게 연결될 수 있다. 연결이 성공한 다이얼러의 ID를 res 채널에서 받는다(#5).
이후 콘텍스트를 취소하여 다른 다이얼러들의 연결 시도를 중단한다. 이 지점에서 wg.Wait 함수는 다른 다이얼러들의 연결 시도 중단이 끝나고 고루틴이 종료될 때까지 블로킹
된다.

마지막으로, 발생한 콘텍스트 취소가 코드상의 취소였음을 확인한다(#6).
*/
