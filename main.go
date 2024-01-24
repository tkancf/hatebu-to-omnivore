package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
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
	apiToken      = flag.String("k", "", "required: Omnivore API token")
	apiUrl        = flag.String("u", "https://api-prod.omnivore.app/api/graphql", "optional: Omnivore API URL")
)

func main() {
	flag.Parse()

	if *inputFilePath == "" {
		fmt.Println("Input file path is required")
		flag.PrintDefaults()
		return
	}

	if *apiToken == "" {
		fmt.Println("API token is required")
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

	csv, err := OutputCSV(relatedLinks)
	if err != nil {
		panic(err)
	}

	url, err := GetSignedURL()
	if err != nil {
		panic(err)
	}

	if err := UploadToSignedUrl(url, csv); err != nil {
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

func OutputCSV(relatedLinks []RelatedLink) (io.Reader, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	header := []string{"url", "state", "labels", "saved_at", "published_at"}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	for _, link := range relatedLinks {
		record := []string{
			link.URL,
			link.State,
			formatLabels(link.Tags),
			strconv.FormatInt(link.SavedAt, 10),
			"", // set published_at to blank because it does not exist in the atom file of hatebu
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf, nil
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

func GetSignedURL() (string, error) {
	const mutation = `
		mutation UploadImportFile($type: UploadImportFileType!, $contentType: String!) {
			uploadImportFile(type: $type, contentType: $contentType) {
				... on UploadImportFileError {
					errorCodes
				}
				... on UploadImportFileSuccess {
					uploadSignedUrl
				}
			}
		}
	`
	variables := map[string]string{
		"type":        "URL_LIST",
		"contentType": "text/csv",
	}

	body := map[string]interface{}{
		"query":     mutation,
		"variables": variables,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON: %w", err)
	}

	req, err := http.NewRequest("POST", *apiUrl, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", *apiToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	uploadSignedUrl, ok := response["data"].(map[string]interface{})["uploadImportFile"].(map[string]interface{})["uploadSignedUrl"].(string)
	if !ok {
		return "", fmt.Errorf("error retrieving signed URL: %s", response["data"].(map[string]interface{})["uploadImportFile"].(map[string]interface{})["errorCodes"])
	}

	return uploadSignedUrl, nil
}

func UploadToSignedUrl(signedUrl string, data io.Reader) error {
	req, err := http.NewRequest(http.MethodPut, signedUrl, data)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "text/csv")
	req.Header.Set("type", "URL_LIST")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload: %s", resp.Status)
	}

	return nil
}
