# RSSFeed Plugin ![CircleCI branch](https://img.shields.io/circleci/project/github/mattermost/mattermost-plugin-sample/master.svg)

This plugin allows a user to subscribe a channel to an RSS (Version 2 only) Feed.

## Getting Started
Upload tar.gz to Mattermost using the plugin screen.
Asign a user under the settings for posting.

To use the plugin, navigate to the channel you want subscribed and use the following commands:
```
/feed help                  // to see the help menu
/feed subscribe <url>       // to subscribe the channel to an rss feed
/feed unsubscribe <url>     // to unsubscribe the channel from an rss feed
/feed list                  // to list the feeds the channel is subscribed to
```

## Developers
Clone the repository:
```
git clone https://github.com/wbernest/mattermost-plugin-rssfeed
```

Build your plugin:
```
make dist
```

This will produce a single plugin file (with support for multiple architectures) for upload to your Mattermost server:

```
rssfeed.0.0.1.tar.gz
```
