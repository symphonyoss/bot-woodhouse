package commands

import (
	"regexp"
	"github.com/SymphonyOSF/botexample/client"
	"bytes"
)

type HelpCommand struct{}

func (helpCommand HelpCommand) MatchRegex() (regex *regexp.Regexp) {
	regex = regexp.MustCompile(`(?i)/help`)
	return regex
}

func (helpCommand HelpCommand) Help() (help string) {
	return "<b>/help</b> - list of commands that are currently enabled"
}

func (helpCommand HelpCommand) OnMessage(message client.V2Message, client client.BotClient, handlers []CommandHandler) {
	var buffer bytes.Buffer
	buffer.WriteString("<messageML>")

	buffer.WriteString("Here are the things I can do:<br/>")
	for _, commandHandler := range handlers {
		buffer.WriteString(commandHandler.Help() + "<br/>")

	}
	buffer.WriteString("</messageML>")

	client.SendMessageMLMessage(message.StreamId, buffer.String())
}
