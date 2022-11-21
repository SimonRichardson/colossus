package store

import (
	"fmt"
	"strconv"
	"strings"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/scripts"
	"github.com/SimonRichardson/colossus/typex"
	"github.com/garyburd/redigo/redis"
)

const (
	prefix    = "s:"
	prefixLen = len(prefix)

	separator    = ","
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
	raw, err := scripts.Asset("../scripts/store/script.lua")
	if err != nil {
		typex.Fatal(err)
	}
	script := string(raw)

	genericScript = strings.NewReplacer(
		"SEPARATOR", separator,
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
	expiry int64,
	txn bs.Key,
	value string,
) (interface{}, error) {
	return insertScript.Do(conn,
		prefix+key.String(),
		field.String(),
		score,
		txn.String(),
		PackageScoreTxnExpiryValue(score, txn, expiry, value),
	)
}

func sendInsertScript(conn redis.Conn,
	key, field bs.Key,
	score float64,
	expiry int64,
	txn bs.Key,
	value string,
) error {
	return insertScript.Send(conn,
		prefix+key.String(),
		field.String(),
		score,
		txn.String(),
		PackageScoreTxnExpiryValue(score, txn, expiry, value),
	)
}

func doDeleteScript(conn redis.Conn,
	key, field bs.Key,
	score float64,
	expiry int64,
	txn bs.Key,
	value string,
) (interface{}, error) {
	return deleteScript.Do(conn,
		prefix+key.String(),
		field.String(),
		score,
		txn.String(),
		PackageScoreTxnExpiryValue(score, txn, expiry, value),
	)
}

func sendDeleteScript(conn redis.Conn,
	key, field bs.Key,
	score float64,
	expiry int64,
	txn bs.Key,
	value string,
) error {
	return deleteScript.Send(conn,
		prefix+key.String(),
		field.String(),
		score,
		txn.String(),
		PackageScoreTxnExpiryValue(score, txn, expiry, value),
	)
}

func PackageScoreTxnExpiryValue(score float64, txn bs.Key, expiry int64, value string) string {
	return fmt.Sprintf("%f%s%s%s%d%s%s",
		score, separator,
		txn.String(), separator,
		expiry, separator,
		value,
	)
}

func ExtractScoreTxnExpiryValue(value string) (float64, string, int64, string, error) {
	parts := strings.SplitN(value, separator, 4)
	if num := len(parts); num != 4 {
		return 0, "", 0, "", typex.Errorf(errors.Source, errors.UnexpectedResults, "Received %d parts, expected 4", num)
	}

	score, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, "", 0, "", err
	}

	expiry, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return 0, "", 0, "", err
	}
	return score, parts[1], expiry, parts[3], nil
}
