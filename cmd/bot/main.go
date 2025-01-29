package main

import (
	"log"

	"github.com/memuraFath/pocket__tg/pkg/config"
	"github.com/memuraFath/pocket__tg/pkg/repository"
	"github.com/memuraFath/pocket__tg/pkg/repository/boltdb"
	"github.com/memuraFath/pocket__tg/pkg/server"
	"github.com/memuraFath/pocket__tg/pkg/telegram"

	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
)

func main() {

	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg.TelegramToken)
	botApi, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	botApi.Debug = true

	pocketClient, err := pocket.NewClient(cfg.PocketConsumerKey)
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// log.Println("========================")
	// log.Println(cfg)
	// log.Println(cfg.Messages)
	// log.Println("========================")

	TokenRepository := boltdb.NewTokenRepository(db)
	bot := telegram.NewBot(botApi, pocketClient, cfg.AuthServerURL, TokenRepository, cfg)

	redirectServer := server.NewAuthprizationServer(pocketClient, TokenRepository, cfg.TelegramBotUrl)

	log.Printf("Authorized on account %s", botApi.Self.UserName)
	go func() {
		err := bot.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if err = redirectServer.Start(); err != nil {
		log.Fatal(err)
	}
}

func initDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DBPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	if err := db.Batch(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessToken))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestToken))
		return err

	}); err != nil {
		return nil, err
	}
	return db, err
}
