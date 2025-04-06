package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gad/gophercises-cyoa/story"
)



func main() {

	jsonFile := flag.String("file","gopher.json", "the json file containing the story")
	tmplFile := flag.String("tmpl", "templates/arc.tmpl", "the html template file")
	flag.Parse()

	chapters, err := story.StoryParsing(*jsonFile)
	if err != nil {
		log.Panic(err)
	}

	err = story.HtmlGenerate(chapters, *tmplFile, story.HtmlTempDir)
	if err != nil {
		log.Panic(err)
	}

	story.RegisterHandlers(chapters)

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
	if err = story.Cleaning(story.HtmlTempDir); err != nil {
		log.Fatalf("Error cleaning temporary html files : %v", err)
	} else {
		log.Println("Done")
	}
}

