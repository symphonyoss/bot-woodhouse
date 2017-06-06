package nytimes

import (
	"gopkg.in/resty.v0"
	"fmt"
	"encoding/json"
)

var newsClient = resty.New()

type News struct {
	APIKey string
}

type Result struct {
	Title string
	Url string
	Abstract string
}

type NewsResult struct {
	Results []Result
}

func (news *News) TopStories() (result []Result ) {
	resp, err := newsClient.R().
	Get("https://api.nytimes.com/svc/topstories/v2/home.json?apikey=" + news.APIKey)
	if err != nil {
		fmt.Println(err)
		fmt.Println(resp)
		return
	} else {
		var newsResult NewsResult
		if err := json.Unmarshal(resp.Body(), &newsResult); err != nil {
			fmt.Println(err)
			return
		}
		result = newsResult.Results
		fmt.Println(resp.String())
	}
	return
}
