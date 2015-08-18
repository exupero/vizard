package main

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/toqueteos/webbrowser"
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

func staticTemplate(name, contentType string, data interface{}) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		asset, err := Asset(name)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", contentType)

		t := template.Must(template.New(name).Parse(string(asset)))
		t.Execute(w, data)
	}
}

func data(content []byte, contentType string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.Write(content)
	}
}

func handleError(w http.ResponseWriter) {
	if err := recover(); err != nil {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusInternalServerError)
	}
}

func main() {
	port := flag.Int("port", 8000, "HTTP service port")
	svgFilename := flag.String("svg", "", "Filename to save SVG to")
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("Script name required as first argument")
	}
	scriptFilename := flag.Arg(0)

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

	http.HandleFunc("/svg", func(w http.ResponseWriter, req *http.Request) {
		defer handleError(w)

		if len(*svgFilename) == 0 {
			panic(errors.New("No filename given for saving an SVG"))
		}

		bin, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}

		contents := strings.TrimPrefix(string(bin), "data:image/svg+xml;base64,")
		svg, err := base64.StdEncoding.DecodeString(contents)
		if err != nil {
			panic(err)
		}

		f, err := os.Create(*svgFilename)
		if err != nil {
			panic(err)
		}

		_, err = f.Write(svg)
		if err != nil {
			panic(err)
		}
	})

	http.HandleFunc("/bundle.js", static("bundle.js", "application/javascript"))
	http.HandleFunc("/styles.css", static("styles.css", "text/css"))
	http.HandleFunc("/script.js", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, scriptFilename)
	})
	http.HandleFunc("/", staticTemplate("index.html", "text/html", template.JS(content)))

	go func(){
		webbrowser.Open(fmt.Sprintf("0.0.0.0:%d", *port))
	}()

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), nil))
}
