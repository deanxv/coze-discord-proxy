package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
)

var NotifyTelegramBotToken = os.Getenv("NOTIFY_TELEGRAM_BOT_TOKEN")
var NotifyTelegramUserId = os.Getenv("NOTIFY_TELEGRAM_USER_ID")

var NotifyTelegramUserIdInt64 int64

var TgBot *tgbotapi.BotAPI

func InitTelegramBot() (err error) {

	TgBot, err = tgbotapi.NewBotAPI(NotifyTelegramBotToken)
	if err != nil {
		return err
	}
	TgBot.Debug = false
	return nil
}

func SendMessage(chattable tgbotapi.Chattable) (err error) {
	_, err = TgBot.Send(chattable)
	if err != nil {
		return err
	}
	return nil
}
