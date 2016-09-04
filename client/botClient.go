package client

import (
	"fmt"
	"encoding/json"
	"gopkg.in/resty.v0"
	"log"
	"crypto/tls"
	"sync"
)

type BotClient struct {
	AgentUrl          string
	SessionAuthUrl    string
	KeyManagerAuthUrl string
	PodUrl            string
	CertFilePath      string
	KeyFilePath       string
	isAuthenticated   bool
	skey              string
	kmsessionKey      string
}

func (botclient *BotClient) Authenticate() {
	loadCerts(resty.DefaultClient, botclient.CertFilePath, botclient.KeyFilePath)
	fmt.Println("\nauthenticating to " + botclient.SessionAuthUrl)
	botclient.setSkey(doCertificateAuthentication(botclient.SessionAuthUrl))
	fmt.Printf("Retrieved pod sessionKey=%s\n", botclient.skey)

	fmt.Println("\nauthenticating to " + botclient.KeyManagerAuthUrl)
	botclient.setKmsessionKey(doCertificateAuthentication(botclient.KeyManagerAuthUrl))
	fmt.Printf("Retrieved keymanager sessionKey=%s\n", botclient.kmsessionKey)
}

func (botClient *BotClient) setSkey(skey string) {
	botClient.skey = skey
}

func (botClient BotClient) Skey() string {
	return botClient.skey
}

func (botClient BotClient) KmsessionKey() string {
	return botClient.kmsessionKey
}

func (botClient *BotClient) setKmsessionKey(kmsessionKey string) {
	botClient.kmsessionKey = kmsessionKey
}

func doCertificateAuthentication(authUrl string) (sessionToken string) {
	resp, err := resty.R().
	Post(authUrl)

	if err != nil {
		fmt.Print(err)
	} else {

		var dat PodSession

		if err := json.Unmarshal(resp.Body(), &dat); err != nil {
			panic(err)
		}
		sessionToken = dat.Token
	}
	return
}

type PodSession struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

type SessionInfo struct {
	UserId int64
}

type V2MessageSubmission struct {
	Format  string `json:"format"`
	Message string `json:"message"`
}

func (botClient BotClient) SendMessage(streamId, message string, format string) {
	var messageSendUrl = botClient.AgentUrl + "/v2/stream/" + streamId + "/message/create"

	fmt.Println("send message to " + messageSendUrl)
	_, err := resty.R().
	SetHeader("Accept", "application/json").
	SetHeader("Content-Type", "application/json").
	SetHeader("sessionToken", botClient.Skey()).
	SetHeader("keyManagerToken", botClient.KmsessionKey()).
	SetBody(V2MessageSubmission{Format:format, Message:message}).
	Post(messageSendUrl)
	fmt.Println("Message sent")

	if err != nil {
		fmt.Print(err)
	}
}

func (botClient BotClient) SendPlainTextMessage(streamId, message string) {
	botClient.SendMessage(streamId, message, "TEXT")
}

func (botClient BotClient) SendMessageMLMessage(streamId, message string) {
	botClient.SendMessage(streamId, message, "MESSAGEML")
}

func (botClient BotClient) GetCurrentUserId() (userId int64) {
	resp, err := resty.R().
	SetHeader("Accept", "application/json").
	SetHeader("sessionToken", botClient.Skey()).
	Get(botClient.PodUrl + "/pod/v1/sessioninfo")

	if err != nil {
		fmt.Print(err)
	}
	var sessionInfo SessionInfo
	if err := json.Unmarshal(resp.Body(), &sessionInfo); err != nil {
		log.Fatal(err)
	}
	userId = sessionInfo.UserId
	return
}

func loadCerts(r *resty.Client, certFile, keyFile string) {
	cert1, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("ERROR client certificate: %s", err)
	}

	r.SetCertificates(cert1)
}

func (botClient BotClient) StartStreaming(wg sync.WaitGroup) <-chan interface{} {

	fmt.Println("\nCreating datafeed")

	resp, err := resty.R().
	SetHeader("Accept", "application/json").
	SetHeader("sessionToken", botClient.Skey()).
	SetHeader("keyManagerToken", botClient.KmsessionKey()).
	Post(botClient.AgentUrl + "/v1/datafeed/create")

	channel := make(chan interface{})

	if err != nil {
		fmt.Print(err)
	} else {
		var dat Datafeed

		if err := json.Unmarshal(resp.Body(), &dat); err != nil {
			panic(err)
		}
		fmt.Println(resp.String())
		datafeedId := dat.Id
		wg.Add(1)
		go botClient.readDatafeed(datafeedId, channel, wg);
	}

	return channel
}

