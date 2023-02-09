package tftp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strings"
)

const (
	DatagramSize = 516              //#1
	BlockSize    = DatagramSize - 4 //#2
)

type OpCode uint16 //#3

const (
	OpRRQ OpCode = iota + 1
	_
	OpData
	OpAck
	OpErr
)

type ErrCode uint16 //#4

const (
	ErrUnknown ErrCode = iota
	ErrNotFound
	ErrAccessViolation
	ErrDiskFill
	ErrIllagelOp
	ErrUnknownID
	ErrFileExists
	ErrNoUser
)

type ReadReq struct { // #1
	Filename string
	Mode     string
}

func (q ReadReq) MarshalBinary() ([]byte, error) {
	mode := "octet"
	if q.Mode != "" {
		mode = q.Mode
	}

	// OP코드+파일명+0바이트+모드 정보+0바이트

	cap := 2 + 2 + len(q.Filename) + 1 + len(q.Mode) + 1

	b := new(bytes.Buffer)
	b.Grow(cap)

	err := binary.Write(b, binary.BigEndian, OpRRQ) //OP 코드 쓰기 #2
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(q.Filename) // 파일명 쓰기
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0) // 0바이트 쓰기 //#3
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(mode) // 모드 정보 쓰기
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0) // 0바이트 쓰기 //#3
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (q *ReadReq) UnmarshalBinary(p []byte) error { //#1
	r := bytes.NewBuffer(p)

	var code OpCode

	err := binary.Read(r, binary.BigEndian, &code) // #2 OP 코드 읽기
	if err != nil {
		return err
	}

	if code != OpRRQ {
		return errors.New("invalid RRQ")
	}

	q.Filename, err = r.ReadString(0) // 파일명 읽기//#3
	if err != nil {
		return errors.New("invalid RRQ")
	}

	q.Filename = strings.TrimRight(q.Filename, "\x00") // 0바이트 제거 #4
	if len(q.Filename) == 0 {
		return errors.New("invalid RRQ")
	}

	q.Mode, err = r.ReadString(0) // 모드 정보 읽기
	if err != nil {
		return errors.New("invalid RRQ")
	}

	q.Mode = strings.TrimRight(q.Mode, "\x00") // 0바이트 제거
	if len(q.Mode) == 0 {
		return errors.New("invalid RRQ")
	}

	actual := strings.ToLower(q.Mode) // 강제 octet 모드 설정

	if actual != "octet" {
		return errors.New("only binary transfer supported")
	}

	return nil
}

type Data struct { //#1
	Block   uint16
	PayLoad io.Reader
}

func (d *Data) MarshalBinary() ([]byte, error) { //#2
	b := new(bytes.Buffer)
	b.Grow(DatagramSize)

	d.Block++

	err := binary.Write(b, binary.BigEndian, OpData) // OP 코드 쓰기
	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, d.Block) // 블록 번호 쓰기
	if err != nil {
		return nil, err
	}

	//BlockSize 만큼 쓰기

	_, err = io.CopyN(b, d.PayLoad, BlockSize) //#3
	if err != nil && err != io.EOF {
		return nil, err
	}

	return b.Bytes(), nil
}

func (d *Data) UnmarshalBinary(p []byte) error {

	if l := len(p); l < 4 || l > DatagramSize { //#1
		return errors.New("invalid DATA")
	}

	var opcode OpCode

	err := binary.Read(bytes.NewReader(p[:2]), binary.BigEndian, &opcode) //#2

	if err != nil || opcode != OpData {
		return errors.New("invalid DATA")
	}

	err = binary.Read(bytes.NewReader(p[2:4]), binary.BigEndian, &d.Block) //#3
	if err != nil {
		return errors.New("invalid DATA")
	}

	d.PayLoad = bytes.NewBuffer(p[4:]) //#4

	return nil
}

type Ack uint16 //#1

func (a Ack) MarshalBinary() ([]byte, error) {
	cap := 2 + 2 // 코드 번호 + 블록 번호

	b := new(bytes.Buffer)
	b.Grow(cap)

	err := binary.Write(b, binary.BigEndian, OpAck) // OP 코드 쓰기
	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, a) // 블록 번호 쓰기
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (a *Ack) UnmarshalBinary(p []byte) error {
	var code OpCode

	r := bytes.NewReader(p)

	err := binary.Read(r, binary.BigEndian, &code) // OP 코드 읽기
	if err != nil {
		return err
	}

	if code != OpAck {
		return errors.New("invalid ACK") // 블록 번호 읽기
	}

	return binary.Read(r, binary.BigEndian, a) // 블록 번호 읽기
}

type Err struct {
	Error   ErrCode
	Message string
}

func (e Err) MarshalBinary() ([]byte, error) {

	//OP 코드+에러 코드+메시지+ 0바이트
	cap := 2 + 2 + len(e.Message) + 1

	b := new(bytes.Buffer)
	b.Grow(cap)

	err := binary.Write(b, binary.BigEndian, OpErr)
	if err != nil {
		return nil, err
	}

	err = binary.Write(b, binary.BigEndian, e.Error)
	if err != nil {
		return nil, err
	}

	_, err = b.WriteString(e.Message)
	if err != nil {
		return nil, err
	}

	err = b.WriteByte(0)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (e *Err) UnmarshalBinary(p []byte) error {
	r := bytes.NewBuffer(p)
	var code OpCode

	err := binary.Read(r, binary.BigEndian, &code) //#1 OP 코드 읽기
	if err != nil {
		return err
	}

	if code != OpErr {
		return errors.New("invalid ERROR")
	}

	err = binary.Read(r, binary.BigEndian, &e.Error) //  #2 에러 메시지 읽기
	if err != nil {
		return err
	}

	e.Message, err = r.ReadString(0)                 //#3
	e.Message = strings.TrimRight(e.Message, "\x00") //#4

	return err

}
