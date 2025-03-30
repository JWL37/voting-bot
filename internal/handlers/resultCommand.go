package handlers

import (
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"strconv"
	"strings"
	"voting-bot/internal/models"
)

type PollGettter interface {
	GetPoll(pollID uint64) (map[string]models.Option, error)
}

func HandleResultCommand(log *slog.Logger, message string, repo PollGettter) string {
	const op = "handlers.HandleResultCommand"

	templog := log.With(
		slog.String("op", op),
	)

	parts := strings.Fields(message)
	if len(parts) != 3 {
		return "Ошибка: Используйте формат poll results <ID голосования>"
	}

	pollID, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		templog.Error(fmt.Errorf(" Error: %w", err).Error())

		return "Ошибка: ID голосования должен быть числом"
	}

	pollOptions, err := repo.GetPoll(pollID)
	if err != nil {
		return "Ошибка: Голосование с таким ID не найдено"
	}

	response := fmt.Sprintf("Результаты голосования %d:\n", pollID)
	sortedKeysPollOptions := slices.Sorted(maps.Keys(pollOptions))

	for _, optionIdx := range sortedKeysPollOptions {
		response += fmt.Sprintf("%s. %s: %d голосов\n", optionIdx, pollOptions[optionIdx].TextOption, pollOptions[optionIdx].Votes)
	}

	return response
}
