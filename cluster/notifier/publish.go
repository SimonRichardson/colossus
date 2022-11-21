package notifier

import (
	"github.com/SimonRichardson/colossus/schemas/pool"
	"github.com/SimonRichardson/colossus/schemas/records"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/garyburd/redigo/redis"
)

func publish(conn redis.Conn, channel s.Channel, values []s.KeyFieldScoreSizeExpiry) error {
	fb := pool.Get()
	defer pool.Put(fb)

	kfs := records.KeyFieldScoreSizeExpiry{}

	for _, v := range values {
		fb.Reset()

		kfs.Key = v.Key
		kfs.Field = v.Field
		kfs.Score = v.Score
		kfs.Size = v.Size
		kfs.Expiry = v.Expiry

		bytes, err := kfs.Write(fb)
		if err != nil {
			return err
		}

		if err := conn.Send("RPUSH", channel.String(), bytes); err != nil {
			return err
		}
	}

	if err := conn.Flush(); err != nil {
		return err
	}

	if !defaultVerifyResults {
		return nil
	}

	for range values {
		if _, err := conn.Receive(); err != nil {
			return err
		}
	}

	return nil
}
