package conf

import (
	"encoding/json"
	"os"
	"fmt"
)

type Configuration struct {
	AgentUrl    string
	SessionAuthUrl   string
	KeyManagerAuthUrl string
	PodUrl string
	KeyFilePath string
	CertFilePath string
	NytApiKey string
}

type ConfigurationLoader struct {
	ConfigurationFileName string
}

func (configurationLoader ConfigurationLoader) Load(configurationFile string)(configuration Configuration) {
	file, _ := os.Open(configurationFile)
	decoder := json.NewDecoder(file)
	configuration = Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("Configuration=%+v\n", configuration)
	return
}
