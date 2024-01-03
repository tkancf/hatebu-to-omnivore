package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type AtomFeed struct {
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Title string `xml:"title"`
	Link  []Link `xml:"link"`
}

type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

// relatedのリンクだけ残した構造体
type RelatedLink struct {
	Title string
	Link  string
}

func main() {
	f, err := os.Open("testdata/sample.atom")
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
		for _, link := range entry.Link {
			if link.Rel == "related" {
				relatedLinks = append(relatedLinks, RelatedLink{entry.Title, link.Href})
			}
		}
	}
	for _, link := range relatedLinks {
		fmt.Printf("%+v\n", link)
	}
}
