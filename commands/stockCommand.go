package commands

import (
	"regexp"
	"fmt"
	"bytes"
	"github.com/SymphonyOSF/botexample/client"
	"gopkg.in/resty.v0"
	"encoding/json"
	"time"
)

type StockCommand struct{}

func (sc StockCommand) MatchRegex() (regex *regexp.Regexp) {
	regex = regexp.MustCompile(`(?i)/stock`)
	return regex
}

func (sc StockCommand) Help() (help string) {
	return "<b>/stock SYMBOL</b> - provides data from Markit"
}

var stockRegex = regexp.MustCompile(`(?i)/stock `)

type StockResult struct {
	Name      string
	LastPrice float64
	Timestamp string
	Status    string
}

const referenceTimeFormat = "Mon Jan 2 15:04:05 UTC-06:00 2006"

func (stockCommand StockCommand) OnMessage(message client.V2Message, client client.BotClient, handlers []CommandHandler) {
	locStockStart := stockRegex.FindStringIndex(message.Message)
	if locStockStart == nil {
		fmt.Println("replying back with missingstock symbol in streamId=" + message.StreamId)
		client.SendPlainTextMessage(message.StreamId, "I'm sorry. You must provide a stock symbol to lookup")
	}
	symbol := message.Message[locStockStart[1]:(len(message.Message) - len("</messageML>"))]
	resp, err := resty.R().Get("http://dev.markitondemand.com/MODApis/Api/v2/Quote/json?symbol=" + symbol)
	var stockResult StockResult
	if err != nil {
		fmt.Println(err)
		fmt.Println(resp)
		return
	} else {
		if err := json.Unmarshal(resp.Body(), &stockResult); err != nil || stockResult.Status != "SUCCESS" {
			fmt.Println("replying back with bad stock symbol in streamId=" + message.StreamId)
			client.SendPlainTextMessage(message.StreamId, "I'm sorry. I couldn't find a stock matching " + symbol)
			return
		}
	}
	fmt.Println("replying back with stock price in streamId=" + message.StreamId)
	var buffer bytes.Buffer
	buffer.WriteString("<messageML>")
	buffer.WriteString("<b>" + stockResult.Name + "</b><br/>")
	prettyPrice := fmt.Sprintf("%.2f", stockResult.LastPrice)

	t, _ := time.Parse(referenceTimeFormat, stockResult.Timestamp)
	prettyTime := fmt.Sprint(t.Format(time.Kitchen) + " EST " + t.Format("Mon Jan 2"))

	buffer.WriteString("Last price $" + prettyPrice + "<br/>At " + prettyTime)
	buffer.WriteString("</messageML>")
	client.SendMessageMLMessage(message.StreamId, buffer.String())
}
