package main

import (
	"io"
	"net"
)

func proxyConn(source, destination string) error {
	connSource, err := net.Dial("tcp", source) //#1

	if err != nil {
		return err
	}

	defer connSource.Close()

	connDestination, err := net.Dial("tcp", destination) //#2

	if err != nil {
		return err
	}

	defer connDestination.Close()

	go func() {

		//connSource에 대응하는 connDestination
		_, _ = io.Copy(connSource, connDestination)
	}()
	// connDestination 으로 메시지를 보내는 connSource
	_, err = io.Copy(connDestination, connSource)

	return err

}
