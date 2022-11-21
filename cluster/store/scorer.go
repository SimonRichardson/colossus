package store

import (
	"time"

	"github.com/SimonRichardson/colossus/errors"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/typex"
	"github.com/garyburd/redigo/redis"
)

func score(conn redis.Conn, members []s.KeyFieldTxnValue) (map[s.KeyFieldTxnValue]s.Presence, error) {
	for _, keyField := range members {
		var (
			key   = keyField.Key.String()
			field = keyField.Field.String()
		)
		if err := conn.Send("HGET", prefix+key+insertSuffix, field); err != nil {
			return map[s.KeyFieldTxnValue]s.Presence{}, err
		}

		if err := conn.Send("HGET", prefix+key+deleteSuffix, field); err != nil {
			return map[s.KeyFieldTxnValue]s.Presence{}, err
		}
	}

	if err := conn.Flush(); err != nil {
		return map[s.KeyFieldTxnValue]s.Presence{}, err
	}

	var (
		m   = map[s.KeyFieldTxnValue]s.Presence{}
		now = time.Now().UnixNano()
	)
	for i := 0; i < len(members); i++ {
		var (
			insertValue, insertErr = redis.String(conn.Receive())
			deleteValue, deleteErr = redis.String(conn.Receive())
		)

		switch {
		case insertErr == nil && deleteErr == redis.ErrNil:
			score, _, expiry, _, err := ExtractScoreTxnExpiryValue(insertValue)
			if err != nil {
				return map[s.KeyFieldTxnValue]s.Presence{}, err
			}
			m[members[i]] = s.Presence{
				Present:  expiry >= now,
				Inserted: true,
				Score:    score,
			}
		case insertErr == redis.ErrNil && deleteErr == nil:
			score, _, expiry, _, err := ExtractScoreTxnExpiryValue(deleteValue)
			if err != nil {
				return map[s.KeyFieldTxnValue]s.Presence{}, err
			}
			m[members[i]] = s.Presence{
				Present:  expiry >= now,
				Inserted: false,
				Score:    score,
			}
		case insertErr == redis.ErrNil && deleteErr == redis.ErrNil:
			m[members[i]] = s.Presence{
				Present: false,
			}
		default:
			err := typex.Errorf(errors.Source, errors.UnexpectedResults,
				"Score Member Error for %v (%v/%v)",
				members[i],
				insertErr,
				deleteErr,
			)
			return map[s.KeyFieldTxnValue]s.Presence{}, err
		}
	}

	return m, nil
}
