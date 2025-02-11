package telegram

import (
	"errors"

	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	errInvalidUrl   = errors.New("url is invalid")
	errUnathorized  = errors.New("user is unathorized")
	errUnableToSave = errors.New("unable to save link")
)

func (b *Bot) handleError(chatId int64, err error) {
	var msgTxt string
	switch err {
	case errInvalidUrl:
		msgTxt = "Invalid url"
	case errUnathorized:
		msgTxt = "You're not authortised. \nUse /start comand to generate authorization link"
	case errUnableToSave:
		msgTxt = "Something went wrong.Try again later"
	default:
		msgTxt = "Unknown error occured"
	}
	msg := tgbotapi.NewMessage(chatId, msgTxt)
	log.WithFields(log.Fields{
		"handler": "telegram.handleError",
		"problem": msgTxt,
	}).Error(err)
	b.bot.Send(msg)
}
