package ch03

import (
	"net"
	"syscall"
	"testing"
	"time"
)

// #1
func DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {

	d := net.Dialer{
		Control: func(_, addr string, c syscall.RawConn) error {
			// #1
			return &net.DNSError{
				Err:         "connection timed out",
				Name:        addr,
				Server:      "127.0.0.1",
				IsTimeout:   true,
				IsTemporary: true,
			}
		},
		Timeout: timeout, // time.Duration
	}
	return d.Dial(network, address)
}

func TestDialTimeout(t *testing.T) {

	c, err := DialTimeout("tcp", "10.0.0.1:http", 3*time.Second)
	// #3
	if err == nil {
		c.Close()
		t.Fatal("Connection did not time out")
	}

	nErr, ok := err.(net.Error) // #4

	if !ok {
		t.Fatal(err)
	}

	if !nErr.Timeout() { // #5
		t.Fatal("error is not a timeout")
	}
}
/*
	DialTimeout 함수를 이용한 연결 시도에 대한 타임 아웃


Dial 함수를 이용하는데 잠재적인 문제가 있는데, 각 연결 시도의 타임아웃을 운영체제의 타임아웃 시간에
의존해야 한다는 것이다. 

네트워크 연결을 위해 Dial 함수를 호출했더니 연결이 되지 않아서 타임아웃이 되어야 한다고 가정하면, 운영체제가
타임아웃을 두 시간 뒤에 시킨다면 프로그램은 더 이상 실시간 반응형이 아니게 되며, 사용자는 아주 오랜 기간 기다리게
된다.


프로그램을 예측 가능하도록 유지하고 사용자들을 행복하게 유지하려면 코드상에서 타임아웃을 제어해야 한다. 
서비스를 사용할 수 없다면 빠르게 타임아웃시키고 다음 서비스로 넘어 가는 게 좋다.


한 가지 방법으로는 명시적으로 연결마다 타임아웃 기간을 정하고 DialTimout 함수를 이용하는 것이다.




함수 net.DialTimeout(#1)가 net.Dialer 인터페이스에 대한 제어권을 제공하지 않기 때문에 테스트 코드에서
다이얼러 동작을 흉내낼 수 없다. 따라서 동일한 인터페이스를 갖는 별도의 구현체를 사용한다. 우리의 구현체인
DialTimeout 함수는 에러를 반환하기 위한 net.Dialer 인터페이스의 Control 함수(#2)를 오버라이딩한다.
DNS의 타임아웃 에러도 흉내낸다.


net.Dial 함수와는 달리 DialTimeout 함수는 추가적으로 타임아웃 기간(#3)에 대한 매개변수를 받는다. 이 경우
타임아웃 기간이 5초이기 때문에 연결이 5초 안에 되지 못할 경우 연결 시도는 아웃된다. 이 테스트에서, 라우팅을
할 수 없는 IP 주소인 10.0.0.0로 다이얼하면 연결 시도가 확실히 타임아웃될 것이다. 테스트가 성공하기 위해서는
Timeout 메서드(#5)에서 확인하기 이전에 에러를 net.Error(#4)로 타입 어설션해야 한다.


여러 IP 주소로 해석되는 호스트에 다이얼 시 Go에서는 각 IP주소 중 먼저 연결되는 주소를 기본 IP주소로 연결을
시도한다. 첫번째로 성공한 연결 시도에 대한 연결만 유지되고 그 외에 모든 연결 시도는 취소된다. 모든 연결 시도가
실패하거나 타임아웃되면 net.DialTimeout 에러를 반환한다.



*/