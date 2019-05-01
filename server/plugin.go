package main

import (
	"errors"
	"github.com/lunny/html2md"
	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/wbernest/rss-v2-parser"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const RSSFEED_ICON_URL = "https://mattermost.gridprotectionalliance.org/plugins/rssfeed/images/rss.png"
const RSSFEED_USERNAME = "RSSFeed Plugin"

// RSSFeedPlugin Object
type RSSFeedPlugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

// ServeHTTP hook from mattermost plugin
func (p *RSSFeedPlugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	switch path := r.URL.Path; path {
	case "/images/rss.png":
		data, err := ioutil.ReadFile(string("plugins/rssfeed/server/dist/images/rss.png"))
		if err == nil {
			w.Header().Set("Content-Type", "image/png")
			w.Write(data)
		} else {
			w.WriteHeader(404)
			w.Write([]byte("404 Something went wrong - " + http.StatusText(404)))
			p.API.LogInfo("/imags/rss.png err = ", err.Error())
		}
	default:
		w.Header().Set("Content-Type", "application/json")
		http.NotFound(w, r)
	}
}

// OnActivate is a plugin hook from the Mattermost plugin API
func (p *RSSFeedPlugin) OnActivate() error {
	p.API.RegisterCommand(getCommand())
	go p.setupHeartBeat()
	return nil
}

func (p *RSSFeedPlugin) setupHeartBeat() {
	heartbeatTime, err := p.getHeartbeatTime()
	if err != nil {
		p.API.LogError(err.Error())
	}
	//p.API.LogDebug("Heartbeat time = " + string(heartbeatTime))
	for true {
		p.API.LogDebug("Heartbeat")

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
	config := p.getConfiguration()

	if len(subscription.URL) == 0 {
		return errors.New("no url supplied")
	}

	// get new rss feed string from url
	newRssFeed, newRssFeedString, err := rssv2parser.RssParseURL(subscription.URL)
	if err != nil {
		return err
	}

	// retrieve old xml feed from database
	oldRssFeed, err := rssv2parser.RssParseString(subscription.XML)
	if err != nil {
		return err
	}

	items := rssv2parser.CompareItemsBetweenOldAndNew(oldRssFeed, newRssFeed)

	for _, item := range items {
		post := item.Title + "\n" + item.Link + "\n"
		if config.ShowDescription {
			post = post + html2md.Convert(item.Description) + "\n"
		}
		p.createBotPost(subscription.ChannelID, post, model.POST_DEFAULT)
	}

	if len(items) > 0 {
		subscription.XML = newRssFeedString
		p.updateSubscription(subscription)
	}

	return nil
}

func (p *RSSFeedPlugin) createBotPost(channelID string, message string, postType string) error {
	config := p.getConfiguration()
	user, err := p.API.GetUserByUsername(config.Username)
	if err != nil {
		mlog.Error(err.Error())
		return err
	}
	post := &model.Post{
		UserId:    user.Id,
		ChannelId: channelID,
		Message:   message,
		Type:      postType,
		Props: map[string]interface{}{
			"override_username": RSSFEED_USERNAME,
			"override_icon_url": RSSFEED_ICON_URL,
			"from_webhook":      true,
		},
	}

	if _, err := p.API.CreatePost(post); err != nil {
		mlog.Error(err.Error())
		return err
	}

	return nil
}
