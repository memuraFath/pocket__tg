package telegram

import (
	"context"
	"log"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
)

const (
	comandStart = "start"
)

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) error {

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
			continue
		}

		if err := b.handleMessage(update.Message); err != nil {
			b.handleError(update.Message.Chat.ID, err)
		}
	}
	return nil
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	accessToken, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return errUnathorized
	}
	if err := b.saveLink(message, accessToken); err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.LinkSaved)
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case comandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.AlreadyAuthorized)
	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {

		return b.initAuthorizationProcess(message)
	}
	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.Responses.UnknownCommand)
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) saveLink(message *tgbotapi.Message, accessToken string) error {
	if err := b.validateURL(message.Text); err != nil {
		return errInvalidUrl
	}

	err := b.pocketClient.Add(context.Background(), pocket.AddInput{
		URL:         message.Text,
		AccessToken: accessToken,
	})
	if err != nil {
		return errInvalidUrl
	}
	return nil
}
func (b *Bot) validateURL(text string) error {
	_, err := url.ParseRequestURI(text)
	return err
}

/*func (b *Bot) getAccessToken(chatId int64) (string, error) {
	return b.TokenRepository.GetToken(chatId, repository.AccessToken)
}*/
