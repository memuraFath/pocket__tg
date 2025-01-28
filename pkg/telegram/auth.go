package telegram

import (
	"context"
	"fmt"
	"pocket_tg/pkg/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *Bot) generateAuthorizationLink(chatId int64) (string, error) {
	redirectUrl := b.generateRedirectUrl(chatId)
	requestToken, err := b.pocketClient.GetRequestToken(context.Background(), redirectUrl)
	if err != nil {
		return "", err
	}
	err = b.TokenRepository.SaveToken(chatId, requestToken, repository.RequestToken)
	if err != nil {
		return "", err
	}
	return b.pocketClient.GetAuthorizationURL(requestToken, redirectUrl)

}

func (b *Bot) generateRedirectUrl(chatId int64) string {
	return fmt.Sprintf("%s?chat_id=%d", b.redirectUrl, chatId)
}

func (b *Bot) getAccessToken(chatId int64) (string, error) {
	return b.TokenRepository.GetToken(chatId, repository.AccessToken)
}

func (b *Bot) initAuthorizationProcess(message *tgbotapi.Message) error {
	authLink, err := b.generateAuthorizationLink(message.Chat.ID)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	msg.Text = fmt.Sprintf(b.messages.Responses.Start, message.From.UserName, authLink)
	_, err = b.bot.Send(msg)
	return err
}
