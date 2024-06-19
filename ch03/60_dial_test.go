package ch03

import (
	"io"
	"net"
	"testing"
)

func TestDial(t *testing.T) {

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})

	go func() { // #1
		defer func() { done <- struct{}{} }()

		for {
			conn, err := listener.Accept() // #2
			if err != nil {
				t.Log(err)
				return
			}

			go func(c net.Conn) { // #3

				defer func() {
					c.Close()
					done <- struct{}{}
				}()
				buf := make([]byte, 1024)

				for {
					n, err := c.Read(buf) // #4
					if err != nil {
						if err != io.EOF {
							t.Error(err)
						}
						return
					}
					t.Logf("received: %q", buf[:n])
				}

			}(conn)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	// #5				//#6		//#7

	if err != nil {
		t.Fatal(err)
	}

	conn.Close() // #8
	<-done
	listener.Close() // #9
	<-done
}

//		서버와 연결 수립

/*	클라이언트 측면에서 Go 표준 라이브러리의 net 패키지를 사용하면 간단하게 서버로 연결하고
연결을 수립할 수 있다. 목록 3-3은 랜덤 포트의 127.0.0.1에서 리스너를 바인딩하고 있는 서버와
TCP 연결을 수립하는 절차를 보여 주는 테스트이다.
*/

/*
	IP 주소 127.0.0.1 에 클라이언트가 접속할 수 있는 리스너를 생성한다. 현재 랜덤 포트를 할당하고 있으면 리스너를
	고루틴에서 시작해서(#1)이후의 테스트에서 클라이언트 측에서 연결할 수 있게 한다.

	리스너 고루틴에서는 TCP 수신 연결을 루프에서 받아들이고 각 연결 처리 로직을 담당하고 고루틴을 시작한다.
	이 고루틴을 '핸들러(handler)'라고 부른다.

	지금은 소켓으로부터 1024 바이트를 일거서 수신한 데이터를 로깅한다는 것만 알면 된다.

	표준 라이브러리의 net.Dial 함수는 tcp 같은 네트워크의 종류(#6)와 IP 주소, 포트의 조합(#7)을 매개변수로 받는다는
	 점에서 net.Listen 함수와 유사하다. Dial 함수에서 두 번째 매개변수로 받은 IP 주소, 포트를 이용하여 리스너로 연결을
	 시도한다.

	IP 주소 대시에 호스트명을 사용할 수도 있으며, 포트 번호 대신에 http와 같은 서비스명을 사용할 수도 있다. 호스트명이 하나
	이상의 IP 주소로 해석되는 경우 Go는 연결이 성공할 때까지의 IP 주소에 연결을 시도한다.

	Dial 함수는 연결 객체(#5)와 에러 인터페이스 값을 반환한다.const





	리스너로 연결을 성공적으로 수립한 후 클라이언트 측에서는 우아한 종료를 시작한다(#8). FIN 패킷을 받고 나면 Read 메서드는
	io.EOF 에러를 반환하는데(#4), 이는 리스너 측에서는 반대편 연결이 종료되었다는 의미이다. 커넥션 핸들러는 연결 객체의 Close
	메서드를 호출하며 종료된다(#3). Close 메서드는 FIN 패킷을 전송하여 TCP 세션의 우아한 종료를 마무리 한다.const


	마지막으로 리스너를 종료한다(#9). 리스너의 Accept 메서드(#2)는 즉시 블로킹이 해제되고 에러를 반환한다. 이 에러는 무언가 실패했다는
	의미가 아니고 그냥 로깅하고 가면 된다.

	리스너의 고루틴(#1)이 종료되며 테스트가 완료된다.

*/
