package main

type Config struct {
	RefreshRate struct {
		Currency  int `json:"currency"`
		Messages  int `json:"messages"`
		MainCycle int `json:"mainCycle"`
	} `json:"refreshRate"`
	TelegramToken    string `json:"telegramToken"`
	UsersFilePath    string `json:"usersFilePath"`
	SaveTheSentRates bool   `json:"saveTheSentRates"`
}
