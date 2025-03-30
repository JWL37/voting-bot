package tarantool

import (
	"fmt"
	"voting-bot/internal/models"
)

func (s *Storage) NextPollID() (uint64, error) {
	const op = "storage.tarantool.NextPollID"

	resp, err := s.DB.Call("box.sequence.poll_id_seq:next", []interface{}{})
	if err != nil || len(resp.Data) == 0 {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	pollID, ok := resp.Data[0].([]interface{})
	if !ok || len(pollID) == 0 {
		return 0, fmt.Errorf("%s: invalid pollID received: %w", op, err)
	}

	id, ok := pollID[0].(uint64)
	if !ok {
		return 0, fmt.Errorf("%s:  pollID is not uint64: %w", op, err)
	}

	return id, nil
}

func (s *Storage) SavePoll(id uint64, options map[string]models.Option, creatorID string) error {
	const op = "storage.tarantool.SavePoll"

	_, err := s.DB.Insert("polls", []interface{}{id, options, creatorID, true})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
