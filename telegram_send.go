package main

import (
	"encoding/json"
	"fmt"
)

type TelegramSendMessageRequest struct {
	ChatId              int    `json:"chat_id"`
	Text                string `json:"text"`
	DisableNotification   bool   `json:"disable_notification"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
}

type TelegramEditMessageRequest struct {
	MessageId int `json:"message_id"`
	TelegramSendMessageRequest
}

type TelegramMessageResponse struct {
	Ok        bool `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
	Result struct {
		MessageId   int    `json:"message_id"`
	} `json:"result"`
}

func sendTelegramMessage(token string, chatId int, message string) TelegramMessageResponse {
	var url = fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	request := TelegramSendMessageRequest{
		ChatId: chatId,
		Text: message,
		DisableNotification: true,
		ParseMode: "HTML",
		DisableWebPagePreview: true,
	}

	requestBytes, _ := json.Marshal(request)

	telegramSendMessageParcel := TelegramMessageResponse{}

	err := postJson(url, &telegramSendMessageParcel, requestBytes)
	if err != nil {
		return TelegramMessageResponse{}
	}

	return telegramSendMessageParcel
}

func editTelegramMessage(token string, chatId int, messageId int, message string) TelegramMessageResponse {
	var url = fmt.Sprintf("https://api.telegram.org/bot%s/editMessageText", token)

	request := TelegramEditMessageRequest{
		MessageId: messageId,
		TelegramSendMessageRequest: TelegramSendMessageRequest {
			ChatId: chatId,
			Text: message,
			DisableNotification: true,
			ParseMode: "HTML",
			DisableWebPagePreview: true,
		},
	}

	requestBytes, _ := json.Marshal(request)

	telegramSendMessageParcel := TelegramMessageResponse{}

	postJson(url, &telegramSendMessageParcel, requestBytes)

	return telegramSendMessageParcel
}
