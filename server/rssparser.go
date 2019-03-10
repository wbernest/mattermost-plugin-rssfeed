package main

import (
	"encoding/xml"
	//"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	//"github.com/mattermost/mattermost-server/mlog"
)

type Rss struct {
	XMLName     xml.Name `xml:"rss"`
	Version     string   `xml:"version,attr"`
	Title       string   `xml:"channel>title"`
	Link        string   `xml:"channel>link"`
	Description string   `xml:"channel>description"`
	ItemList    []Item   `xml:"channel>item"`
}

type Item struct {
	Title       string `xml:"title"`
	Creator     string `xml:"dc:creator"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Guid        string `xml:"guid"`
}

// RssParseString will be used to parse strings and will return the Rss object
func RssParseString(s string) (*Rss, error) {
	byteValue, err := ioutil.ReadAll(strings.NewReader(s))
	if err != nil {
		return nil, err
	}

	rss := Rss{}
	xml.Unmarshal(byteValue, &rss)
	return &rss, nil
}

// RssParseURL will be used to parse a string returned from a url and will return the Rss object
func RssParseURL(url string) (*Rss, string, error) {

	byteValue, err := getContent(url)
	if err != nil {
		return nil, "", err
	}

	rss := Rss{}
	//	mlog.Info(fmt.Sprintf("RSS URL Parser xml = %v", string(byteValue)))

	xml.Unmarshal(byteValue, &rss)
	//mlog.Info(fmt.Sprintf("RSS URL Parser rss = %v", rss.Title))

	return &rss, string(byteValue), nil
}

func getContent(url string) ([]byte, error) {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
