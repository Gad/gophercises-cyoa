package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

type arc string
type ArcHandler struct {
	arc arc
}

func main() {

	f, err := os.Open("gopher.json")
	if err != nil {
		log.Panicln(err.Error())
	}

	st, err := io.ReadAll(f)
	if err != nil {
		log.Println(err.Error())
	}

	type arcContent struct {
		Title   string   `json:"title"`
		Story   []string `json:"story"`
		Options []struct {
			Text string `json:"text"`
			Arc  string `json:"arc"`
		} `json:"options"`
	}

	chapters := map[arc]arcContent{}

	err = json.Unmarshal(st, &chapters)
	if err != nil {
		log.Panicln(err.Error())
	}

	var tmplFile = "templates/arc.tmpl"
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}

	for arc := range chapters {

		htmlPath := fmt.Sprintf("html/arc_%s.html", arc)
		f, err = os.Create(htmlPath)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		err = tmpl.Execute(f, chapters[arc])
		if err != nil {
			panic(err)
		}
	}
	for arc := range chapters {
		archHandler := &ArcHandler{
			arc: arc,
		}
		path := fmt.Sprintf("/%s", arc)
		http.Handle(path, archHandler)
	}

	http.ListenAndServe(":8000", nil)
	// need to do the cleaning
}

func (s *ArcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
	htmlPath := fmt.Sprintf("html/arc_%s.html", s.arc)
	f, err := os.Open(htmlPath)
	if err != nil {
		log.Panicln(err.Error())
	}

	st, err := io.ReadAll(f)
	if err != nil {
		log.Println(err.Error())
	}
	w.Write([]byte(st))
}
