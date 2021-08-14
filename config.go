package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	RefreshRate struct {
		Currency int `json:"currency"`
		Messages int `json:"messages"`
		MainCycle int `json:"mainCycle"`
	} `json:"refreshRate"`
	TelegramToken string `json:"telegramToken"`
	UsersFilePath string `json:"usersFilePath"`
	SaveTheSentRates bool `json:"saveTheSentRates"`
}

func loadConfig(path string, target interface{}) error  {
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