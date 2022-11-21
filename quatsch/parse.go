package quatsch

import (
	"strings"

	"fmt"

	"time"

	"github.com/SimonRichardson/colossus/quatsch/pool/bson"
	"github.com/SimonRichardson/colossus/quatsch/pool/uuid"
)

type PoolOptions struct {
	Duration    time.Duration
	PerDuration int64
}

func Parse(variance string, max int, opts PoolOptions) (Variance, error) {
	switch strings.ToLower(variance) {
	case "bson":
		return bson.Parse(max, opts.Duration, opts.PerDuration)
	case "uuid":
		return uuid.Parse(max, opts.Duration, opts.PerDuration)
	}
	return nil, fmt.Errorf("No variance found for %q", variance)
}
