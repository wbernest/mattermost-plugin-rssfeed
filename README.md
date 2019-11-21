# RSSFeed Plugin [![CircleCI branch](https://img.shields.io/circleci/project/github/wbernest/mattermost-plugin-rssfeed/master.svg)](https://circleci.com/gh/wbernest/mattermost-plugin-rssfeed/tree/master)

This plugin allows a user to subscribe a channel to an RSS (Version 2 only) or an Atom Feed.

- Version 0.1.0+ requires Mattermost 5.10
- Version < 0.1.0 requires Mattermost 5.6

## Getting Started
Upload tar.gz to Mattermost using the plugin screen.
Assign a user under the settings for posting.

To allow for the plugin to override user user name and icon on posts you must enable the feature in your Mattermost system settings:

* <a class="reference external" href="https://docs.mattermost.com/administration/config-settings.html#enable-integrations-to-override-usernames">Enable integrations to override usernames</a> must be set to `true` in `config.json` to override usernames. Enable them from <strong>System Console &gt; Integrations &gt; Custom Integrations</strong> or ask your System Administrator to do so. If not enabled, the username is set to `webhook`.
* <a class="reference external" href="https://docs.mattermost.com/administration/config-settings.html#enable-integrations-to-override-profile-picture-icons">Enable integrations to override profile picture icons</a> must be set to `true` in `config.json` to override profile picture icons. Enable them from <strong>System Console &gt; Integrations &gt; Custom Integrations</strong> or ask your System Administrator to do so. If not enabled, the icon of the creator of the webhook URL is used to post messages.

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
