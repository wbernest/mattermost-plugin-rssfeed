module github.com/wbernest/mattermost-plugin-rssfeed

go 1.12

require (
	github.com/lunny/html2md v0.0.0-20181018071239-7d234de44546
	github.com/mattermost/mattermost-server/v5 v5.35.1
	github.com/pkg/errors v0.9.1
	github.com/wbernest/atom-parser v0.0.0-20190507183633-f862cce5996a
	github.com/wbernest/rss-v2-parser v0.0.0-20190507183749-19659d6a25f2
	golang.org/x/tools v0.1.0
)

replace willnorris.com/go/imageproxy => willnorris.com/go/imageproxy v0.8.1-0.20190422234945-d4246a08fdec
