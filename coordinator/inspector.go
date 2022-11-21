package coordinator

import (
	"gopkg.in/mgo.v2/bson"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/farm/store"
	"github.com/SimonRichardson/colossus/schemas/records"
	"github.com/SimonRichardson/colossus/schemas/schema"
	s "github.com/SimonRichardson/colossus/selectors"
)

type inspector struct {
	s.LifeCycleManager

	co          *Coordinator
	store       *store.Farm
	transformer s.Transformer
}

func newInspector(co *Coordinator, store *store.Farm, transformer s.Transformer) *inspector {
	return &inspector{
		LifeCycleManager: newLifeCycleService(),

		co:          co,
		store:       store,
		transformer: transformer,
	}
}

func (i *inspector) Query(key bs.Key,
	options s.QueryOptions,
	sizeExpiry s.SizeExpiry,
) ([]s.QueryRecord, error) {
	size, err := i.store.Size(key)
	if err != nil {
		return nil, err
	}

	members, err := i.store.SelectRange(key, size, s.KeySizeExpiry{
		key: sizeExpiry,
	})
	if err != nil {
		return nil, err
	}

	var (
		ownerId  = options.OwnerId.String()
		filtered = make([]s.QueryRecord, 0)
	)
	for _, v := range members {
		if header, err := records.ReadType(v.Value); err != nil || header != schema.TypePost {
			continue
		}

		if record, err := i.transformer(v); err == nil {
			if id, ok := record["owner_id"].(bson.ObjectId); ok {
				if ownerId == id.Hex() {
					filtered = append(filtered, s.QueryRecord{
						Key:    v.Key,
						Field:  v.Field,
						Record: record,
					})
				}
			}
		}
	}

	return filtered, nil
}
