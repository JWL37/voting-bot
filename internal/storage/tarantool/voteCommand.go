package tarantool

import (
	"errors"
	"fmt"
	"strconv"
	"voting-bot/internal/models"

	"github.com/tarantool/go-tarantool"
)

func (s *Storage) Vote(pollID uint64, userID string, optionNumber int) (string, error) {
	resp, err := s.DB.Select("polls", "primary", 0, 1, tarantool.IterEq, []interface{}{pollID})
	if err != nil || len(resp.Tuples()) == 0 {
		return "", errors.New("ошибка: Голосование с таким ID не найдено")
	}

	tuple := resp.Tuples()[0]
	rawPollOptions := tuple[1].(map[interface{}]interface{})
	pollOptions := make(map[string]models.Option)
	for k, v := range rawPollOptions {
		optionMap := v.(map[interface{}]interface{})
		pollOptions[k.(string)] = models.Option{
			TextOption: optionMap["TextOption"].(string),
			Votes:      int(optionMap["Votes"].(uint64)),
		}
	}
	creatorID := tuple[2].(string)
	active := tuple[3].(bool)

	if !active {
		return "", errors.New("ошибка: Это голосование уже завершено")
	}

	if optionNumber < 1 || optionNumber > len(pollOptions) {
		return "", errors.New("ошибка: Неверный номер варианта")
	}

	voterResp, err := s.DB.Select("voters", "primary", 0, 1, tarantool.IterEq, []interface{}{pollID, userID})
	if err != nil {
		return "", errors.New("ошибка: Не удалось проверить статус голосования")
	}
	if len(voterResp.Tuples()) > 0 {
		return "", errors.New("ошибка: Вы уже голосовали в этом голосовании")
	}

	selectedOption := pollOptions[strconv.Itoa(optionNumber)]
	selectedOption.Votes++
	pollOptions[strconv.Itoa(optionNumber)] = selectedOption
	_, err = s.DB.Replace("polls", []interface{}{pollID, pollOptions, creatorID, active})
	if err != nil {
		return "", errors.New("ошибка: Не удалось обновить голосование")
	}

	_, err = s.DB.Insert("voters", []interface{}{pollID, userID})
	if err != nil {
		return "", errors.New("ошибка: Не удалось сохранить информацию о голосовании")
	}

	return fmt.Sprintf("ваш голос за вариант \"%s\" в голосовании %d учтён!", selectedOption.TextOption, pollID), nil
}
