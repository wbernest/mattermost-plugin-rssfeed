/*
Package atomparser
*/

package atomparser

import (
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"golang.org/x/tools/blog/atom"
	"io/ioutil"
	"net/http"
	"strings"
)

// ParseString will be used to parse strings and will return the Atom object
func ParseString(s string) (*atom.Feed, error) {
	feed := atom.Feed{}
	if len(s) == 0 {
		return &feed, nil
	}

	decoder := xml.NewDecoder(strings.NewReader(s))
	decoder.CharsetReader = charset.NewReaderLabel
	err := decoder.Decode(&feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}

// ParseURL will be used to parse a string returned from a url and will return the Rss object
func ParseURL(url string) (*atom.Feed, string, error) {
	byteValue, err := getContent(url)
	if err != nil {
		return nil, "", err
	}

	decoder := xml.NewDecoder(strings.NewReader(string(byteValue)))
	decoder.CharsetReader = charset.NewReaderLabel
	feed := atom.Feed{}
	err = decoder.Decode(&feed)
	if err != nil {
		return nil, "", err
	}

	return &feed, string(byteValue), nil
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

// CompareItems - This function will used to compare 2 atom feed xml item objects
// and will return a list of differing items
func CompareItems(feedOne *atom.Feed, feedTwo *atom.Feed) []*atom.Entry {
	biggerFeed := feedOne
	smallerFeed := feedTwo
	itemList := []*atom.Entry{}
	if len(feedTwo.Entry) > len(feedOne.Entry) {
		biggerFeed = feedTwo
		smallerFeed = feedOne
	} else if len(feedTwo.Entry) == len(feedOne.Entry) {
		return itemList
	}

	for _, item1 := range biggerFeed.Entry {
		exists := false
		for _, item2 := range smallerFeed.Entry {
			if item1.ID == item2.ID {
				exists = true
				break
			}
		}
		if !exists {
			itemList = append(itemList, item1)
		}
	}
	return itemList
}

// CompareItemsBetweenOldAndNew - This function will used to compare 2 atom xml event objects
// and will return a list of items that are specifically in the newer feed but not in
// the older feed
func CompareItemsBetweenOldAndNew(feedOld *atom.Feed, feedNew *atom.Feed) []*atom.Entry {
	itemList := []*atom.Entry{}

	for _, item1 := range feedNew.Entry {
		exists := false
		for _, item2 := range feedOld.Entry {
			if item1.ID == item2.ID {
				exists = true
				break
			}
		}
		if !exists {
			itemList = append(itemList, item1)
		}
	}
	return itemList
}

// IsValidFeed checks feed to see if it is an Atom feed
func IsValidFeed(url string) bool {
	_, _, err := ParseURL(url)
	if err == nil {
		return true
	}

	return false
}
