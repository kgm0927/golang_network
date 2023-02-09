package tftp

import (
	"bytes"
	"errors"
	"log"
	"net"
	"time"
)

type Server struct {
	Payload []byte        // #1 // 모든 읽기 요청에 반환될 페이로드
	Retries uint8         // #2// 전송 실패 시 재시도 횟수
	Timeout time.Duration // #3 // 전송 승인을 기다릴 시간
}

func (s Server) ListenAndServe(addr string) error {

	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}

	defer func() {
		_ = conn.Close()
	}()

	log.Printf("Listening on %s ...\n", conn.LocalAddr())

	return s.Serve(conn)
}

func (s *Server) Serve(conn net.PacketConn) error {
	if conn == nil {
		return errors.New("nil connection")
	}
	if s.Payload == nil {
		return errors.New("payload is required")
	}

	if s.Retries == 0 {
		s.Retries = 10
	}

	if s.Timeout == 0 {
		s.Timeout = 6 * time.Second
	}

	var rrq ReadReq

	for {

		buf := make([]byte, DatagramSize)

		_, addr, err := conn.ReadFrom(buf)
		if err != nil {
			return err
		}

		err = rrq.UnmarshalBinary(buf)
		if err != nil {
			log.Printf("[%s] bad request: %v", addr, err)
			continue
		}
		go s.handle(addr.String(), rrq)
	}
}

func (s Server) handle(clientAddr string, rrq ReadReq) { //#1

	log.Printf("[%s] requested file: %s", clientAddr, rrq.Filename)

	conn, err := net.Dial("udp", clientAddr) //#2

	if err != nil {
		log.Printf("[%s] dial: %v", clientAddr, err)
		return
	}

	defer func() {
		_ = conn.Close()
	}()

	var (
		ackPkt  Ack
		errPkt  Err
		dataPkt = Data{PayLoad: bytes.NewReader(s.Payload)} //#3
		buf     = make([]byte, DatagramSize)
	)

NEXTPACKET:
	for n := DatagramSize; n == DatagramSize; {
		data, err := dataPkt.MarshalBinary()
		if err != nil {
			log.Printf("[%s] preparing data packet: %v", clientAddr, err)
			return
		}

	RETRY:
		for i := s.Retries; i > 0; i-- { // #5
			n, err = conn.Write(data) // 데이터 패킷 전송 #6
			if err != nil {
				log.Printf("[%s] write %v", clientAddr, err)
				return
			}
			// 클라이언트의 ACK 패킷 대기
			_ = conn.SetReadDeadline(time.Now().Add(s.Timeout))

			_, err = conn.Read(buf)
			if err != nil {
				if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
					continue RETRY
				}
				log.Printf("[%s] waiting for ACK: %v", clientAddr, err)
				return
			}

			switch {
			case ackPkt.UnmarshalBinary(buf) == nil:
				if uint16(ackPkt) == dataPkt.Block { //#7
					// 클라이언트의 ACK 패킷 대기
					continue NEXTPACKET
				}
			case errPkt.UnmarshalBinary(buf) == nil:
				log.Printf("[%s] received error: %v", clientAddr, errPkt.Message)
				return

			default:
				log.Printf("[%s] bad packet", clientAddr)
			}
		}
		log.Printf("[%s] exhausted retries", clientAddr)
		return
	}

	log.Printf("[%s] sent %d blocks", clientAddr, dataPkt.Block)

}
