package notifier

import (
	bs "github.com/SimonRichardson/colossus/blist/selectors"
	t "github.com/SimonRichardson/colossus/cluster"
	"github.com/SimonRichardson/colossus/schemas/records"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/garyburd/redigo/redis"
)

func subscribe(conn redis.Conn, channel s.Channel) <-chan t.Element {
	var (
		out = make(chan t.Element)
		key = bs.Key(defaultSubscribeKey)
	)

loop:
	for conn.Err() == nil {
		if err := conn.Send("LPOP", channel.String()); err != nil {
			out <- t.NewErrorElement(key, err)
			break loop
		}

		if err := conn.Flush(); err != nil {
			out <- t.NewErrorElement(key, err)
			break loop
		}

		reply, err := redis.Bytes(conn.Receive())
		if err != nil {
			out <- t.NewErrorElement(key, err)
			break loop
		}

		kfs := &records.KeyFieldScoreSizeExpiry{}
		if err := kfs.Read(reply); err != nil {
			out <- t.NewErrorElement(key, err)
			break loop
		}

		// If a value has been unpublished but we're unsure if it's already been
		// sent try and capture it on the other side to prevent it from saving.
		key := unpublishKey(channel, kfs.Key, kfs.Field)
		if val, err := redis.Bool(conn.Do("GET", key)); err != nil && val {
			continue
		}

		out <- t.NewKeyFieldScoreSizeExpiryElement(s.KeyFieldScoreSizeExpiry{
			Key:    kfs.Key,
			Field:  kfs.Field,
			Score:  kfs.Score,
			Size:   kfs.Size,
			Expiry: kfs.Expiry,
		})
	}

	close(out)
	return out
}
