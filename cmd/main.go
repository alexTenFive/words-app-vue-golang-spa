package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var words map[string]int

// LargeText represent large string
type LargeText struct {
	Text   string `json:"text"`
	Number int    `json:"number"`
}

func main() {
	text := make(chan string)
	results := make(chan map[string]int)
	handleText := getText(text)
	handleResults := getResults(results)
	go WordsCount(text, results)

	http.HandleFunc("/api/send", handleText)
	http.HandleFunc("/api/results", handleResults)
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func getResults(results <-chan map[string]int) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		rw.Header().Set("Content-Type", "application/json")

		var jsWords []byte
		var err error

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

		rw.WriteHeader(http.StatusOK)
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

func WordsCount(text <-chan string, results chan<- map[string]int) {
	for {
		results <- count_words(get_words_from(strings.ToLower(<-text)))
	}
}

func get_words_from(text string) []string {
	words := regexp.MustCompile("\\w+")
	return words.FindAllString(text, -1)
}

func count_words(words []string) map[string]int {
	wordCounts := make(map[string]int)
	for _, word := range words {
		wordCounts[word]++
	}
	return wordCounts
}
