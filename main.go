package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"os"
)

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

	type arc string

	stories := map[arc]arcContent{}

	err = json.Unmarshal(st, &stories)
	if err != nil {
		log.Panicln(err.Error())
	}

	var tmplFile = "templates/arc.tmpl"
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}

	f, err = os.Create("arc.html")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = tmpl.Execute(f, stories["intro"])
	if err != nil {
		panic(err)
	}

}
