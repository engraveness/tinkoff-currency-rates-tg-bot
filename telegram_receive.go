package main

import "fmt"

type TelegramResponse struct {
	Ok     bool                     `json:"ok"`
	Result []TelegramResponseResult `json:"result"`
}

type TelegramResponseResult struct {
	UpdateId      int                           `json:"update_id"`
	Message       TelegramResponseResultMessage `json:"message,omitempty"`
	EditedMessage TelegramResponseResultMessage `json:"edited_message,omitempty"`
}

type TelegramResponseResultMessage struct {
	MessageId int `json:"message_id"`
	Chat      struct {
		Id        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	} `json:"chat"`
	From struct {
		LanguageCode string `json:"language_code"`
	} `json:"from"`
	Text string `json:"text"`
}

type TelegramMessage struct {
	UpdateId int
	ChatId int
	Text string
	FirstName string
	LastName string
	LanguageCode string
}

func getTelegramMessages(token string, offset int) []TelegramMessage {
	var url = fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d", token, offset)

	telegramParcel := TelegramResponse{}

	err := getJson(url, &telegramParcel)
	if err != nil {
		fmt.Printf("Error getting telegram messages: %s", err)
		return nil
	}

	result := make([]TelegramMessage, len(telegramParcel.Result))

	for i := 0; i < len(telegramParcel.Result); i++ {
		var part = telegramParcel.Result[i]
		var messageSource = TelegramResponseResultMessage{}

		if part.Message.MessageId != 0 {
			messageSource = part.Message
		} else {
			messageSource = part.EditedMessage
		}

		result[i] = TelegramMessage{
			UpdateId: part.UpdateId,
			ChatId: messageSource.Chat.Id,
			Text: messageSource.Text,
			FirstName: messageSource.Chat.FirstName,
			LastName: messageSource.Chat.LastName,
			LanguageCode: messageSource.From.LanguageCode,
		}
	}

	return result
}