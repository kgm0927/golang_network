package ch04

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	BinaryType uint8 = iota + 1 //#1
	StringType                  //#2

	MaxPayloadSize uint32 = 10 << 20 //10 MB #3
)

var ErrMaxPayloadSize = errors.New("maximum payload size exceeded")

// binary 타입의 Payload 인터페이스 구현 마무리하기

// #4
type Payload interface {
	fmt.Stringer
	io.ReaderFrom
	io.WriterTo
	Bytes() []byte
}

type Binary []byte

//#1

func (m Binary) Bytes() []byte { // #2
	return m
}

func (m Binary) String() string { return string(m) } //#3

func (m Binary) WriteTo(w io.Writer) (int64, error) { //#4
	err := binary.Write(w, binary.BigEndian, BinaryType) // 1바이트 타입 #5
	if err != nil {
		return 0, err
	}
	var n int64 = 1

	err = binary.Write(w, binary.BigEndian, uint32(len(m))) // 4 바이트 타입 32bit이므로 //#6

	if err != nil {
		return n, err
	}
	n += 4
	o, err := w.Write(m) // 페이로드 #7

	return n + int64(o), err
}

func (m *Binary) ReadFrom(r io.Reader) (int64, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ) //#1
	if err != nil {
		return 0, err
	}

	var n int64 = 1
	if typ != BinaryType { //#2
		return n, errors.New("invaild Binary")

	}

	var size uint32

	err = binary.Read(r, binary.BigEndian, &size) //#3

	if err != nil {
		return n, err
	}
	n += 4
	if size > MaxPayloadSize { //4
		return n, ErrMaxPayloadSize
	}

	*m = make([]byte, size) //#5
	o, err := r.Read(*m)    //#6

	return n + int64(o), err

}

//String 타입 생성하기

type String string

func (m String) Bytes() []byte { // #1
	return []byte(m)
}
func (m String) String() string { // #2
	return string(m)
}

func (m String) WriteTo(w io.Writer) (int64, error) { //#3
	err := binary.Write(w, binary.BigEndian, StringType) //1 바이트 #4
	if err != nil {
		return 0, err
	}

	var n int64 = 1

	err = binary.Write(w, binary.BigEndian, uint32(len(m))) // 4바이트 크기
	if err != nil {
		return n, err
	}
	n += 4
	o, err := w.Write([]byte(m)) // #5 페이로드
	return n + int64(o), err
}

func (m *String) ReadFrom(r io.Reader) (int64, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)

	if err != nil {
		return 0, err
	}

	var n int64 = 1

	if typ != StringType { // #1
		return n, errors.New("invalid String")
	}

	var size uint32

	err = binary.Read(r, binary.BigEndian, &size)
	if err != nil {
		return n, err
	}

	n += 4

	buf := make([]byte, size)

	o, err := r.Read(buf)

	if err != nil {
		return n, err
	}
	*m = String(buf) // #2
	return n + int64(o), nil
}

// reader에서 바이트를 읽ㅇ서 Binary와 String타입으로 디코딩하기

// #1
func decode(r io.Reader) (Payload, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ) //#2

	if err != nil {
		return nil, err
	}

	var payload Payload //#3

	switch typ { //#4
	case BinaryType:
		payload = new(Binary)

	case StringType:
		payload = new(String)

	default:
		return nil, errors.New("unknown type")
	}

	_, err = payload.ReadFrom(io.MultiReader(bytes.NewReader([]byte{typ}), r)) //#5
	if err != nil {
		return nil, err
	}

	return payload, nil
}
