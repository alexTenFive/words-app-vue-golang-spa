package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"
	"sync"
)

// SafeMap for storing words until new data not will be recieved
type SafeMap struct {
	val map[string]int
	m   sync.Mutex
}

// Get SafeMap
func (i *SafeMap) Get() map[string]int {
	i.m.Lock()
	defer i.m.Unlock()
	return i.val
}

// Set data to SafeMap
func (i *SafeMap) Set(val map[string]int) {
	i.m.Lock()
	defer i.m.Unlock()
	i.val = val
}

var words SafeMap

const staticDir = "pkg/http/web/app/dist/"

// LargeText represent large string
type LargeText struct {
	Text string `json:"text"`
}

func main() {
	// Serve front server
	handleVue := http.FileServer(Vue(staticDir))

	// Channel for recieve text
	text := make(chan string)
	// For send results
	results := make(chan map[string]int)

	handleText := getText(text)
	// handle text
	go WordsCount(text, results)
	// get results
	go Results(results)
	// Front
	http.Handle("/", handleVue)
	// API
	http.HandleFunc("/api/send", handleText)
	http.HandleFunc("/api/results", getResults)

	log.Fatal(http.ListenAndServe(":9090", nil))
}

func getResults(rw http.ResponseWriter, req *http.Request) {
	var jsWords []byte
	var err error

	/*
		If new data waiting save it like json and send to client
		else send data from store var words
	*/
	jsWords, err = json.Marshal(words.Get())

	if err != nil {
		http.Error(rw, "Error while encoding", 400)
		return
	}

	setHeaders(rw)
	// 200
	rw.WriteHeader(http.StatusOK)
	// JSON
	rw.Write(jsWords)
}

func getText(text chan<- string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		if req.Body == nil {
			http.Error(rw, "Please send a request body", 400)
			return
		}

		lt := LargeText{}

		err := json.NewDecoder(req.Body).Decode(&lt)
		if err != nil {
			http.Error(rw, "Error while decoding", 400)
			return
		}

		text <- lt.Text

		setHeaders(rw)
		rw.WriteHeader(http.StatusOK)
	}
}

func setHeaders(rw http.ResponseWriter) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "OPTIONS")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rw.Header().Set("Content-Type", "application/json")
}

/*
WordsCount
1. Receives large text from /api/send
2. Make all text to lower case
3. Split text to words by regexp
4. Count words quantity and format it like map[string]int
5. Pass data to results channel
*/
func WordsCount(text <-chan string, results chan<- map[string]int) {
	for {
		results <- countWords(getWordsFrom(strings.ToLower(<-text)))
	}
}

// Results save results
func Results(results <-chan map[string]int) {
	for {
		words.Set(<-results)
	}
}

////////////////

func getWordsFrom(text string) []string {
	words := regexp.MustCompile("\\w+")
	return words.FindAllString(text, -1)
}

func countWords(words []string) map[string]int {
	wordCounts := make(map[string]int)
	for _, word := range words {
		wordCounts[word]++
	}
	return wordCounts
}

// Vue custom string type
type Vue string

// Open opening index.html core dist file
func (v Vue) Open(name string) (http.File, error) {
	if ext := path.Ext(name); name != "/" && (ext == "" || ext == ".html") {
		name = "index.html"
	}
	return http.Dir(v).Open(name)
}
