package handlers

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

type PollDeleter interface {
	DeletePoll(pollID uint64, userID string) error
}

func HandleDeleteCommand(log *slog.Logger, message string, userID string, repo PollDeleter) string {
	parts := strings.Fields(message)
	if len(parts) != 3 {
		return "Ошибка: Используйте формат poll delete <ID голосования>"
	}

	pollID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return "Ошибка: ID голосования должен быть числом"
	}

	err = repo.DeletePoll(pollID, userID)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("Голосование %d и все связанные с ним голоса удалены!", pollID)
}
