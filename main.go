package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	host = "https://bbs.byr.cn"
)

var (
	myHost = []byte("https://byr.zeegg.com")
)

func replace(raw []byte) []byte {
	raw = bytes.Replace(raw, []byte(host), myHost, -1)
	return raw
}

func handler(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(r.Method, host+r.URL.String(), nil)
	if err != nil {
		log.Print(err)
		return
	}

	var client http.Client
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	header := w.Header()
	for k, v := range resp.Header {
		header[k] = v
	}

	w.WriteHeader(resp.StatusCode)

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return
	}

	raw = replace(raw)
	_, err = w.Write(raw)
	if err != nil {
		log.Print(err)
		return
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	port := "33037"
	http.HandleFunc("/", handler)

	log.Printf("ByrRss started on port %s...", port)
	http.ListenAndServe(":"+port, nil)
}
