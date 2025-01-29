package telegram

import (
	"github.com/memuraFath/pocket__tg/pkg/config"
	"github.com/memuraFath/pocket__tg/pkg/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	pocketClient    *pocket.Client
	redirectUrl     string
	TokenRepository repository.TokenRepository
	messages        config.Messages
}

func NewBot(bot *tgbotapi.BotAPI, pocketClient *pocket.Client, redirectUrl string, tr repository.TokenRepository, cfg *config.Config) *Bot {

	return &Bot{
		bot:             bot,
		pocketClient:    pocketClient,
		redirectUrl:     redirectUrl,
		TokenRepository: tr,
		messages:        cfg.Messages,
	}
}

func (b *Bot) Start() error {

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}
	b.handleUpdates(updates)

	return nil

}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}
