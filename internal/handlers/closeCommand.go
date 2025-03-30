package handlers

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

type pollCloser interface {
	ClosePoll(pollID uint64, userID string) error
}

const (
	errFormatMessage = "Ошибка: Используйте формат poll close <ID голосования>"
	errIDMessage     = "Ошибка: ID голосования должен быть числом"
	successMessage   = "Голосование %d завершено!"
)

func HandleCloseCommand(log *slog.Logger, message string, userID string, repo pollCloser) string {
	const op = "handlers.HandleCloseCommand"

	templog := log.With(
		slog.String("op", op),
		slog.String("user_id", userID),
	)

	parts := strings.Fields(message)
	if len(parts) != 3 {
		return errFormatMessage
	}

	pollID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		templog.Error(fmt.Errorf(" Error: %w", err).Error())
		return errIDMessage
	}

	err = repo.ClosePoll(pollID, userID)
	if err != nil {
		templog.Error(fmt.Errorf(" Error: %w", err).Error())
		return err.Error()
	}

	return fmt.Sprintf(successMessage, pollID)
}
