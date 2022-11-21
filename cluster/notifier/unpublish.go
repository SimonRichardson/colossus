package notifier

import (
	"fmt"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	s "github.com/SimonRichardson/colossus/selectors"
	"github.com/garyburd/redigo/redis"
)

const (
	prefix       = "u:"
	deleteSuffix = "-"
)

func unpublish(conn redis.Conn, channel s.Channel, values []s.KeyFieldScoreSizeExpiry) error {
	for _, v := range values {

		key := unpublishKey(channel, v.Key, v.Field)
		if err := conn.Send("SETEX", key, fmt.Sprintf("%.0f", v.Expiry.Seconds()), true); err != nil {
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

func unpublishKey(channel s.Channel, key, field bs.Key) string {
	return fmt.Sprintf("%s%s%s%s%s", prefix, channel.String(),
		key.String(), field.String(), deleteSuffix,
	)
}
