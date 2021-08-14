package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const euroTicker = "EURRUB"
const usdTicker = "USDRUB"

func main() {
	var configPath = "config.json"

	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	var config = Config{}

	err := loadConfig(configPath, &config)
	if err != nil {
		fmt.Printf("Can't load config at path %s: %s. Pass the full path to the config.json as the first argument or create config.json near the executable file", configPath, err)
		return
	}

	var token, status = getToken()
	if !status {
		fmt.Println("Cant get token!")
		return
	}

	nextChatCheck := time.Now().Unix()
	nextRateCheck := nextChatCheck

	var prevUsd float32 = 0
	var prevEur float32 = 0

	var rateForSend = ""
	var newRateForSend = ""

	var messagesMap = map[int]int{}

	var mapPath = config.UsersFilePath
	if mapPath == "" {
		mapPath = "map.json"
	}

	// we're good if config isn't found, we'll create it ourselves
	_ = loadConfig(mapPath, &messagesMap)

	var offset = 0

	var refreshRate = time.Duration(config.RefreshRate.MainCycle) * time.Second

	fmt.Println("Starting main cycle")

	for {
		time.Sleep(refreshRate)

		var now = time.Now()
		var unix = now.Unix()
		if nextRateCheck - unix <= 0 {
			result := getExchangeRate(token)

			newRateForSend = fmt.Sprintf("%s\n%s", getCurrencyString(euroTicker, result[euroTicker], prevEur), getCurrencyString(usdTicker, result[usdTicker], prevUsd))

			prevEur = result[euroTicker]
			prevUsd = result[usdTicker]

			nextRateCheck = unix + int64(config.RefreshRate.Currency)
		}

		var needSaveMap = false

		if nextChatCheck - unix <= 0 {
			messages := getTelegramMessages(config.TelegramToken, offset)

			for _, message := range messages {
				switch message.Text {
				case "/stop":
					delete(messagesMap, message.ChatId)
				default:
					// since we are editing the last sent message, we want our message to be the last one in the conversation
					// so, we send a new message as response to any other user message
					messagesMap[message.ChatId] = 0
				}

				needSaveMap = true

				offset = message.UpdateId + 1
			}
		}

		if newRateForSend != rateForSend {
			var rateWithDate = fmt.Sprintf("%s\nUpdated at %s", newRateForSend, now.Format("2006-01-02 15:04:05"))

			for k, v := range messagesMap {
				if v != 0 {
					editTelegramMessage(config.TelegramToken, k, v, rateWithDate)
				} else {
					things := sendTelegramMessage(config.TelegramToken, k, rateWithDate)

					if things.Ok{
						messagesMap[k] = things.Result.MessageId
						needSaveMap = true
					}
				}
			}

			rateForSend = newRateForSend
		}

		if needSaveMap {
			file, _ := json.MarshalIndent(messagesMap, "", " ")

			err = ioutil.WriteFile("map.json", file, 0644)
			if err != nil {
				fmt.Printf("error during messagesMap save: %s", err)
			}
		}
	}
}
