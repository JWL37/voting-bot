package tarantool

import (
	"errors"
	"log"

	"github.com/tarantool/go-tarantool"
)

func (s *Storage) DeletePoll(pollID uint64, userID string) error {
	resp, err := s.DB.Select("polls", "primary", 0, 1, tarantool.IterEq, []interface{}{pollID})
	if err != nil || len(resp.Tuples()) == 0 {
		return errors.New("Ошибка: Голосование с таким ID не найдено")
	}

	tuple := resp.Tuples()[0]
	creatorID := tuple[2].(string)

	if creatorID != userID {
		return errors.New("Ошибка: Только создатель голосования может его удалить")
	}

	_, err = s.DB.Delete("polls", "primary", []interface{}{pollID})
	if err != nil {
		return errors.New("Ошибка: Не удалось удалить голосование")
	}

	_, err = s.DB.Delete("voters", "poll_id_index", []interface{}{pollID})
	if err != nil {
		log.Printf("Ошибка при удалении записей о голосах: %v", err)
		return errors.New("Ошибка: Не удалось удалить записи о голосах")
	}

	return nil
}
