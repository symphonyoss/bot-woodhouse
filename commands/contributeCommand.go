package commands

import (
	"regexp"
	"github.com/SymphonyOSF/botexample/client"
	"bytes"
)

type ContributeCommand struct{}

func (cc ContributeCommand) MatchRegex() (regex *regexp.Regexp) {
	regex = regexp.MustCompile(`(?i)/contribute`)
	return regex
}

func (hc ContributeCommand) Help() (help string) {
	return "<b>/contribute</b> - find out how you can add your own commands"
}

func (cc ContributeCommand) OnMessage(message client.V2Message, client client.BotClient, handlers []CommandHandler) {
	var buffer bytes.Buffer
	buffer.WriteString("<messageML>")
	buffer.WriteString("Go to <a href=\"https://github.com/SymphonyOSF/botexample\"/>")
	buffer.WriteString("</messageML>")
	client.SendMessageMLMessage(message.StreamId, buffer.String())
}
