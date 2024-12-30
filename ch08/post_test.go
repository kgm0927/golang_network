package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type User struct {
	First string
	Last  string
}

func handlePostUser(t *testing.T) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = /*#2*/ io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		if r.Method != /*#3*/ http.MethodPost {
			/*#4*/ http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		var u User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			t.Error(err)
			http.Error(w, "Decode Failed", http.StatusBadRequest)
			return
		}
		/*#5*/ w.WriteHeader(http.StatusAccepted)
	}

}

func TestPostUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlePostUser(t)))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	if /*#1*/ resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d; actual status %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	u := User{First: "Adam", Last: "Woodbeck"}
	/*#2*/ err = json.NewEncoder(buf).Encode(&u)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = /*#3*/ http.Post(ts.URL, "application/json", buf)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != /*#4*/ http.StatusAccepted {
		t.Fatalf("expected status %d; actual status %d", http.StatusAccepted, resp.StatusCode)
	}
	_ = resp.Body.Close()
}

func TestMultipartPost(t *testing.T) {
	reqBody := /*#1*/ new(bytes.Buffer)
	w := /*#2*/ multipart.NewWriter(reqBody)

	for k, v := range map[string]string{
		"date":        time.Now().Format(time.RFC3339),
		"description": "Form values with attached files",
	} {
		err := /*#3*/ w.WriteField(k, v)
		if err != nil {
			t.Fatal(err)
		}
	}

	for i, file := range []string{
		"./files/hello.txt",
		"./files/goodbye.txt",
	} {
		filepart, err := /*#1*/ w.CreateFormFile(fmt.Sprintf("file%d", i+1), filepath.Base(file))
		if err != nil {
			t.Fatal(err)
		}

		f, err := os.Open(file)
		if err != nil {
			t.Fatal(err)
		}

		_, err = /*#2*/ io.Copy(filepart, f)
		_ = f.Close()
		if err != nil {
			t.Fatal(err)
		}
	}

	err := /*#3*/ w.Close()
	if err != nil {
		t.Fatal(err)
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost /*#1*/, "https://httpbin.org/post" /*#2*/, reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type" /*#3*/, w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d; actual status %d", http.StatusOK, resp.StatusCode)
	}
	t.Logf("\n%s", b)
}
