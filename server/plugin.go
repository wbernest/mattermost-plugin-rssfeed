package main

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"

	"github.com/mmcdole/gofeed"
)

const (
	API_ERROR_ID_NOT_CONNECTED = "not_connected"
	RSSFEED_ICON_URL           = "https://en.wikipedia.org/wiki/RSS#/media/File:Feed-icon.svg"
	RSSFEED_USERNAME           = "RSSFeed Plugin"
)

// RSSFeedPlugin Object
type RSSFeedPlugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

// OnActivate is a plugin hook from the Mattermost plugin API
func (p *RSSFeedPlugin) OnActivate() error {

	p.API.RegisterCommand(getCommand())
	p.setupHeartBeat()
	return nil
}

func (p *RSSFeedPlugin) setupHeartBeat() error {
	heartbeatTime := p.getHeartbeatTime()

	for true {
		time.Sleep(heartbeatTime * time.Minute)
		p.processHeartBeat()
	}
	return nil
}

func (p *RSSFeedPlugin) processHeartBeat() error {
	dictionaryOfSubscriptions, err := p.getSubscriptions()
	if err != nil {
		mlog.Error(err.Error())
		return nil
	}

	for _, value := range dictionaryOfSubscriptions.Subscriptions {
		p.processSubscription(value)
	}

	return nil
}

func (p *RSSFeedPlugin) getHeartbeatTime() time.Duration {
	config := p.getConfiguration()
	heartbeatTime := 15
	if len(config.Heartbeat) > 0 {
		time, err := strconv.Atoi(config.Heartbeat)
		if err == nil {
			heartbeatTime = time
		}
	}

	return time.Duration(heartbeatTime)
}

func (p *RSSFeedPlugin) processSubscription(subscription *Subscription) error {
	fp := gofeed.Parser{}

	// retrieve old xml feed from database
	oldRssFeed, err := fp.Parse(strings.NewReader(subscription.XML))
	if err != nil {
		return err
	}

	// pull null xml feed from url
	newRssFeed, err1 := fp.ParseURL(subscription.URL)
	if err1 != nil {
		return err1
	}

	// check each item in new feed to see if they exist in old feed
	// if they do not exist post the new item to the channel and update
	// xml in the subscription object
	postsMade := false
	for _, item := range newRssFeed.Items {
		exists := false
		for _, oldItem := range oldRssFeed.Items {
			if oldItem.GUID == item.GUID {
				exists = true
			}
		}

		// if the item does not exist post it to the correct channel
		if !exists {
			postsMade = true
		}
	}

	if postsMade {
		subscription.XML = newRssFeed.String()
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
			"from_webhook":      "true",
			"override_username": RSSFEED_USERNAME,
			"override_icon_url": RSSFEED_ICON_URL,
		},
	}

	if _, err := p.API.CreatePost(post); err != nil {
		mlog.Error(err.Error())
		return err
	}

	return nil
}
