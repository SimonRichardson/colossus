package strategies

import (
	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/farm/counter"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/typex"
)

// InsertStrategy defines if we can insert a new item.
type InsertStrategy func(*counter.Farm,
	map[bs.Key][]s.KeyFieldScoreTxnValue,
	s.KeySizeExpiry,
) (map[bs.Key][]s.KeyFieldScoreTxnValue, error)

func insertNoop(c *counter.Farm,
	buckets map[bs.Key][]s.KeyFieldScoreTxnValue,
	sizeExpiry s.KeySizeExpiry,
) (map[bs.Key][]s.KeyFieldScoreTxnValue, error) {
	return buckets, nil
}

func insertCounter(c *counter.Farm,
	buckets map[bs.Key][]s.KeyFieldScoreTxnValue,
	sizeExpiry s.KeySizeExpiry,
) (map[bs.Key][]s.KeyFieldScoreTxnValue, error) {
	sized := map[bs.Key][]s.KeyFieldScoreTxnValue{}

	for k, v := range buckets {
		size, err := c.Size(k)
		if err != nil {
			return sized, err
		}

		if sizeExpiry, err := sizeExpiry.Get(k); err == nil && int64(size+len(v)) <= sizeExpiry.Size {
			sized[k] = append(sized[k], v...)
		} else {
			return sized, typex.Errorf(errors.Source, errors.MaxSize,
				"Reached max size")
		}
	}

	if len(sized) < 1 {
		return sized, typex.Errorf(errors.Source, errors.Complete,
			"Complete insertion failure")
	}

	return sized, nil
}
