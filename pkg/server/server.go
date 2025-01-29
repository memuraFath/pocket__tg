package server

import (
	"context"
	"log"
	"net/http"
	"strconv"

	repository "github.com/memuraFath/pocket__tg/pkg/repository"

	"github.com/pkg/errors"
	pocket "github.com/zhashkevych/go-pocket-sdk"
)

type AuthorizationServer struct {
	server          *http.Server               // base server func
	pocketClient    *pocket.Client             // interchange with pocket API
	tokenRepository repository.TokenRepository // to use DB
	redirectUrl     string                     // base redirect url
}

func NewAuthprizationServer(pocketClient *pocket.Client, tokenRepository repository.TokenRepository, redirectURL string) *AuthorizationServer {
	return &AuthorizationServer{
		pocketClient:    pocketClient,
		tokenRepository: tokenRepository,
		redirectUrl:     redirectURL,
	}
}

func (s *AuthorizationServer) Start() error {
	s.server = &http.Server{
		Addr:    ":80",
		Handler: s,
	}
	return s.server.ListenAndServe()
}

func (s *AuthorizationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatIdQuery := r.URL.Query().Get("chat_id")
	if chatIdQuery == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatId, err := strconv.ParseInt(chatIdQuery, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = s.getAndSaveAccessToken(r.Context(), chatId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Add("Location", s.redirectUrl)
	w.WriteHeader(http.StatusMovedPermanently)
}

func (s *AuthorizationServer) getAndSaveAccessToken(ctx context.Context, chatId int64) error {
	requestToken, err := s.tokenRepository.GetToken(chatId, repository.RequestToken)
	if err != nil {
		return errors.WithMessage(err, "failed to get request token")
	}

	accessToken, err := s.pocketClient.Authorize(ctx, requestToken)
	if err != nil {
		return errors.WithMessage(err, "failed to get access token")
	}
	if err = s.tokenRepository.SaveToken(chatId, accessToken.AccessToken, repository.AccessToken); err != nil {

		return errors.WithMessage(err, "failed to save access token to DB")
	}
	log.Printf("chat_id:\t%d\nrequstToken:\t%s\naccessToken:\t%s", chatId, requestToken, accessToken)
	return nil
}
