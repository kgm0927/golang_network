package ch03

import (
	"context"
	"fmt"
	"io"
	"time"
)

func ExamplePinger() {
	ctx, cancel := context.WithCancel(context.Background())
	r, w := io.Pipe()
	done := make(chan struct{})
	resetimer := make(chan time.Duration, 1) // #1
	resetimer <- time.Second                 // 초기 핑 간격

	go func() {
		Pinger(ctx, w, resetimer)
		close(done)
	}()

	receivePing := func(d time.Duration, r io.Reader) {
		if d >= 0 {
			fmt.Printf("resetting timer (%s)\n", d)
			resetimer <- d
		}

		now := time.Now()
		buf := make([]byte, 1024)
		n, err := r.Read(buf)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("receive %q (%s) \n", buf[:n], time.Since(now).Round(100*time.Millisecond))

	}

	for i, v := range []int64{0, 200, 300, 0, -1, -1, -1} { //#2
		fmt.Printf("Run %d:\n", i+1)
		receivePing(time.Duration(v)*time.Millisecond, r)
	}

	cancel()
	<-done // 콘텍스트가 취소된 이후 pinger가 종료되었는지 확인
	// Output:
	// .

}

/*

실행 결과

Run 1:				// #3
resetting timer (0s)
receive "ping" (1s)
Run 2:				// #4
resetting timer (200ms)
receive "ping" (200ms)
Run 3:				// #5
resetting timer (300ms)
receive "ping" (300ms)
Run 4:				// #6
resetting timer (0s)
receive "ping" (300ms)
Run 5:				// #7
receive "ping" (300ms)
Run 6:
receive "ping" (300ms)
Run 7:
receive "ping" (300ms)
want:


*/

/*

이 예시에서 Pinger의 타이머를 리셋하기 위해 사용하는 시그널로 버퍼 채널(#1)을 생성한다. 채널을 Pinger 함수로 넘기기
전에 resetTimer 채널에서 초기 핑 타이머 간격으로 1초를 설정한다. Pinger의 타이머를 초기화하고 핑 메시지를 writer에
쓸 때 이 간격을 사용한다.

일련의 시간 간격을 밀리초 단위로 정의하여 만든 int64배열을 for 루프(#2)에서 순회하여 각각의 값을 receivePing 함수로 전달한다.
receivePing 함수는 핑 타이머를 주어진 값으로 초기화 하고 주어진 reader로부터 핑 메시지를 받을 때까지 대기한다. 마지막으로,
핑 메시지 수신에 걸린 시간을 표준 출력으로 출력한다. Go는 표준 출력의 결과가 예시의 실행 결과와 같은지 확인하며, 같으면 테스트가
성공으로 끝나게 된다.


for 루프의 첫 순회에서(#3) 시간 간격으로 0을 전달하였다. 0을 전달하면 Pinger는 이전에 사용한 시간 간격으로 타이머를 리셋한다(이 예시에서는
1초로 리셋된다). 예상대로 reader는 1초 후에 핑 메시지를 수신하게 된다. 두 번째 순회에서 핑 타이머를 200ms로 리셋한다(#4). 200ms가 지나면
reader는 핑 메시지를 수신한다. 세 번째 순회에서 핑 타이머를 300ms로 리셋하고(#5), 300ms 후에 핑 메시지를 수신한다.

네 번째 순회에서도 시간 간격으로 0을 전달하였는데(#6), 직전 값이 300ms로 핑 타이머가 리셋된다. 이렇게 시간 간격을 0으로 설정하여 직전 시간
간격을 활용하면 초기 핑 타이머 시간 간격을 별도로 관리할 필요가 없게 된다.

타이머의 시간 간격을 적당히 원하는 값으로 초기화시키고, 다음 핑 메시지를 미리 전송해야 할 때마다 타이머 시간 간격으로 0을 전달하여 리셋할 수
있다. 향후에는 핑 타이머 시간 간격을 변경하기 위해 resetTime 채널로 매번 보내지 않고 코드 한 줄만 변경하면 된다.


다섯 번째부터 일곱 번재 순회까지는(#7) 핑 타이머를 리셋하지 않고 그냥 핑 메시지를 수신한다. 예상대로 reader는 해당 순회에서 300ms 간격으로 핑
메시지를 수신한다.
*/
