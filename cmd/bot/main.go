package main

import (
	"log"
	"os"
	"pocket_tg/pkg/config"
	"pocket_tg/pkg/repository"
	"pocket_tg/pkg/repository/boltdb"
	"pocket_tg/pkg/server"
	"pocket_tg/pkg/telegram"

	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
)

func main() {

	os.Setenv("TOKEN", "7538467353:AAHvKXCiLG7Es_FZffLK0MZ0BcYmnFG2EUk")
	os.Setenv("CONSUMER_KEY", "113171-ed6f3927e3622fd715be0d9")
	os.Setenv("AUTH_SERVER_URL", "http://localhost:80")
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
