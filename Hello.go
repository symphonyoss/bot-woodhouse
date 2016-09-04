package main

import (
	"fmt"
	"sync"
	"github.com/SymphonyOSF/botexample/client"
	"os"
	"regexp"
	"github.com/SymphonyOSF/botexample/commands"
	"github.com/SymphonyOSF/botexample/conf"
	"log"
)

// Globally accessible variables
var currentUserId int64
var currentHost string

// Required so that the main program thread wait sfor the datafeed to finish running before exiting
var wg sync.WaitGroup

var whatTimeIsItRegex = regexp.MustCompile(`(?i)what time is it`)
var botathonRegex = regexp.MustCompile(`(?i)#botathon`)

func main() {
	currentHost, _ = os.Hostname()
	fmt.Println("hello " + currentHost)

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		log.Fatal("Cannot start Woodhouse without telling him which environment you'd like to use. Your command line needs to look like \"go run Hello.go resources/nexus.json\"")
	}
	configurationLoader := conf.ConfigurationLoader{ConfigurationFileName:argsWithoutProg[0]}
	config := configurationLoader.Load(configurationLoader.ConfigurationFileName)
	botClient := client.BotClient{
		AgentUrl: config.AgentUrl,
		SessionAuthUrl: config.SessionAuthUrl,
		KeyManagerAuthUrl: config.KeyManagerAuthUrl,
		PodUrl: config.PodUrl,
		CertFilePath: config.CertFilePath,
		KeyFilePath: config.KeyFilePath,
	}
	botClient.Authenticate()

	currentUserId = botClient.GetCurrentUserId()

	messageHandlers := registerMessageHandlers(botClient, config)

	channel := botClient.StartStreaming(wg)
	for {
		message := <-channel
		switch message := message.(type) {
		case client.V2Message:
			if message.FromUserId == currentUserId {
				break
			}
			replyIfScotch(message.Message, message.StreamId, botClient)
			replyIfBotathon(message.Message, message.StreamId, botClient)
			for _, commandHandler := range messageHandlers {
				loc := commandHandler.MatchRegex().FindStringIndex(message.Message)
				if loc != nil && loc[0] == 11 {
					// first location after messageML tag
					commandHandler.OnMessage(message, botClient, messageHandlers)
				}
			}
			break
		case client.UserJoinedRoomMessage:
			if message.MemberAddedUserId == currentUserId {
				sendHelloRoomMessage(message.StreamId, botClient)
			}
			break
		}
	}
	wg.Wait()
}

func registerMessageHandlers(botclient client.BotClient, config conf.Configuration) ([]commands.CommandHandler) {
	messageHandlers := make([]commands.CommandHandler, 0)
	// You can get a developer key at http://developer.nytimes.com/
	if config.NytApiKey != "" {
		messageHandlers = append(messageHandlers, commands.NewsCommand{ApiKey:config.NytApiKey});
	}
	messageHandlers = append(messageHandlers, commands.HelpCommand{})
	messageHandlers = append(messageHandlers, commands.ContributeCommand{})
	messageHandlers = append(messageHandlers, commands.StockCommand{})
	messageHandlers = append(messageHandlers, commands.AboutCommand{})

	return messageHandlers
}

func replyIfScotch(message, streamId string, client client.BotClient) {
	if whatTimeIsItRegex.FindString(message) != "" {
		fmt.Println("replying back that it's time for scotch in streamId=" + streamId)
		client.SendPlainTextMessage(streamId, "Time for scotch!" + "\nsent from " + currentHost)

	}
}

func replyIfBotathon(message, streamId string, client client.BotClient) {
	if botathonRegex.FindString(message) != "" {
		fmt.Println("replying back that I'm going to win the botathon streamId=" + streamId)
		client.SendPlainTextMessage(streamId, "Have I won yet?")

	}
}

func sendHelloRoomMessage(streamId string, client client.BotClient) {
	fmt.Println("Sending hello room message for streamId=" + streamId)
	client.SendPlainTextMessage(streamId, "Hi! Thanks for adding me to the room." + "\nsent from " + currentHost)
}






