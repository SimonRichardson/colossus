package store

import (
	"time"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/errors"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/typex"
	"github.com/garyburd/redigo/redis"
)

// ErrExpiredNode defines an error where the node has expired and can't be viewed
var (
	ErrExpiredNode = typex.Errorf(errors.Source, errors.ExpiredNode, "Expired node")
)

func selection(conn redis.Conn, key, field bs.Key) (s.KeyFieldScoreTxnValue, error) {
	var (
		now        = time.Now().UnixNano()
		value, err = redis.String(conn.Do("HGET", prefix+key.String()+insertSuffix, field))
	)
	if err != nil {
		return s.KeyFieldScoreTxnValue{}, err
	}

	return extractKeyFieldScoreTxnValue(now, key, field.String(), value)
}

func selectionWithRange(conn redis.Conn, key bs.Key, limit int, sizeExpiry s.SizeExpiry) ([]s.KeyFieldScoreTxnValue, error) {
	var (
		cursor = 0
		result = []s.KeyFieldScoreTxnValue{}
		now    = time.Now().UnixNano()
	)

	for {
		values, err := redis.Values(conn.Do("HSCAN", prefix+key.String()+insertSuffix, cursor, "COUNT", limit))
		if err != nil {
			return nil, err
		}

		if n := len(values); n != 2 {
			return nil, typex.Errorf(errors.Source, errors.UnexpectedResults,
				"Received %d values from redis, expected 2",
				n,
			)
		}

		newCursor, err := redis.Int(values[0], nil)
		if err != nil {
			return nil, err
		}

		pairs, err := redis.Values(values[1], nil)
		if err != nil {
			return nil, err
		}

		if num := len(pairs); num > 0 && num%2 != 0 {
			return nil, typex.Errorf(errors.Source, errors.UnexpectedResults,
				"Received %d pairs from redis, expected multipules of 2",
				num,
			)
		}

		members, err := extractMembers(now, key, pairs)
		if err != nil {
			return nil, err
		}

		result = append(result, members...)
		if len(result) >= limit {
			return result[:limit], nil
		}

		if newCursor <= 0 {
			break
		}

		cursor = newCursor
	}

	return result, nil
}

func extractMembers(now int64, key bs.Key, pairs []interface{}) ([]s.KeyFieldScoreTxnValue, error) {
	var (
		num = len(pairs)
		res = make([]s.KeyFieldScoreTxnValue, 0, num/2)
	)

	for i := 0; i < num; i += 2 {
		field, err := redis.String(pairs[i], nil)
		if err != nil {
			return nil, err
		}

		value, err := redis.String(pairs[i+1], nil)
		if err != nil {
			return nil, err
		}

		member, err := extractKeyFieldScoreTxnValue(now, key, field, value)
		if err != nil {
			if err == ErrExpiredNode {
				continue
			}
			return nil, err
		}
		res = append(res, member)
	}
	return res, nil
}

func extractKeyFieldScoreTxnValue(now int64, key bs.Key, field, value string) (s.KeyFieldScoreTxnValue, error) {
	score, txn, expiry, value, err := ExtractScoreTxnExpiryValue(value)
	if err != nil {
		return s.KeyFieldScoreTxnValue{}, err
	}

	var expireErr error
	if expiry < now {
		expireErr = ErrExpiredNode
	}

	return s.KeyFieldScoreTxnValue{
		Key:   key,
		Field: bs.Key(field),
		Score: score,
		Txn:   bs.Key(txn),
		Value: value,
	}, expireErr
}
