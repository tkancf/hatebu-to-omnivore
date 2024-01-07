package main

import (
	"strings"
	"testing"
)

func TestParseAtomFeed(t *testing.T) {
	// Create a sample Atom feed XML
	xmlData := `
<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://purl.org/atom/ns#" xmlns:dc="http://purl.org/dc/elements/1.1/" xml:lang="ja">
  <title>tkancfのブックマーク</title>
  <link type="text/html" rel="alternate" href="https://b.hatena.ne.jp/tkancf/bookmark"/>
  <link type="application/x.atom+xml" rel="service.post" title="tkancfのブックマーク" href="https://b.hatena.ne.jp/atom/post"/>
  <entry>
    <id>tag:hatena.ne.jp,2005:bookmark-tkancf-4747196714236931727</id>
    <title>2023年の振り返り</title>
    <link type="text/html" rel="related" href="https://tkancf.com/blog/2023-summary/"/>
    <link type="text/html" rel="alternate" href="https://b.hatena.ne.jp/tkancf/20231231#bookmark-4747196714236931727"/>
    <link type="application/x.atom+xml" rel="service.edit" title="2023年の振り返り" href="https://b.hatena.ne.jp/atom/edit/4747196714236931727"/>
    <summary>年末なのでね、振り返りましたよ。</summary>
    <issued>2023-12-31T09:25:43Z</issued>
    <author>
      <name>tkancf</name>
    </author>
    <dc:subject>tkancf</dc:subject>
  </entry>
  <entry>
    <id>tag:hatena.ne.jp,2005:bookmark-tkancf-4746065390374578287</id>
    <title>Vimの設定整理した - 2020年版</title>
    <link type="text/html" rel="related" href="https://tkancf.com/blog/vim%E3%81%AE%E8%A8%AD%E5%AE%9A%E6%95%B4%E7%90%86%E3%81%97%E3%81%9F-2020%E5%B9%B4%E7%89%88/"/>
    <link type="text/html" rel="alternate" href="https://b.hatena.ne.jp/tkancf/20231207#bookmark-4746065390374578287"/>
    <link type="application/x.atom+xml" rel="service.edit" title="Vimの設定整理した - 2020年版" href="https://b.hatena.ne.jp/atom/edit/4746065390374578287"/>
    <summary></summary>
    <issued>2023-12-07T00:04:45Z</issued>
    <author>
      <name>tkancf</name>
    </author>
    <dc:subject>vim</dc:subject>
    <dc:subject>tkancf</dc:subject>
  </entry>
  <entry>
    <id>tag:hatena.ne.jp,2005:bookmark-tkancf-4743218953048896367</id>
    <title>GitHub Mobile + GitHub issueでメモが良い感じ</title>
    <link type="text/html" rel="related" href="https://tkancf.com/blog/2023-10-05/"/>
    <link type="text/html" rel="alternate" href="https://b.hatena.ne.jp/tkancf/20231007#bookmark-4743218953048896367"/>
    <link type="application/x.atom+xml" rel="service.edit" title="GitHub Mobile + GitHub issueでメモが良い感じ" href="https://b.hatena.ne.jp/atom/edit/4743218953048896367"/>
    <summary>書きました</summary>
    <issued>2023-10-06T15:19:42Z</issued>
    <author>
      <name>tkancf</name>
    </author>
  </entry>
</feed>
	`

	// Create a reader from the XML data
	reader := strings.NewReader(xmlData)

	// Call the ParseAtomFeed function
	relatedLinks, err := ParseAtomFeed(reader)
	if err != nil {
		t.Errorf("ParseAtomFeed returned an error: %v", err)
	}

	// Check the number of related links
	if len(relatedLinks) != 3 {
		t.Errorf("Expected 2 related links, got %d", len(relatedLinks))
	}

	// Check the values of the related links
	expectedLinks := []RelatedLink{
		{
			Title: "2023年の振り返り",
			Link:  "https://tkancf.com/blog/2023-summary/",
			Tags:  []string{"tkancf"},
		},
		{
			Title: "Vimの設定整理した - 2020年版",
			Link:  "https://tkancf.com/blog/vim%E3%81%AE%E8%A8%AD%E5%AE%9A%E6%95%B4%E7%90%86%E3%81%97%E3%81%9F-2020%E5%B9%B4%E7%89%88/",
			Tags:  []string{"vim", "tkancf"},
		},
		{
			Title: "GitHub Mobile + GitHub issueでメモが良い感じ",
			Link:  "https://tkancf.com/blog/2023-10-05/",
			Tags:  []string{},
		},
	}

	for i, expectedLink := range expectedLinks {
		if relatedLinks[i].Title != expectedLink.Title || relatedLinks[i].Link != expectedLink.Link || !equalSlices(relatedLinks[i].Tags, expectedLink.Tags) {
			t.Errorf("Expected related link %d to be %+v, got %+v", i, expectedLink, relatedLinks[i])
		}
	}
}

// Function to compare slices of strings
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
