package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

type arc string
type ArcHandler struct {
	arc arc
}
var tmplFile = "templates/arc.tmpl"
var htmlTempDir = "tmp/html"


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

	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}

	for arc := range chapters {

		fileName := fmt.Sprintf("%s/arc_%s.html", htmlTempDir, arc)

		f, err = os.Create(fileName)
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

	go func() {
		if err := http.ListenAndServe(":8000", nil); err != http.ErrServerClosed {
			log.Fatalf("Server error : %v", err)
		}

	}()
	// handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// cleaning temporary html files
	
	globPattern := fmt.Sprintf("%s/*.html",htmlTempDir)
	matches, err := filepath.Glob(globPattern)
	if err!=nil {
		log.Println(err)
	}
	for _,m := range matches {
		if err = os.Remove(m); err != nil {
			log.Fatalf("Error cleaning temporary html files : %v", err)
		}

	}

}

func (s *ArcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	htmlPath := fmt.Sprintf("%s/arc_%s.html", htmlTempDir, s.arc)
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
