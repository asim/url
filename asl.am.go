package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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
	latest chan string
	reqCh  chan chan []string
	words  []string
	facts  map[string][]string
)

func loadFacts() {
	facts = make(map[string][]string)

	files, err := ioutil.ReadDir("/home/asim/facts")
	if err != nil {
		panic("cant get facts")
	}

	for _, file := range files {
		log.Println("loading ", file.Name())
		b, err := ioutil.ReadFile("/home/asim/facts/" + file.Name())
		if err != nil {
			panic("cant get words")
		}

		var quotes []string

		for _, quote := range bytes.Split(b, []byte("\n")) {
			quotes = append(quotes, strings.ToLower(string(quote)))
		}

		facts[file.Name()] = quotes
	}
}

func loadWords() {
	b, err := ioutil.ReadFile("/usr/share/dict/words")
	if err != nil {
		panic("cant get words")
	}

	for _, word := range bytes.Split(b, []byte("\n")) {
		words = append(words, strings.ToLower(string(word)))
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())

	loadFacts()
	loadWords()

	latest = make(chan string, 100)
	reqCh = make(chan chan []string)

	go func() {
		var urls []string

		for {
			select {
			case u := <-latest:
				urls = append(urls, u)
				if len(urls) > 5 {
					urls = urls[1:]
				}
			case ch := <-reqCh:
				ch <- urls
			}
		}
	}()
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

func randomFact(category string) string {
	qs, ok := facts[category]
	if !ok {
		return random()
	}

	return qs[rand.Int()%(len(qs)-1)]
}

func main() {
	http.HandleFunc("/latest", func(w http.ResponseWriter, r *http.Request) {
		ch := make(chan []string)
		reqCh <- ch
		urls := <-ch
		b, _ := json.Marshal(urls)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(b))
	})

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

		var fact string
		c := r.Form.Get("category")
		switch c {
		case "0":
			fact = randomFact("cats")
		case "1":
			fact = randomFact("chucknorris")
		case "2":
			fact = random()
		case "3":
			fact = randomFact("numbers")
		case "4":
			fact = randomFact("archer")
		default:
			fact = random()
		}

		ul := fmt.Sprintf("http://%s/u/%s/%s", domain, encode(u), url.QueryEscape(fact))
		fmt.Fprint(w, ul)

		select {
		case latest <- ul:
		default:
		}
	})

	log.Fatal(http.ListenAndServe(":9999", nil))
}
