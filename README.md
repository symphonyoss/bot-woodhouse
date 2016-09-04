# Woodhouse
![Woodhouse](http://i.imgur.com/1cESlQ2.png)

This is a sample Symphony bot written in Go. It authenticates as bot.user1 on nexus.symphony.com and then listens for messages. If you add bot.user1 to a room while this program is running, it will say hello. Try asking it "what time is it"

#Prerequisites
1. You must be running on Symphony's VPN
2. You must have go installed: https://golang.org/

# How to run this bot
1. clone this into your ~/go/src/github.com/SymphonyOSF 
```git clone git@github.com:SymphonyOSF/botexample.git```
2. for Mac users, add this to your .bash_profile ```export GOPATH=/Users/{YOUR USERNAME HERE}/go```
3. Then open a new terminal window so that it takes effect
4. Change to the botexample directory
```cd /src/github.com/SymphonyOSF/botexample```
5. Download all dependencies
```go get```
6. Run the program
```go run *.go resources/nexus.json```
7. Log on to nexus.symphony.com and add bot.user1 to a chatroom

# How to enable the /news feature
1. go to http://developer.nytimes.com/
2. register for an API key for "Top Stories V2"
3. Add the nytApiKey to your environment's config file
3. Run with that specified config file

# Adding a new environment (and using your own certificate)

First, you need a valid keystore, and you need to know the password to that keystore.

Then, you need to create a public certificate pem file from that. Use:

```openssl pkcs12 -in bot.userX.p12 -nokeys -out bot.userX-cert.pem```

Now, export the private key from the keystore:
```openssl pkcs12 -in bot.userX.p12 -nocerts -out bot.userX-key.pem```

Finally, convert the private key into a plain key format that is usable by Go, Python, etc.
```openssl rsa -in bot.userX-key.pem -out bot.userX-plainkey.pem```

Good! Now you need the bot.userX-cert.pem and the bot.userX-plainkey.pem to use in your environment config like this corporate.json config file:
```javascript
{
  "agentUrl": "https://sym-corp-stage-guse1-aha1.symphony.com:8444/agent",
  "sessionAuthUrl": "https://sym-corp-stage-guse1-aha1.symphony.com:8444/sessionauth/v1/authenticate",
  "keyManagerAuthUrl": "https://sym-corp-stage-guse1-aha1.symphony.com:8444/keyauth/v1/authenticate",
  "podUrl": "https://corporate.symphony.com",
  "keyFilePath": "private/corp/plain-bot.userX-plainkey.pem",
  "certFilePath": "private/corp/bot.userX-cert.pem",
  "nytApiKey": ""
}
```

I can now start Woodhouse using
```go run Hello.go private/corp/corporate.json```
