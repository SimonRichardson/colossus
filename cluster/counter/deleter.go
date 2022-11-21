package counter

import (
	t "github.com/SimonRichardson/colossus/cluster"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/garyburd/redigo/redis"
)

func deletion(conn redis.Conn, members []s.KeyFieldScoreTxnValue, sizeExpiry s.SizeExpiry) ([]s.KeyCount, error) {
	for _, member := range members {
		if err := sendDeleteScript(conn,
			member.Key,
			member.Field,
			member.Score,
			sizeExpiry.Size,
		); err != nil {
			return generateResult(members, 0), err
		}
	}

	if err := conn.Flush(); err != nil {
		return generateResult(members, 0), err
	}

	if !defaultVerifyResults {
		return generateResult(members, 1), nil
	}

	result := make([]s.KeyCount, 0, len(members))

	for _, m := range members {
		res, err := redis.Int(conn.Receive())
		if err != nil {
			return result, err
		}

		result = append(result, s.KeyCount{Key: m.Key, Count: abs(res)})
	}

	if len(result) < len(members) {
		return result, t.ErrPartialDeletions
	}

	return result, nil
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func generateResult(members []s.KeyFieldScoreTxnValue, count int) []s.KeyCount {
	result := make([]s.KeyCount, 0, len(members))
	for _, m := range members {
		result = append(result, s.KeyCount{Key: m.Key, Count: count})
	}
	return result
}
