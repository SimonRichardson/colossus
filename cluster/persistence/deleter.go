package persistence

import (
	bs "github.com/SimonRichardson/colossus/blist/selectors"
	t "github.com/SimonRichardson/colossus/cluster"
	"github.com/SimonRichardson/colossus/fusion/mongo"
	s "github.com/SimonRichardson/colossus/selectors"
	"gopkg.in/mgo.v2/bson"
)

func deletion(db mongo.Database, fn s.Transformer, members []s.KeyFieldScoreTxnValue) ([]s.KeyCount, error) {
	m := make(map[bs.Key][]s.KeyFieldScoreTxnValue, 0)
	for _, member := range members {
		m[member.Key] = append(m[member.Key], member)
	}

	result := make([]s.KeyCount, 0, len(members))

	for k, v := range m {
		collection := db.C(collectionName(k))
		if _, err := collection.RemoveAll(bson.M{
			"txn": bson.M{"$in": memberTxnToObjectIds(v)},
		}); err != nil {
			continue
		} else {
			for range v {
				result = append(result, s.KeyCount{
					Key:   k,
					Count: 1,
				})
			}
		}
	}

	if len(result) < len(members) {
		return result, t.ErrPartialDeletions
	}

	return result, nil
}

func memberTxnToObjectIds(m []s.KeyFieldScoreTxnValue) []bson.ObjectId {
	res := make([]bson.ObjectId, 0, len(m))
	for _, v := range m {
		res = append(res, bson.ObjectIdHex(v.Txn.String()))
	}
	return res
}
