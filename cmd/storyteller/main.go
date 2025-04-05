package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	story "github.com/gad/gophercises-cyoa"
)

// Handler of each arc/chapter. implements ServeHTTP Handler interface
type ArcHandler struct {
	arc story.Arc
}

// pour the html file corresponding to each handler into the http paquet
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


const (
	tmplFile    = "templates/arc.tmpl"
	htmlTempDir = "tmp/html"
)




// generate html files for each chapter based on a template file
func htmlGenerate(chapters map[story.Arc]story.ArcContent, tmplFile string, htmlTempDir string) error {

	// 1- creates a new template from tmplFile
	// 2- creates as many html file as chapters in htmlTempDir
	// 3- applies the template to each chapter to create the html and save them into files

	for arc := range chapters {

		tmpl, err := template.ParseFiles(tmplFile)
		if err != nil {
			return err
		}

		if _, err := os.Stat(htmlTempDir); errors.Is(err, os.ErrNotExist){
			err = os.MkdirAll(htmlTempDir, os.ModePerm)
			if err != nil {
				return err
			}
		}

		fileName := fmt.Sprintf("%s/arc_%s.html", htmlTempDir, arc)

		f, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer f.Close()

		err = tmpl.Execute(f, chapters[arc])
		if err != nil {
			return err
		}
	}

	return nil
}

// create a new endpoint for each chapter and register the handler
func registerHandlers(chapters map[story.Arc]story.ArcContent) {

	
	for arc := range chapters {
		archHandler := &ArcHandler{
			arc: arc,
		}
		path := fmt.Sprintf("/%s", arc)
		http.Handle(path, archHandler)
	}
}

// cleans the html temp dir
func cleaning(htmlTempDir string) error {
	globPattern := fmt.Sprintf("%s/*.html", htmlTempDir)
	matches, err := filepath.Glob(globPattern)
	if err != nil {
		return err
	}
	for _, m := range matches {
		if err = os.Remove(m); err != nil {
			return err
		}

	}
	return nil
}



func main() {


	chapters, err := story.StoryParsing("gopher.json")
	if err != nil {
		log.Panic(err)
	}

	err = htmlGenerate(chapters, tmplFile, htmlTempDir)
	if err != nil {
		log.Panic(err)
	}

	registerHandlers(chapters)

	// run http.ListenAndServe in a go routine to manage graceful shutdown
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
	if err = cleaning(htmlTempDir); err != nil {
		log.Fatalf("Error cleaning temporary html files : %v", err)
	} else {
		log.Println("Done")
	}
}

