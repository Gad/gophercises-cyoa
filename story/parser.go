package story

import (
	"encoding/json"
	"os"
)

const HtmlTempDir = "tmp/html"

type Arc string

type ArcContent struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  Arc    `json:"Arc"`
	} `json:"options"`
}

type FullStory map[Arc]ArcContent



// parse json file into a map of ArcContent
func StoryParsing(fileName string) (chapters FullStory, err error) {

	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	d := json.NewDecoder(f)

	err = d.Decode(&chapters)
	if err != nil {
		return nil, err
	}

	return chapters, nil
}
