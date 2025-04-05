package story

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Handler of each arc/chapter. implements ServeHTTP Handler interface
type ArcHandler struct {
	arc Arc
}



func NewHandler(arc Arc) http.Handler {
	return &ArcHandler{arc}
}

// create a new endpoint for each chapter and register the handler
func RegisterHandlers(chapters FullStory) {

	for arc := range chapters {
		handler := NewHandler(arc)
		path := fmt.Sprintf("/%s", arc)
		http.Handle(path, handler)
	}
	http.Handle("/{$}", NewHandler("intro"))
	//all other paths, will return an empty ArcHandler
	// but serveHTTP will still process incoming request 
	http.Handle("/", NewHandler(""))
}


func (s *ArcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	log.Println(path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	
	// pour the html file corresponding to each handler into the http paquet
	arc := path[1:]
	htmlPath := fmt.Sprintf("%s/arc_%s.html", HtmlTempDir, arc)
	f, err := os.Open(htmlPath)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "No story here", http.StatusNotFound)
	}

	st, err := io.ReadAll(f)
	if err != nil {
		log.Println(err.Error())
	}
	w.Write([]byte(st))
}
