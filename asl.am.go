package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	dir    = "/home/asim/cache"
	domain = "asl.am"
	words  []string
)

func init() {
	rand.Seed(time.Now().UnixNano())

	b, err := ioutil.ReadFile("/usr/share/dict/words")
	if err != nil {
		panic("cant get words")
	}

	for _, word := range bytes.Split(b, []byte("\n")) {
		words = append(words, strings.ToLower(string(word)))
	}
}

func encode(u string) string {
	return base64.StdEncoding.EncodeToString([]byte(u))
}

func decode(u string) string {
	d, _ := base64.StdEncoding.DecodeString(u)
	return string(d)
}

func random() string {
	var word string
	i := len(words) - 1

	for len(word) < 512 {
		word += (words[rand.Int()%i] + "-")
	}

	return word[:len(word)-1]
}

func main() {
	http.HandleFunc("/u/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 {
			http.Redirect(w, r, "/", 302)
			return
		}

		var u string
		if len(parts) >= 3 {
			u = decode(strings.Join(parts[2:len(parts)-1], "/"))
		} else {
			u = decode(parts[2])
		}

		pu, err := url.Parse(u)
		if err != nil || !pu.IsAbs() {
			http.Redirect(w, r, "/", 302)
			return
		}

		qu, _ := url.QueryUnescape(u)
		http.Redirect(w, r, qu, 301)
	})

	http.HandleFunc("/lengthen", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		u := r.Form.Get("url")
		pu, err := url.Parse(u)
		if err != nil || !pu.IsAbs() {
			http.Error(w, "Invalid url", 502)
			return
		}

		ul := fmt.Sprintf("http://%s/u/%s/%s", domain, encode(u), random())
		fmt.Fprintf(w, ul)

	})

	log.Fatal(http.ListenAndServe(":9999", nil))
}
