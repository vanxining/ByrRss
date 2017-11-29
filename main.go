package main

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	host = "https://bbs.byr.cn"
)

var (
	myHost = []byte("https://byr.zeegg.com")
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	search = [][]byte{
		[]byte(` src="/`),
		[]byte(` href="/att/`),
		[]byte(`byr.zeegg.com/#!`),
	}
	replace = [][]byte{
		[]byte(` src="` + host + "/"),
		[]byte(` href="` + host + "/att/"),
		[]byte(`byr.zeegg.com/`),
	}
)

func modifyPage(raw []byte) []byte {
	raw = bytes.Replace(raw, []byte(host), myHost, -1)

	for index := range search {
		raw = bytes.Replace(raw, search[index], replace[index], -1)
	}

	return raw
}

func handler(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(r.Method, host+r.URL.String(), nil)
	if err != nil {
		log.Print(err)
		return
	}

	client := http.Client{Transport: tr}
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

	raw = modifyPage(raw)
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
