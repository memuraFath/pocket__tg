package repository

const (
	AccessToken  Bucket = "access_token"
	RequestToken Bucket = "request_token"
)

type Bucket string
type TokenRepository interface {
	SaveToken(chatId int64, token string, bucket Bucket) error
	GetToken(chatId int64, bucket Bucket) (string, error)
}
