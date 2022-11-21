package counter

import (
	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/errors"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/typex"
	"github.com/garyburd/redigo/redis"
)

func keys(conn redis.Conn, batchSize int) ([]bs.Key, error) {
	var (
		cursor = 0
		result = []bs.Key{}
	)

	// Iterate the whole list.
	for {
		values, err := redis.Values(conn.Do("SCAN", cursor, "MATCH", prefix+"*", "COUNT", batchSize))
		if err != nil {
			return result, err
		}

		if n := len(values); n != 2 {
			return result, typex.Errorf(errors.Source, errors.UnexpectedResults,
				"Received %d values from redis, expected 2", n)
		}

		newCursor, err := redis.Int(values[0], nil)
		if err != nil {
			return result, err
		}

		keys, err := redis.Strings(values[1], nil)
		if err != nil {
			return result, err
		}

		// We only want insertions, not deletions
		for _, key := range keys {
			l := len(key) - insertSuffixLen
			if key[l:] == insertSuffix {
				// Remove the prefix
				result = append(result, bs.Key(key[prefixLen:l]))
			}
		}

		if newCursor <= 0 {
			break
		}
		cursor = newCursor
	}

	return result, nil
}

func members(conn redis.Conn, key bs.Key) ([]bs.Key, error) {
	m, err := redis.Strings(conn.Do("ZRANGE", prefix+key+insertSuffix, 0, -1))
	if err != nil {
		return nil, err
	}
	res := make([]bs.Key, 0, len(m))
	for _, v := range m {
		res = append(res, bs.Key(v))
	}
	return res, nil
}

func cardinality(conn redis.Conn, key bs.Key) ([]s.KeyCount, error) {
	res, err := redis.Int(conn.Do("ZCARD", prefix+key+insertSuffix))
	if err != nil {
		res = 0
	}
	return []s.KeyCount{
		s.KeyCount{Key: key, Count: res},
	}, err
}
