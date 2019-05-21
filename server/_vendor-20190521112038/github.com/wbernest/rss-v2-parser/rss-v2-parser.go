/*Package rssv2parser -

rss-v2-parser.go

structures have been defined using the standards laid out
in https://validator.w3.org/feed/docs/rss2.html

Documentation for each structure has been pulled from the above link.*/
package rssv2parser

import (
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"net/http"
	"strings"
)

/*RSSV2 - What is RSS?
RSS is a Web content syndication format.

Its name is an acronym for Really Simple Syndication.

RSS is dialect of XML. All RSS files must conform to the XML 1.0 specification, as published on the World Wide Web Consortium (W3C) website.

At the top level, a RSS document is a <rss> element, with a mandatory attribute called version, that specifies the version of RSS that the document conforms to. If it conforms to this specification, the version attribute must be 2.0.

Subordinate to the <rss> element is a single <channel> element, which contains information about the channel (metadata) and its contents.

Sample files
Here are sample files for: RSS 0.91, 0.92 and 2.0.

Note that the sample files may point to documents and services that no longer exist. The 0.91 sample was created when the 0.91 docs were written. Maintaining a trail of samples seems like a good idea.

About this document
This document represents the status of RSS as of the Fall of 2002, version 2.0.1.

It incorporates all changes and additions, starting with the basic spec for RSS 0.91 (June 2000) and includes new features introduced in RSS 0.92 (December 2000) and RSS 0.94 (August 2002).

Change notes are here.

First we document the required and optional sub-elements of <channel>; and then document the sub-elements of <item>. The final sections answer frequently asked questions, and provide a roadmap for future evolution, and guidelines for extending RSS.*/
type RSSV2 struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

/*Channel - Required channel elements
Here's a list of the required channel elements, each with a brief description, an example, and where available, a pointer to a more complete description.

Element	Description	Example
title	The name of the channel. It's how people refer to your service. If you have an HTML website that contains the same information as your RSS file, the title of your channel should be the same as the title of your website.	GoUpstate.com News Headlines
link	The URL to the HTML website corresponding to the channel.	http://www.goupstate.com/
description	Phrase or sentence describing the channel.	The latest news from GoUpstate.com, a Spartanburg Herald-Journal Web site.


Optional channel elements
Here's a list of optional channel elements.

Element	Description	Example
language	The language the channel is written in. This allows aggregators to group all Italian language sites, for example, on a single page. A list of allowable values for this element, as provided by Netscape, is here. You may also use values defined by the W3C.	en-us
copyright	Copyright notice for content in the channel.	Copyright 2002, Spartanburg Herald-Journal
managingEditor	Email address for person responsible for editorial content.	geo@herald.com (George Matesky)
webMaster	Email address for person responsible for technical issues relating to channel.	betty@herald.com (Betty Guernsey)
pubDate	The publication date for the content in the channel. For example, the New York Times publishes on a daily basis, the publication date flips once every 24 hours. That's when the pubDate of the channel changes. All date-times in RSS conform to the Date and Time Specification of RFC 822, with the exception that the year may be expressed with two characters or four characters (four preferred).	Sat, 07 Sep 2002 0:00:01 GMT
lastBuildDate	The last time the content of the channel changed.	Sat, 07 Sep 2002 9:42:31 GMT
category	Specify one or more categories that the channel belongs to. Follows the same rules as the <item>-level category element. More info.	<category>Newspapers</category>
generator	A string indicating the program used to generate the channel.	MightyInHouse Content System v2.3
docs	A URL that points to the documentation for the format used in the RSS file. It's probably a pointer to this page. It's for people who might stumble across an RSS file on a Web server 25 years from now and wonder what it is.	http://backend.userland.com/rss
cloud	Allows processes to register with a cloud to be notified of updates to the channel, implementing a lightweight publish-subscribe protocol for RSS feeds. More info here.	<cloud domain="rpc.sys.com" port="80" path="/RPC2" registerProcedure="pingMe" protocol="soap"/>
ttl	ttl stands for time to live. It's a number of minutes that indicates how long a channel can be cached before refreshing from the source. More info here.	<ttl>60</ttl>
image	Specifies a GIF, JPEG or PNG image that can be displayed with the channel. More info here.
textInput	Specifies a text input box that can be displayed with the channel. More info here.
skipHours	A hint for aggregators telling them which hours they can skip. More info here.
skipDays	A hint for aggregators telling them which days they can skip. More info here.*/
type Channel struct {
	Title          string    `xml:"title"`
	Link           string    `xml:"link"`
	Description    string    `xml:"description"`
	Language       string    `xml:"language"`
	Copyright      string    `xml:"copyright"`
	ManagingEditor string    `xml:"managingEditor"`
	WebMaster      string    `xml:"webMaster"`
	PubDate        string    `xml:"pubDate"`
	LastBuildDate  string    `xml:"lastBuildDate"`
	Category       string    `xml:"category"`
	Generator      string    `xml:"generator"`
	Docs           string    `xml:"docs"`
	Image          Image     `xml:"image"`
	Cloud          Cloud     `xml:"cloud"`
	TTL            string    `xml:"ttl"`
	ItemList       []Item    `xml:"item"`
	TextInput      TextInput `xml:"textInput"`
	SkipHours      []Hour    `xml:"skipHours"`
	SkipDays       []Day     `xml:"skipDays"`
}

