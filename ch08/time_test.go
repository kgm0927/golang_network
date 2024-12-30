package main

import (
	"net/http"
	"testing"
	"time"
)

func TestHeadTime(t *testing.T) {
	resp, err := /*#1*/ http.Get("https://www.time.gov/")
	if err != nil {
		t.Fatal(err)
	}
	_ = /*#2*/ resp.Body.Close() // 예외 상황 처리 없이 항상 보디를 닫습니다.

	now := time.Now().Round(time.Second)
	date := /*#3*/ resp.Header.Get("Date")
	if date == "" {
		t.Fatal("no Date header received from time.gov")
	}

	dt, err := time.Parse(time.RFC1123, date)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("time.gov %s: (skew %s)", dt, now.Sub(dt))
}
