package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"
)

// Variable for storing words until new data not wiil be recieved
var words map[string]int

const staticDir = "pkg/http/web/app/dist/"

// LargeText represent large string
type LargeText struct {
	Text   string `json:"text"`
	Number int    `json:"number"`
}

func main() {
	// Serve from server
	handleVue := http.FileServer(Vue(staticDir))

	// Channel for recieve text
	text := make(chan string)
	// For send results
	results := make(chan map[string]int)

	handleText := getText(text)
	handleResults := getResults(results)

	go WordsCount(text, results)

	// Front
	http.Handle("/", handleVue)
	// API
	http.HandleFunc("/api/send", handleText)
	http.HandleFunc("/api/results", handleResults)
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func indexHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}
	return http.HandlerFunc(fn)
}

func getResults(results <-chan map[string]int) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		rw.Header().Set("Content-Type", "application/json")

		var jsWords []byte
		var err error

		/*
			If new data waiting save it like json and send to client
			else send data from store var words
		*/
		select {
		case result := <-results:
			words = result
			jsWords, err = json.Marshal(result)
		default:
			jsWords, err = json.Marshal(words)
		}

		if err != nil {
			panic(err)
		}
		// 200
		rw.WriteHeader(http.StatusOK)
		// JSON
		rw.Write(jsWords)
	}
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
			panic(err)
		}

		text <- lt.Text

		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		rw.Header().Set("Content-Type", "application/json")

		rw.WriteHeader(http.StatusOK)
	}
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