func (botClient BotClient) readDatafeed(datafeedId string, channel chan <- interface{}, wg sync.WaitGroup) {
	defer wg.Done()
	streamingRestyClient := resty.New()
	loadCerts(streamingRestyClient, botClient.CertFilePath, botClient.KeyFilePath)
	for {
		fmt.Println("making read request for datafeed=" + datafeedId)
		resp, err := streamingRestyClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("sessionToken", botClient.skey).
		SetHeader("keyManagerToken", botClient.kmsessionKey).
		Get(botClient.AgentUrl + "/v2/datafeed/" + datafeedId + "/read")

		if err != nil {
			fmt.Println(err)
			fmt.Println(resp)
			fmt.Println("*** Shutting down datafeed ***")
			return
		} else {
			handleDatafeedResponse(resp, channel)
		}
	}
}

func handleDatafeedResponse(resp *resty.Response, channel chan <- interface{}) {
	var messageList = make([]*json.RawMessage, 0)
	fmt.Println("Raw JSON received from Datafeed=" + resp.String())

	if err := json.Unmarshal(resp.Body(), &messageList); err != nil {
		if resp.StatusCode() != 200 && resp.StatusCode() != 204 {
			fmt.Printf("status code=%d\n", resp.StatusCode())
			panic(err)
		}
		if resp.StatusCode() == 204 {
			fmt.Println("No messages returned during this poll.")
		}
	} else {
		// this ugliness is here because of the deserialization of different types that is
		// required for the message types possible on the response
		for _, rawMessage := range messageList {
			handleRawMessage(rawMessage, channel)
		}
	}
}

func handleRawMessage(message *json.RawMessage, channel chan <- interface{}) {
	// first marshal the raw JSON to a string
	jsonString, err := json.Marshal(&message)
	if err != nil {
		panic(err)
	}
	// then we do our first unmarshalling to a basemessage in order to get the type
	var v2BaseMessage V2BaseMessage = getV2BaseMessage(jsonString)

	switch v2BaseMessage.V2messageType {
	case "V2Message":
		handleV2Message(jsonString, v2BaseMessage, channel)
	case "UserJoinedRoomMessage":
		handleUserJoinedRoomMessage(jsonString, v2BaseMessage, channel)
	case "UserLeftRoomMessage":
		handleUserLeftRoomMessage(jsonString, v2BaseMessage, channel)
	case "RoomMemberDemotedFromOwnerMessage":
		handleRoomMemberDemotedFromOwnerMessage(jsonString, v2BaseMessage, channel)
	case "RoomMemberPromotedToOwnerMessage":
		handleRoomMemberPromotedToOwnerMessage(jsonString, v2BaseMessage, channel)
	default:
		fmt.Printf("unknown message type: %q", v2BaseMessage.V2messageType)
	}
}

func handleV2Message(jsonString []byte, v2BaseMessage V2BaseMessage, channel chan <- interface{}) {
	var v2Message V2Message
	if err := json.Unmarshal(jsonString, &v2Message); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received V2Message %+v\n", v2Message)
	channel <- v2Message
}

func handleUserJoinedRoomMessage(jsonString []byte, v2BaseMessage V2BaseMessage, channel chan <- interface{}) {
	var userJoinedRoomMessage UserJoinedRoomMessage
	if err := json.Unmarshal(jsonString, &userJoinedRoomMessage); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received UserJoinedRoomMessage %+v\n", userJoinedRoomMessage)
	channel <- userJoinedRoomMessage
}

func handleUserLeftRoomMessage(jsonString []byte, v2BaseMessage V2BaseMessage, channel chan <- interface{}) {
	var userLeftRoomMessage UserLeftRoomMessage
	if err := json.Unmarshal(jsonString, &userLeftRoomMessage); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received UserLeftRoomMessage %+v\n", userLeftRoomMessage)
	channel <- userLeftRoomMessage
}

func handleRoomMemberPromotedToOwnerMessage(jsonString []byte, v2BaseMessage V2BaseMessage, channel chan <- interface{}) {
	var roomMemberPromotedToOwnerMessage RoomMemberPromotedToOwnerMessage
	if err := json.Unmarshal(jsonString, &roomMemberPromotedToOwnerMessage); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received RoomMemberPromotedToOwnerMessage %+v\n", roomMemberPromotedToOwnerMessage)
	channel <- roomMemberPromotedToOwnerMessage
}

func handleRoomMemberDemotedFromOwnerMessage(jsonString []byte, v2BaseMessage V2BaseMessage, channel chan <- interface{}) {
	var roomMemberDemotedFromOwnerMessage RoomMemberDemotedFromOwnerMessage
	if err := json.Unmarshal(jsonString, &roomMemberDemotedFromOwnerMessage); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("RoomMemberDemotedFromOwnerMessage %+v\n", roomMemberDemotedFromOwnerMessage)
	channel <- roomMemberDemotedFromOwnerMessage
}

func getV2BaseMessage(jsonString []byte) (v2BaseMessage V2BaseMessage) {
	var err = json.Unmarshal(jsonString, &v2BaseMessage)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("V2BaseMessage %+v\n", v2BaseMessage)
	return
}

type Datafeed struct {
	Id string `json:"id"`
}
