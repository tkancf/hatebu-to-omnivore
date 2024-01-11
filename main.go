package main

import (
	"encoding/xml"
	"flag"
	"io"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

type AtomFeed struct {
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Title  string    `xml:"title"`
	Links  []Link    `xml:"link"`
	Tags   []string  `xml:"subject"`
	Issued time.Time `xml:"issued"`
}

type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

// url,state,labels,saved_at,published_at
type RelatedLink struct {
	Title   string   `csv:"title"`
	URL     string   `csv:"url"`
	State   string   `csv:"state"` // ARCHIVED, SUCCEEDED
	Tags    []string `csv:"tags"`
	SavedAt int64    `csv:"saved_at"`
}

var (
	inputFilePath = flag.String("i", "", "required: File path for the input Atom feed file")
	stateBool     = flag.Bool("a", false, "optional: Set the state to ARCHIVED")
	saveState     = "SUCCEEDED"
)

func main() {
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
	// Marshal related links to CSV and output to standard output
	if err := gocsv.MarshalCSV(&relatedLinks, gocsv.DefaultCSVWriter(os.Stdout)); err != nil {
		panic(err)
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
	if *stateBool {
		saveState = "ARCHIVED"
	}
	for _, entry := range feed.Entries {
		for _, link := range entry.Links {
			if link.Rel == "related" {
				relatedLinks = append(relatedLinks,
					RelatedLink{
						Title:   entry.Title,
						URL:     link.Href,
						State:   saveState,
						Tags:    entry.Tags,
						SavedAt: entry.Issued.Unix(),
					})
			}
		}
	}
	return relatedLinks, nil
}
