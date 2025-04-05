package story

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

// generate html files for each chapter based on a template file
func HtmlGenerate(chapters FullStory, tmplFile string, htmlTempDir string) error {

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

// cleans the html temp dir
func Cleaning(htmlTempDir string) error {
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