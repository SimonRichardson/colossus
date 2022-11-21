package persistence

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/common"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/fusion/mongo"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

func repair(db mongo.Database, fn s.Transformer, members []s.KeyFieldScoreTxnValue) ([]s.KeyCount, error) {
	m := make(map[bs.Key][]s.KeyFieldScoreTxnValue, 0)
	for _, member := range members {
		m[member.Key] = append(m[member.Key], member)
	}

	var (
		result = make([]s.KeyCount, 0)
		errs   = make([]error, 0)
	)

	for k, v := range m {
		collection := db.C(k.String())

		for _, element := range v {
			var (
				id       = bson.ObjectIdHex(element.Field.String())
				res, err = fn(element)
			)
			if err != nil {
				return result, err
			}

			changes, err := collection.UpsertId(id, res)
			if err != nil {
				// This can happen if another colossus happens to repair at the
				// same time, which causes this mongo defect.
				// https://jira.mongodb.org/browse/SERVER-14322
				if mgo.IsDup(err) {
					continue
				}
				errs = append(errs, err)
				continue
			}

			if changes.UpsertedId != nil || (changes.Matched == 0 && changes.Updated == 1) {
				// Repair happened!
				result = append(result, s.KeyCount{Key: k, Count: 1})
			}
		}
	}

	if len(errs) > 0 {
		err := common.SumErrors(errs)
		return result, typex.Errorf(errors.Source, errors.Partial, "Partial Error (%s)", err.Error())
	}

	return result, nil
}
