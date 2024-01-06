package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
)

type AtomFeed struct {
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Title string   `xml:"title"`
	Links []Link   `xml:"link"`
	Tags  []string `xml:"subject"`
}

type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

type RelatedLink struct {
	Title string
	Link  string
	Tags  []string
}

func main() {
	inputFilePath := flag.String("i", "", "required: File path for the input Atom feed file")
	flag.Parse()

	// Check if the input file path is provided
	if *inputFilePath == "" {
		flag.PrintDefaults()
		return
	}

	f, err := os.Open(*inputFilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	relatedLinks, err := ParseAtomFeed(f)
	if err != nil {
		panic(err)
	}

	for _, link := range relatedLinks {
		fmt.Printf("%+v\n", link)
	}
}

func ParseAtomFeed(r io.Reader) ([]RelatedLink, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var feed AtomFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, err
	}

	var relatedLinks []RelatedLink
	for _, entry := range feed.Entries {
		for _, link := range entry.Links {
			if link.Rel == "related" {
				relatedLinks = append(relatedLinks, RelatedLink{entry.Title, link.Href, entry.Tags})
			}
		}
	}
	return relatedLinks, nil
}
