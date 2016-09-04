package commands

import (
	"regexp"
	"github.com/SymphonyOSF/botexample/client"
)

type CommandHandler interface {
	MatchRegex() *regexp.Regexp
	OnMessage(client.V2Message, client.BotClient, []CommandHandler)
	Help() string
}
