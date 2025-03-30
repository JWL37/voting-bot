package handlers

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

type voter interface {
	Vote(pollID uint64, userID string, optionNumber int) (string, error)
}

func HandleVoteCommand(log *slog.Logger, message string, userID string, repo voter) string {
	const op = "handlers.HandleVoteCommand"

	templog := log.With(
		slog.String("op", op),
		slog.String("user_id", userID),
	)

	parts := strings.Fields(message)
	if len(parts) != 4 {
		return "Ошибка: Используйте формат poll vote <ID голосования> <номер варианта>"
	}

	pollID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		templog.Error(fmt.Errorf(" Error: %w", err).Error())
		return "Ошибка: ID голосования должен быть числом"
	}
	optionNumber, err := strconv.Atoi(parts[3])
	if err != nil {
		templog.Error(fmt.Errorf(" Error: %w", err).Error())
		return "Ошибка: Номер варианта должен быть числом"
	}

	response, err := repo.Vote(pollID, userID, optionNumber)
	if err != nil {
		return err.Error()
	}

	return response
}
