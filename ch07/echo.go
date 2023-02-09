package ch07

import (
	"context"
	"net"
	"os"
)

// 참고: window와 WSL(Windows Subsystem for Linux)은 unixgram 도메인 소켓을 지원하지 않는다.


func streamingEchoServer(ctx context.Context, network string, addr string) (net.Addr, error) {

	s, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}

	go func() {
		go func() {
			<-ctx.Done() //#1
			_ = s.Close()
		}()
		for {
			conn, err := s.Accept() //#2
			if err != nil {
				return
			}

			go func() {
				defer func() {
					_ = conn.Close()
				}()
				for {

					buf := make([]byte, 1024)
					n, err := conn.Read(buf) //#3
					if err != nil {
						return
					}

					_, err = conn.Write(buf[:n]) //#4
					if err != nil {
						return
					}

				}
			}()
		}
	}()

	return s.Addr(), nil
}

func datagramEchoServer(ctx context.Context, network string, addr string) (net.Addr, error) {

	s, err := net.ListenPacket(network, addr)//#1
	if err != nil {
		return nil, err
	}

	go func() {
		go func() {

			<-ctx.Done()
			_ = s.Close()
		}()
		if network == "unixgram" {
			_ = os.Remove(addr)//#2
		}

		buf := make([]byte, 1024)
		for {
			n, clientAddr, err := s.ReadFrom(buf)
			if err != nil {
				return
			}

			_, err = s.WriteTo(buf[:n], clientAddr)
			if err != nil {
				return
			}

		}
	}()
	return s.LocalAddr(), nil
}
