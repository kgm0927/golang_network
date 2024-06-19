package ch03

import (
	"context"
	"io"
	"time"
)

const defaultPingInterval = 30 * time.Second

func Pinger(ctx context.Context, w io.Writer, reset <-chan time.Duration) {
	var interval time.Duration

	select {
	case <-ctx.Done():
		return
	case interval = <-reset: // reset 채널에서 초기 간격을 받아옴. // #1
	default:

	}

	if interval <= 0 {
		interval = defaultPingInterval
	}

	timer := time.NewTimer(interval) // #2
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}()

	for {
		select {
		case <-ctx.Done(): // #3
			return
		case newInterval := <-reset: // #4
			if !timer.Stop() {
				<-timer.C
			}
			if newInterval > 0 {
				interval = newInterval
			}
		case <-timer.C: // #5
			if _, err := w.Write([]byte("ping")); err != nil {
				return
			}
		}
		_ = timer.Reset(interval) // #6

	}

}

/*

			하트비트 구현하기

네트워크 연결이 지속되어야 하기 때무에 애플리케이션 계층에서 긴 유휴 시간을 가져야만 하는 경우 데드라인을
지속해서 뒤로 설정하기 위해 노드 간에 하트비트(Heartbeat)를 구현해야 한다. 하트비트로 인해 연결상의 문제가
생긴 경우 데이터가 전송될 때까지 기다리는 것이 아니라 네트워크 상의 장애를 빠르게 파악하고 연결을 재시도할 수
있다.

그로써 애플리케이션상에서는 필요 시에 네트워크를 사용할 수 있는 상태임을 확신할 수 있다.



하트비트(heartbeat)는 네트워크 연결의 데드라인을 지속해서 뒤로 설정하기 위한 의도로 응답을 받기 위해 원격지로 보내는
메시지이다. 네트워크 노드에서는 심박(heartbeat)처럼 일정한 간격(interval)으로 이러한 메시지를 전송한다. 이 방법은 다양한
운영체제에서 적용될 뿐만 아니라 애플리케이션이 하트비트를 구현하고 있으니 애플리케이션이 사용하고 있는 네트워크 연결이 응답
가능하다는 것을 확신할 수 있다.

하트비트는 TCP keeplive를 차단할 수 있는 방화벽에 대해서도 잘 작동한다.


하트비트의 예에서 주고받는 메시지로 핑(ping) 메시지와 퐁(pong)메시지를 사용하였다.




Pinger 함수는 일정한 간격(interval)마다 핑 메시지를 전송한다. Pinger 함수는 고루틴에서 동작하도록 설계되었기 때문에 이후에
고루틴을 종료시키거나 메모리 누수를 방지하기 위해 첫 번째 매개변수로 콘텍스트를 받는다. 다른 두 개의 매개변수로 io.Writer
인터페이스와 타이머가 리셋될 경우 시그널을 보낼 수 있는 채널이 있다. 타이머의 초기 간격 설정을 위해 버퍼 채널을 생성하여
시간을 설정한다(#1).

타이머를 interval로 초기화하고(#2) 필요한 경우 defer을 사용하여 타이머의 채널의 값을 소비한다. 종료되지 않은 for 루프에는
select 구문이 있는데, 이 select 구문은 콘텍스트가 취소되거나, 타이머를 리셋하기 위해 시그널을 받았거나, 또는 타이머가
만료되었거나 세 가지 중 하나가 일어날 때 까지 블로킹한다.

콘텍스트가 취소된 경우(#3) 함수는 종료되어 더 이상의 핑은 전송되지 않는다. 리셋을 위한 시그널을 위한(#4) 타이머는 다음 select
구문 실행 전에 리셋된다(#6).


타이머가 만료되면(#5) 핑 메시지를 writer에 쓰고 다음 select 구문 실행 전에 타이머를 리셋한다. 핑 메시지를 writer에 쓴ㄴ 도중 연속적으로
발생하는 타임아웃들을 이 케이스문에 이용하여 추적할 수 있다. 이를 위해 콘텍스트의 cancel 함수를 전달하고, 연속적 타임아웃이 임계값을
넘게 되면 cancel 함수를 호출한다.




*/
