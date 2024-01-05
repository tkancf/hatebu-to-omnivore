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

	// parse inputfile path
	f, err := os.Open(*inputFilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	var feed AtomFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		panic(err)
	}

	var relatedLinks []RelatedLink

	for _, entry := range feed.Entries {
		for _, link := range entry.Links {
			if link.Rel == "related" {
				relatedLinks = append(relatedLinks, RelatedLink{entry.Title, link.Href, entry.Tags})
			}
		}
	}
	for _, link := range relatedLinks {
		fmt.Printf("%+v\n", link)
	}
}
