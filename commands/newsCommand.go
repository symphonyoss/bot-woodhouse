package commands

import (
	"regexp"
	"github.com/SymphonyOSF/botexample/client"
	"fmt"
	"bytes"
	"github.com/SymphonyOSF/botexample/nytimes"
)

type NewsCommand struct {
	ApiKey string
}

var news nytimes.News

func (newsCommand NewsCommand) MatchRegex() (regex *regexp.Regexp) {
	regex = regexp.MustCompile(`(?i)/news`)
	return regex
}


func (newsCommand NewsCommand) Help() (help string) {
	return "<b>/news</b> - provides top 5 stories from the New York Times"
}

func (newsCommand NewsCommand) OnMessage(message client.V2Message, client client.BotClient, handlers []CommandHandler) {
	news = nytimes.News{APIKey:newsCommand.ApiKey}

	// first location after messageML
	fmt.Println("replying back with news in streamId=" + message.StreamId)
	var buffer bytes.Buffer
	buffer.WriteString("<messageML>")
	stories := news.TopStories()
	for i, story := range stories {
		buffer.WriteString("<b>" + story.Title + "</b><br/>")
		buffer.WriteString(story.Abstract + "<br/>")
		buffer.WriteString(
			"<a href=\"" + story.Url + "\"/>")
		if i != len(stories) - 1 {
			buffer.WriteString("<br/><br/>")
		}
		if i > 5 {
			break
		}
	}
	buffer.WriteString("</messageML>")
	client.SendMessageMLMessage(message.StreamId, buffer.String())
}
