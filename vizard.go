package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func static(name, contentType string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		asset, err := Asset(name)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", contentType)
		w.Write(asset)
	}
}

func data(content []byte, contentType string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.Write(content)
	}
}

func main() {
	port := flag.Int("port", 8000, "HTTP service port")
	contentType := flag.String("content-type", "application/json", "Content-Type")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("Script name required as first argument")
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal("Stdin:", err)
	}

	var content []byte
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		content, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal("Input:", err)
		}
	} else {
		content = []byte{}
	}

	js, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		log.Fatal("Script:")
	}

	http.HandleFunc("/data", data(content, *contentType))
	http.HandleFunc("/bundle.js", static("bundle.js", "application/javascript"))
	http.HandleFunc("/styles.css", static("styles.css", "text/css"))
	http.HandleFunc("/script.js", data(js, "application/javascript"))
	http.HandleFunc("/", static("index.html", "text/html"))

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), nil))
}
