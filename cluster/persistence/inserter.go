package persistence

import (
	"fmt"

	"gopkg.in/mgo.v2"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	t "github.com/SimonRichardson/colossus/cluster"
	"github.com/SimonRichardson/colossus/fusion/mongo"
	s "github.com/SimonRichardson/colossus/selectors"
)

func insertion(db mongo.Database, fn s.Transformer, members []s.KeyFieldScoreTxnValue) ([]s.KeyCount, error) {
	m := make(map[bs.Key][]s.KeyFieldScoreTxnValue, 0)
	for _, member := range members {
		m[member.Key] = append(m[member.Key], member)
	}

	result := make([]s.KeyCount, 0, len(members))

	for k, v := range m {
		var (
			collection = db.C(collectionName(k))
			bulk       = collection.Bulk()
		)

		bulk.Unordered()

		for _, element := range v {
			res, err := fn(element)
			if err != nil {
				return result, err
			}
			bulk.Insert(res)
		}

		_, err := bulk.Run()
		if err != nil {
			// We don't care it's a duplicate. The way bulk works is to continue
			// on executing the insertions, so it could be possible that a
			// duplicate error could happen.
			if mgo.IsDup(err) {
				continue
			}
			return result, err
		}

		for range v {
			result = append(result, s.KeyCount{Key: k, Count: 1})
		}
	}

	if len(result) < len(members) {
		return result, t.ErrPartialInsertions
	}

	return result, nil
}

func collectionName(key bs.Key) string {
	return fmt.Sprintf("tickets_%s", key.String())
}
