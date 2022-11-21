package counter

import (
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
		if err := conn.Send("ZSCORE", prefix+key+insertSuffix, field); err != nil {
			return map[s.KeyFieldTxnValue]s.Presence{}, err
		}

		if err := conn.Send("ZSCORE", prefix+key+deleteSuffix, field); err != nil {
			return map[s.KeyFieldTxnValue]s.Presence{}, err
		}
	}

	if err := conn.Flush(); err != nil {
		return map[s.KeyFieldTxnValue]s.Presence{}, err
	}

	m := map[s.KeyFieldTxnValue]s.Presence{}
	for i := 0; i < len(members); i++ {
		var (
			insertScore, insertErr = redis.Float64(conn.Receive())
			deleteScore, deleteErr = redis.Float64(conn.Receive())
		)

		switch {
		case insertErr == nil && deleteErr == redis.ErrNil:
			m[members[i]] = s.Presence{
				Present:  true,
				Inserted: true,
				Score:    insertScore,
			}
		case insertErr == redis.ErrNil && deleteErr == nil:
			m[members[i]] = s.Presence{
				Present:  true,
				Inserted: false,
				Score:    deleteScore,
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
