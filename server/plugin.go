package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/tools/blog/atom"
	"github.com/lunny/html2md"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	atomparser "github.com/wbernest/atom-parser"
	rssv2parser "github.com/wbernest/rss-v2-parser"
)

//const RSSFEED_ICON_URL = "./plugins/rssfeed/assets/rss.png"

// RSSFeedPlugin Object
type RSSFeedPlugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	botUserID            string
	processHeartBeatFlag bool
}

// ServeHTTP hook from mattermost plugin
func (p *RSSFeedPlugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	switch path := r.URL.Path; path {
	case "/images/rss.png":
		data, err := ioutil.ReadFile(string("plugins/rssfeed/assets/rss.png"))
		if err == nil {
			w.Header().Set("Content-Type", "image/png")
			w.Write(data)
		} else {
			w.WriteHeader(404)
			w.Write([]byte("404 Something went wrong - " + http.StatusText(404)))
			p.API.LogInfo("/images/rss.png err = ", err.Error())
		}
	default:
		w.Header().Set("Content-Type", "application/json")
		http.NotFound(w, r)
	}
}

func (p *RSSFeedPlugin) setupHeartBeat() {
	heartbeatTime, err := p.getHeartbeatTime()
	if err != nil {
		p.API.LogError(err.Error())
	}

	for p.processHeartBeatFlag {
		//p.API.LogDebug("Heartbeat")

		err := p.processHeartBeat()
		if err != nil {
			p.API.LogError(err.Error())

		}
		time.Sleep(time.Duration(heartbeatTime) * time.Minute)
	}
}

func (p *RSSFeedPlugin) processHeartBeat() error {
	dictionaryOfSubscriptions, err := p.getSubscriptions()
	if err != nil {
		return err
	}

	for _, value := range dictionaryOfSubscriptions.Subscriptions {
		err := p.processSubscription(value)
		if err != nil {
			p.API.LogError(err.Error())
		}
	}

	return nil
}

func (p *RSSFeedPlugin) getHeartbeatTime() (int, error) {
	config := p.getConfiguration()
	heartbeatTime := 15
	var err error
	if len(config.Heartbeat) > 0 {
		heartbeatTime, err = strconv.Atoi(config.Heartbeat)
		if err != nil {
			return 15, err
		}
	}

	return heartbeatTime, nil
}

func (p *RSSFeedPlugin) processSubscription(subscription *Subscription) error {

	if len(subscription.URL) == 0 {
		return errors.New("no url supplied")
	}

	if rssv2parser.IsValidFeed(subscription.URL) {
		err := p.processRSSV2Subscription(subscription)
		if err != nil {
			return errors.New("invalid RSS v2 feed format - " + err.Error())
		}

	} else if atomparser.IsValidFeed(subscription.URL) {
		err := p.processAtomSubscription(subscription)
		if err != nil {
			return errors.New("invalid atom feed format - " + err.Error())
		}
	} else {
		return fmt.Errorf("invalid feed format for subscription: %s", subscription.URL)
	}

	return nil
}

func (p *RSSFeedPlugin) processRSSV2Subscription(subscription *Subscription) error {
	config := p.getConfiguration()

	// get new rss feed string from url
	newRssFeed, newRssFeedString, err := rssv2parser.ParseURL(subscription.URL)
	if err != nil {
		return err
	}

	// retrieve old xml feed from database
	oldRssFeed, err := rssv2parser.ParseString(subscription.XML)
	if err != nil {
		return err
	}

	items := rssv2parser.CompareItemsBetweenOldAndNew(oldRssFeed, newRssFeed)

	// if this is a new subscription only post the latest
	// and not spam the channel
	if len(oldRssFeed.Channel.ItemList) == 0 {
		items = items[:1]
	}

	for _, item := range items {
		post := ""

		if config.FormatTitle {
			post = post + "##### "
		}
		post = post + newRssFeed.Channel.Title + "\n"

		if config.ShowRSSItemTitle {
			if config.FormatTitle {
				post = post + "###### "
			}
			post = post + item.Title + "\n"
		}

		if config.ShowRSSLink {
			post = post + strings.TrimSpace(item.Link) + "\n"
		}
		if config.ShowDescription {
			post = post + html2md.Convert(item.Description) + "\n"
		}

		p.createBotPost(subscription.ChannelID, post, "custom_git_pr")
	}

	if len(items) > 0 {
		subscription.XML = newRssFeedString
		p.updateSubscription(subscription)
	}

	return nil
}

func (p *RSSFeedPlugin) processAtomSubscription(subscription *Subscription) error {
	config := p.getConfiguration()

	// get new rss feed string from url
	newFeed, newFeedString, err := atomparser.ParseURL(subscription.URL)
	if err != nil {
		return err
	}

	// retrieve old xml feed from database
	oldFeed, err := atomparser.ParseString(subscription.XML)
	if err != nil {
		return err
	}

	items := atomparser.CompareItemsBetweenOldAndNew(oldFeed, newFeed)

	// if this is a new subscription only post the latest
	// and not spam the channel
	if len(oldFeed.Entry) == 0 {
		items = items[:1]
	}

	for _, item := range items {
		post := ""

		if config.FormatTitle {
			post = post + "##### "
		}
		post = post + newFeed.Title + "\n"

		if config.ShowAtomItemTitle {
			if config.FormatTitle {
				post = post + "###### "
			}
			post = post + item.Title + "\n"
		}

		if config.ShowAtomLink {
			for _, link := range item.Link {
				if link.Rel == "alternate" {
					post = post + strings.TrimSpace(link.Href) + "\n"
				}
			}
		}

		if config.ShowSummary {
			if !tryParseRichNode(item.Summary, &post) {
				p.API.LogInfo("Missing summary in atom feed item",
					"subscription_url", subscription.URL,
					"item_title", item.Title)
				post = post + "\n"
			}
		}

		if config.ShowContent {
			if !tryParseRichNode(item.Content, &post) {
				p.API.LogInfo("Missing content in atom feed item",
					"subscription_url", subscription.URL,
					"item_title", item.Title)
				post = post + "\n"
			}
		}

		p.createBotPost(subscription.ChannelID, post, "custom_git_pr")
	}

	if len(items) > 0 {
		subscription.XML = newFeedString
		p.updateSubscription(subscription)
	}

	return nil
}

func tryParseRichNode(node *atom.Text, post *string) bool {
	if node != nil {
		if node.Type != "text" {
			*post = *post + html2md.Convert(strings.TrimSpace(node.Body)) + "\n"
		} else {
			*post = *post + node.Body + "\n"
		}
		return true
	} else {
		return false
	}
}

func (p *RSSFeedPlugin) createBotPost(channelID string, message string, postType string) error {
	post := &model.Post{
		UserId:    p.botUserID,
		ChannelId: channelID,
		Message:   message,
		Type:      postType,
		/*Props: map[string]interface{}{
			"from_webhook":      "true",
			"override_username": botDisplayName,
		},*/
	}

	if _, err := p.API.CreatePost(post); err != nil {
		p.API.LogError(err.Error())
		return err
	}

	return nil
}
