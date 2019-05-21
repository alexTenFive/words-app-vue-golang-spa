package main

import (
	"encoding/json"
	"fmt"
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
	http.HandleFunc("/api/send", getText)
	http.HandleFunc("/api/results", getResults)
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func getResults(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rw.Header().Set("Content-Type", "application/json")
	jsWords, err := json.Marshal(words)

	if err != nil {
		panic(err)
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(jsWords)
}

func getText(rw http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		http.Error(rw, "Please send a request body", 400)
		return
	}

	lt := LargeText{}

	err := json.NewDecoder(req.Body).Decode(&lt)
	if err != nil {
		panic(err)
	}

	words = WordsCount(lt.Text)

	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(http.StatusOK)
	fmt.Fprintln(rw, words)

}

func WordsCount(str string) map[string]int {
	return count_words(get_words_from(strings.ToLower(str)))
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
