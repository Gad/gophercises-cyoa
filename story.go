package story

import (
	"encoding/json"
	"os"
)

type Arc string

type ArcContent struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  Arc `json:"Arc"`
	} `json:"options"`
}

// parse json file into a map of ArcContent
func StoryParsing(fileName string) (chapters map[Arc]ArcContent, err error) {

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
