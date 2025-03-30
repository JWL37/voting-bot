package handlers

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"voting-bot/internal/models"
)

type PollCreater interface {
	NextPollID() (uint64, error)
	SavePoll(id uint64, options map[string]models.Option, creatorID string) error
}

const (
	ErrInvalidFormat      = "Ошибка: Используйте формат poll create \"Вопрос?\" \"Вариант 1\" \"Вариант 2\" ..."
	ErrPollCreationFailed = "Ошибка: Не удалось создать голосование"
	SuccessPollCreated    = "Голосование создано!\nID: %d\nВопрос: %s\n"
)

func HandleCreateCommand(log *slog.Logger, message string, creatorID string, repo PollCreater) string {
	const op = "handlers.HandleCreateCommand"

	templog := log.With(
		slog.String("op", op),
		slog.String("creator_id", creatorID),
	)

	re := regexp.MustCompile(`"([^"]+)"`)
	parts := re.FindAllStringSubmatch(message, -1)
	if len(parts) < 2 {
		templog.Error(fmt.Errorf("%s: Error: %s", op, ErrInvalidFormat).Error())
		return ErrInvalidFormat
	}

	question := parts[0][1]
	options := make([]string, len(parts)-1)
	for i, part := range parts[1:] {
		options[i] = part[1]
	}

	id, err := repo.NextPollID()
	if err != nil {
		templog.Error(fmt.Errorf("%s: Error: %w", op, err).Error())

		return ErrPollCreationFailed
	}

	pollOptions := make(map[string]models.Option)
	for i, text := range options {
		pollOptions[strconv.Itoa(i+1)] = models.Option{
			TextOption: text,
			Votes:      0,
		}
	}

	err = repo.SavePoll(id, pollOptions, creatorID)
	if err != nil {
		templog.Error(fmt.Errorf("%s: Error: %w", op, err).Error())

		return ErrPollCreationFailed
	}

	response := fmt.Sprintf(SuccessPollCreated, id, question)
	for i, option := range options {
		response += fmt.Sprintf("%d. %s\n", i+1, option)
	}

	return response
}
