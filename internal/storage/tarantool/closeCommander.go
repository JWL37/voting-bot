package tarantool

import (
	"errors"
	"fmt"

	"github.com/tarantool/go-tarantool"
)

const (
	ErrPollNotFound      = "Ошибка: Голосование с таким ID не найдено"
	ErrNotCreator        = "Ошибка: Только создатель голосования может его завершить"
	ErrPollAlreadyClosed = "Ошибка: Это голосование уже завершено"
	ErrUpdatePollStatus  = "Ошибка: Не удалось обновить статус голосования"
)

func (s *Storage) ClosePoll(pollID uint64, userID string) error {
	const op = "storage.tarantool.ClosePoll"

	resp, err := s.DB.Select("polls", "primary", 0, 1, tarantool.IterEq, []interface{}{pollID})
	if err != nil || len(resp.Tuples()) == 0 {
		return errors.New(ErrPollNotFound)
	}

	tuple := resp.Tuples()[0]
	creatorID := tuple[2].(string)
	active := tuple[3].(bool)

	if creatorID != userID {
		return errors.New(ErrNotCreator)
	}

	if !active {
		return errors.New(ErrPollAlreadyClosed)
	}

	_, err = s.DB.Replace("polls", []interface{}{pollID, tuple[1], creatorID, false})
	if err != nil {
		s.log.Error(fmt.Errorf("%s: %v", op, err).Error())
		return errors.New(ErrUpdatePollStatus)
	}

	return nil
}
