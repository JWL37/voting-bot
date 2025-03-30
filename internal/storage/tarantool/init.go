package tarantool

import (
	"fmt"
	"log/slog"
	"voting-bot/internal/config"

	"github.com/tarantool/go-tarantool"
)

const UnableToConnectDatabase = "Ошибка подключения к Tarantool: "

type Storage struct {
	DB  *tarantool.Connection
	log *slog.Logger
}

func NewConnection(log *slog.Logger, cfg *config.DatabaseConfig) (*Storage, error) {
	const op = "storage.tarantool.NewConnection"

	tarantoolConn, err := tarantool.Connect(cfg.Server, tarantool.Opts{
		User: cfg.User,
		Pass: cfg.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %s: %w", op, UnableToConnectDatabase, err)
	}

	storage := &Storage{
		DB:  tarantoolConn,
		log: log,
	}

	return storage, nil
}
