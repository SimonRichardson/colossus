package counter

import (
	"fmt"
	"strings"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/scripts"
	"github.com/SimonRichardson/colossus/typex"
	"github.com/garyburd/redigo/redis"
)

const (
	cardinalityPrefix    = "k:"
	cardinalityPrefixLen = len(cardinalityPrefix)

	prefix    = "c:"
	prefixLen = len(prefix)

	insertSuffix = "+"
	deleteSuffix = "-"

	insertSuffixLen = len(insertSuffix)
	deleteSuffixLen = len(deleteSuffix)
)

var (
	genericScript string
	insertScript  *redis.Script
	deleteScript  *redis.Script
)

func init() {
	raw, err := scripts.Asset("../scripts/counter/script.lua")
	if err != nil {
		typex.Fatal(err)
	}
	script := string(raw)

	genericScript = strings.NewReplacer(
		"PREFIX", cardinalityPrefix,
		"INSERTSUFFIX", insertSuffix,
		"DELETESUFFIX", deleteSuffix,
	).Replace(script)

	insertScript = redis.NewScript(1, strings.NewReplacer(
		"REMSUFFIX", deleteSuffix,
		"ADDSUFFIX", insertSuffix,
	).Replace(genericScript))

	deleteScript = redis.NewScript(1, strings.NewReplacer(
		"REMSUFFIX", insertSuffix,
		"ADDSUFFIX", deleteSuffix,
	).Replace(genericScript))
}

func doInsertScript(conn redis.Conn,
	key, field bs.Key,
	score float64,
	maxSize int64,
) (interface{}, error) {
	return insertScript.Do(conn,
		prefix+key.String(),
		field.String(),
		score,
		fmt.Sprintf("%d", maxSize),
	)
}

func sendInsertScript(conn redis.Conn,
	key, field bs.Key,
	score float64,
	maxSize int64,
) error {
	return insertScript.Send(conn,
		prefix+key.String(),
		field.String(),
		score,
		fmt.Sprintf("%d", maxSize),
	)
}

func doDeleteScript(conn redis.Conn,
	key, field bs.Key,
	score float64,
	maxSize int64,
) (interface{}, error) {
	return deleteScript.Do(conn,
		prefix+key.String(),
		field.String(),
		score,
		fmt.Sprintf("%d", maxSize),
	)
}

func sendDeleteScript(conn redis.Conn,
	key, field bs.Key,
	score float64,
	maxSize int64,
) error {
	return deleteScript.Send(conn,
		prefix+key.String(),
		field.String(),
		score,
		fmt.Sprintf("%d", maxSize),
	)
}