/*Image - <image> sub-element of <channel>
<image> is an optional sub-element of <channel>, which contains three required and three optional sub-elements.

<url> is the URL of a GIF, JPEG or PNG image that represents the channel.

<title> describes the image, it's used in the ALT attribute of the HTML <img> tag when the channel is rendered in HTML.

<link> is the URL of the site, when the channel is rendered, the image is a link to the site. (Note, in practice the image <title> and <link> should have the same value as the channel's <title> and <link>.

Optional elements include <width> and <height>, numbers, indicating the width and height of the image in pixels. <description> contains text that is included in the TITLE attribute of the link formed around the image in the HTML rendering.

Maximum value for width is 144, default value is 88.

Maximum value for height is 400, default value is 31.*/
type Image struct {
	URL         string `xml:"url"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Width       string `xml:"width"`
	Height      string `xml:"height"`
	Description string `xml:"description"`
}

/*Cloud - <cloud> sub-element of <channel>
<cloud> is an optional sub-element of <channel>.

It specifies a web service that supports the rssCloud interface which can be implemented in HTTP-POST, XML-RPC or SOAP 1.1.

Its purpose is to allow processes to register with a cloud to be notified of updates to the channel, implementing a lightweight publish-subscribe protocol for RSS feeds.

<cloud domain="radio.xmlstoragesystem.com" port="80" path="/RPC2" registerProcedure="xmlStorageSystem.rssPleaseNotify" protocol="xml-rpc" />

In this example, to request notification on the channel it appears in, you would send an XML-RPC message to radio.xmlstoragesystem.com on port 80, with a path of /RPC2. The procedure to call is xmlStorageSystem.rssPleaseNotify.

A full explanation of this element and the rssCloud interface is here.*/
type Cloud struct {
}

/*Item - Elements of <item>
A channel may contain any number of <item>s. An item may represent a "story" -- much like a story in a newspaper or magazine; if so its description is a synopsis of the story, and the link points to the full story. An item may also be complete in itself, if so, the description contains the text (entity-encoded HTML is allowed), and the link and title may be omitted. All elements of an item are optional, however at least one of title or description must be present.

Element	Description	Example
title	The title of the item.	Venice Film Festival Tries to Quit Sinking
link	The URL of the item.	http://www.nytimes.com/2002/09/07/movies/07FEST.html
description	The item synopsis.	Some of the most heated chatter at the Venice Film Festival this week was about the way that the arrival of the stars at the Palazzo del Cinema was being staged.
author	Email address of the author of the item. More.	oprah@oxygen.net
category	Includes the item in one or more categories. More.	Simpsons Characters
comments	URL of a page for comments relating to the item. More.	http://www.myblog.org/cgi-local/mt/mt-comments.cgi?entry_id=290
enclosure	Describes a media object that is attached to the item. More.	<enclosure url="http://live.curry.com/mp3/celebritySCms.mp3" length="1069871" type="audio/mpeg"/>
guid	A string that uniquely identifies the item. More.	<guid isPermaLink="true">http://inessential.com/2002/09/01.php#a2</guid>
pubDate	Indicates when the item was published. More.	Sun, 19 May 2002 15:21:36 GMT
source	The RSS channel that the item came from. More.	<source url="http://www.quotationspage.com/data/qotd.rss">Quotes of the Day</source>
*/
type Item struct {
	Title       string    `xml:"title"`
	Author      string    `xml:"author"`
	Description string    `xml:"description"`
	Link        string    `xml:"link"`
	PubDate     string    `xml:"pubDate"`
	GUID        string    `xml:"guid"`
	Category    string    `xml:"category"`
	Comments    string    `xml:"comments"`
	Enclosure   Enclosure `xml:"enclosure"`
	Source      string    `xml:"source"`
}

/*TextInput - <textInput> sub-element of <channel>
A channel may optionally contain a <textInput> sub-element, which contains four required sub-elements.

<title> -- The label of the Submit button in the text input area.

<description> -- Explains the text input area.

<name> -- The name of the text object in the text input area.

<link> -- The URL of the CGI script that processes text input requests.

The purpose of the <textInput> element is something of a mystery. You can use it to specify a search engine box. Or to allow a reader to provide feedback. Most aggregators ignore it.
*/
type TextInput struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Name        string `xml:"name"`
	Link        string `xml:"link"`
}

/*Hour -
An XML element that contains up to 24 <hour> sub-elements whose value is a number between 0 and 23,
representing a time in GMT, when aggregators, if they support the feature,
may not read the channel on hours listed in the skipHours element.
The hour beginning at midnight is hour zero.

http://backend.userland.com/skipHoursDays#skiphours
*/
type Hour struct {
	Hour string `xml:"hour"`
}

/*Day -
skipDays

An XML element that contains up to seven <day> sub-elements whose value is Monday,
Tuesday, Wednesday, Thursday, Friday, Saturday or Sunday. Aggregators may not read
the channel during days listed in the skipDays element.

http://backend.userland.com/skipHoursDays#skiphours
*/
type Day struct {
	Day string `xml:"day"`
}

/*Enclosure -
<enclosure> sub-element of <item>
<enclosure> is an optional sub-element of <item>.

It has three required attributes. url says where the enclosure is located, length says how big it is in bytes, and type says what its type is, a standard MIME type.

The url must be an http url.

<enclosure url="http://www.scripting.com/mp3s/weatherReportSuite.mp3" length="12216320" type="audio/mpeg" />

A use-case narrative for this element is here
*/
type Enclosure struct {
	URL    string `xml:"url,attr"`
	Length string `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

// ParseString will be used to parse strings and will return the Rss object
func ParseString(s string) (*RSSV2, error) {
	rss := RSSV2{}
	if len(s) == 0 {
		return &rss, nil
	}

	decoder := xml.NewDecoder(strings.NewReader(s))
	decoder.CharsetReader = charset.NewReaderLabel
	err := decoder.Decode(&rss)
	if err != nil {
		return nil, err
	}
	return &rss, nil
}

// ParseURL will be used to parse a string returned from a url and will return the Rss object
func ParseURL(url string) (*RSSV2, string, error) {
	byteValue, err := getContent(url)
	if err != nil {
		return nil, "", err
	}

	decoder := xml.NewDecoder(strings.NewReader(string(byteValue)))
	decoder.CharsetReader = charset.NewReaderLabel
	rss := RSSV2{}
	err = decoder.Decode(&rss)
	if err != nil {
		return nil, "", err
	}

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

// CompareItems - This function will used to compare 2 RSS xml item objects
// and will return a list of differing items
func CompareItems(firstRSS *RSSV2, secondRSS *RSSV2) []Item {
	biggerRSS := firstRSS
	smallerRSS := secondRSS
	itemList := []Item{}
	if len(secondRSS.Channel.ItemList) > len(firstRSS.Channel.ItemList) {
		biggerRSS = secondRSS
		smallerRSS = firstRSS
	} else if len(secondRSS.Channel.ItemList) == len(firstRSS.Channel.ItemList) {
		return itemList
	}

	for _, item1 := range biggerRSS.Channel.ItemList {
		exists := false
		for _, item2 := range smallerRSS.Channel.ItemList {
			if len(item1.GUID) > 0 && item1.GUID == item2.GUID {
				exists = true
				break
			} else if item1.PubDate == item2.PubDate && item1.Title == item2.Title {
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

// CompareItemsBetweenOldAndNew - This function will used to compare 2 RSS xml item objects
// and will return a list of items that are specifically in the newer feed but not in
// the older feed
func CompareItemsBetweenOldAndNew(oldRSS *RSSV2, newRSS *RSSV2) []Item {
	itemList := []Item{}

	for _, item1 := range newRSS.Channel.ItemList {
		exists := false
		for _, item2 := range oldRSS.Channel.ItemList {
			if len(item1.GUID) > 0 && item1.GUID == item2.GUID {
				exists = true
				break
			} else if item1.PubDate == item2.PubDate && item1.Title == item2.Title {
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

// IsValidFeed checks feed to see if it is RSS v2
func IsValidFeed(url string) bool {
	_, _, err := ParseURL(url)
	if err == nil {
		return true
	}

	return false
}
