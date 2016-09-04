package commands

import (
	"regexp"
	"github.com/SymphonyOSF/botexample/client"
	"bytes"
	"os"
)

type AboutCommand struct{}

func (aboutCommand AboutCommand) MatchRegex() (regex *regexp.Regexp) {
	regex = regexp.MustCompile(`(?i)/about`)
	return regex
}

func (aboutCommand AboutCommand) Help() (help string) {
	return "<b>/about</b> - what is this?"
}

var currentHost,_ = os.Hostname()

func (aboutCommand AboutCommand) OnMessage(message client.V2Message, client client.BotClient, handlers []CommandHandler) {
	var buffer bytes.Buffer
	buffer.WriteString("<messageML>")
	buffer.WriteString("I live on " + currentHost + "<br/>")
	buffer.WriteString("I was written in Go. It's a language written by Google... <a href=\"https://golang.org\"/><br/>")
	buffer.WriteString("I use the API's from <a href=\"https://developers.symphony.com\"/> to authenticate, listen for messages, and reply.<br/>")
	buffer.WriteString("</messageML>")

	client.SendMessageMLMessage(message.StreamId, buffer.String())
}
