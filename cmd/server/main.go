package main

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"image-previewer/internal/app"
)

var shaCommit = "local"

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := log.With().Str("sha_commit", shaCommit).Logger()

	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("Не удалось загрузить файл env")
	}

	srv, err := app.NewServer(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("Ошибка старта сервера")
	}

	srv.WithLogger(logger)

	listenPort, ok := viper.Get("LISTEN_PORT").(int)
	if !ok {
		listenPort = 8090
	}

	ctx := log.Logger.WithContext(context.Background())
	if err := srv.Listen(ctx, listenPort); err != nil {
		logger.Fatal().Err(err).Msg("Не удалось прослушать порт")
	}
}
