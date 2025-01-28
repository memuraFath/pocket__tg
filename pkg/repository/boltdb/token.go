package boltdb

import (
	"errors"
	"log"
	"pocket_tg/pkg/repository"
	"strconv"

	bolt "github.com/boltdb/bolt"
)

type TokenRepository struct {
	db *bolt.DB
}

func NewTokenRepository(db *bolt.DB) *TokenRepository {
	return &TokenRepository{
		db: db,
	}
}

func (r *TokenRepository) SaveToken(chatId int64, token string, bucket repository.Bucket) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(intToByte(chatId), []byte(token))
	})
}

func (r *TokenRepository) GetToken(chatId int64, bucket repository.Bucket) (string, error) {
	token := ""
	err := r.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(bucket))

		token = string(b.Get(intToByte(chatId)))

		log.Println(token)

		return nil
	})
	if err != nil {

		return "", err
	}

	if token == "" {
		return "", errors.New("token is empty")
	}
	return token, err
}

func intToByte(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}
