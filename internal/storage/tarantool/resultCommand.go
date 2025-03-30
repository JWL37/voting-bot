package tarantool

import (
	"voting-bot/internal/models"

	"github.com/tarantool/go-tarantool"
)

func (s *Storage) GetPoll(pollID uint64) (map[string]models.Option, error) {
	resp, err := s.DB.Select("polls", "primary", 0, 1, tarantool.IterEq, []interface{}{pollID})
	if err != nil || len(resp.Tuples()) == 0 {
		return nil, err
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

	return pollOptions, nil
}
