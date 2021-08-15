package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const euroTicker = "EURRUB"
const usdTicker = "USDRUB"

type User struct {
	MessageId   int     `json:"messageId"`
	Threshold   float64 `json:"threshold"`
	LastSentUsd float64 `json:"lastSentUsd"`
	LastSentEur float64 `json:"lastSentEur"`
}

func loadJsonFromFile(path string, target interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&target)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var configPath = "config.json"

	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	var config = Config{}

	err := loadJsonFromFile(configPath, &config)
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

	var prevUsd float64 = 0
	var prevEur float64 = 0

	var rateForSend = ""
	var newRateForSend = ""

	var messagesMap = map[int]*User{}

	var mapPath = config.UsersFilePath
	if mapPath == "" {
		mapPath = "map.json"
	}

	// we're good if config isn't found, we'll create it ourselves
	_ = loadJsonFromFile(mapPath, &messagesMap)

	var offset = 0

	var refreshRate = time.Duration(config.RefreshRate.MainCycle) * time.Second

	fmt.Println("Starting main cycle")

	for {
		time.Sleep(refreshRate)

		var now = time.Now()
		var unix = now.Unix()
		if nextRateCheck-unix <= 0 {
			result := getExchangeRate(token)

			var newEur = math.Floor(result[euroTicker]*100) / 100
			var newUsd = math.Floor(result[usdTicker]*100) / 100

			newRateForSend = fmt.Sprintf("%s\n%s", getCurrencyString(euroTicker, newEur, prevEur), getCurrencyString(usdTicker, newUsd, prevUsd))

			prevEur = newEur
			prevUsd = newUsd

			nextRateCheck = unix + int64(config.RefreshRate.Currency)
		}

		var needSaveMap = false

		if nextChatCheck-unix <= 0 {
			messages := getTelegramMessages(config.TelegramToken, offset)

			for _, message := range messages {
				if strings.HasPrefix(message.Text, "/threshold") {
					if message.Text == "/threshold" {
						messagesMap[message.ChatId].MessageId = 0

						sendTelegramMessage(config.TelegramToken, message.ChatId, fmt.Sprintf("Current threshold is %.2f. To change it send /threshold *VALUE*, i.e. /threshold 0.25", messagesMap[message.ChatId].Threshold))
					} else {
						var stripped = message.Text[len("/threshold "):]
						val, err := strconv.ParseFloat(stripped, 64)

						messagesMap[message.ChatId].MessageId = 0

						if err != nil {
							sendTelegramMessage(config.TelegramToken, message.ChatId, "Error processing argument. Please enter valid floating point value, i.e. /threshold 0.25")
						} else {
							messagesMap[message.ChatId].Threshold = val
							needSaveMap = true

							sendTelegramMessage(config.TelegramToken, message.ChatId, fmt.Sprintf("New threshold is %.2f", val))
						}
					}
				} else {
					switch message.Text {
					case "/stop":
						delete(messagesMap, message.ChatId)

						needSaveMap = true
					case "/start":
						var greetings = fmt.Sprintf("Welcome! from now on, you will receive a message with the current exchange rate for the euro and the dollar. If you want, you can set the minimum threshold by which the course should change, so that you know about it by typing the command /threshold *VALUE*")
						messagesMap[message.ChatId] = new(User)

						sendTelegramMessage(config.TelegramToken, message.ChatId, greetings)

						needSaveMap = true
					default:
						// since we are editing the last sent message, we want our message to be the last one in the conversation
						// so, we send a new message as response to any other user message
						messagesMap[message.ChatId].MessageId = 0
					}
				}

				offset = message.UpdateId + 1
			}
		}

		for k, v := range messagesMap {
			if v.Threshold != 0 {
				if math.Abs(v.LastSentUsd-prevUsd) < v.Threshold && math.Abs(v.LastSentEur-prevEur) < v.Threshold {
					continue
				}

				var message = fmt.Sprintf("%s\n%s", getCurrencyString(euroTicker, prevEur, v.LastSentEur), getCurrencyString(usdTicker, prevUsd, v.LastSentUsd))
				sendTelegramMessage(config.TelegramToken, k, message)

				v.LastSentUsd = prevUsd
				v.LastSentEur = prevEur

				if config.SaveTheSentRates {
					needSaveMap = true
				}
			} else {
				if newRateForSend != rateForSend || v.LastSentUsd == 0 {
					var message = fmt.Sprintf("%s\nUpdated at %s", newRateForSend, now.Format("2006-01-02 15:04:05"))

					if v.MessageId != 0 {
						newMessage := editTelegramMessage(config.TelegramToken, k, v.MessageId, message)

						if newMessage.Ok && config.SaveTheSentRates {
							needSaveMap = true
						}
					} else {
						newMessage := sendTelegramMessage(config.TelegramToken, k, message)

						if newMessage.Ok {
							v.MessageId = newMessage.Result.MessageId
							v.LastSentUsd = prevUsd
							v.LastSentEur = prevEur

							if config.SaveTheSentRates {
								needSaveMap = true
							}
						}
					}
				}
			}

			rateForSend = newRateForSend
		}

		if needSaveMap {
			file, _ := json.MarshalIndent(messagesMap, "", " ")

			err = ioutil.WriteFile(config.UsersFilePath, file, 0644)
			if err != nil {
				fmt.Printf("error during messagesMap save: %s", err)
			}
		}
	}
}
