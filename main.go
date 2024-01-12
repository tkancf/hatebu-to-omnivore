package main

import (
	"encoding/csv"
	"encoding/xml"
	"flag"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
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

type RelatedLink struct {
	Title   string
	URL     string
	State   string // ARCHIVED, SUCCEEDED
	Tags    []string
	SavedAt int64
}

var (
	inputFilePath = flag.String("i", "", "required: File path for the input Atom feed file")
	stateBool     = flag.Bool("a", false, "optional: Set the state to ARCHIVED")
	saveState     = "SUCCEEDED"
)

func main() {
	flag.Parse()

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

	if err := OutputCSV(relatedLinks); err != nil {
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

func OutputCSV(relatedLinks []RelatedLink) error {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()
	header := []string{"url", "state", "labels", "saved_at", "published_at"}
	if err := writer.Write(header); err != nil {
		return err
	}
	for _, link := range relatedLinks {
		record := []string{
			link.Title,
			link.URL,
			link.State,
			formatLabels(link.Tags),
			strconv.FormatInt(link.SavedAt, 10),
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	return nil
}

func formatLabels(labels []string) string {
	if len(labels) == 0 {
		return ""
	}
	quotedLabels := make([]string, len(labels))
	for i, label := range labels {
		quotedLabels[i] = `"` + strings.ReplaceAll(label, `"`, `""`) + `"`
	}
	return "[" + strings.Join(quotedLabels, ",") + "]"
}
