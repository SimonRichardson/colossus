package redis

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/SimonRichardson/colossus/fusion/errors"
	"github.com/SimonRichardson/colossus/typex"
)

func ValidRedisHost(host string) error {
	if strings.Contains(host, ":") {
		url, err := url.Parse(host)
		if err != nil {
			return typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
				"Invalid host %q (%s)", host, err)
		}

		tokens := strings.Split(url.Host, ":")
		if len(tokens) < 2 {
			return typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
				"Invalid host %q", host)
		}
		if _, err := strconv.ParseUint(tokens[1], 10, 16); err != nil {
			return typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
				"Invalid port %q in host %q (%s)", tokens[1], host, err)
		}

		if url.User != nil {
			if password, ok := url.User.Password(); ok && len(password) < 1 {
				return typex.Errorf(errors.Source, errors.UnexpectedParseArgument,
					"Invalid password %q in host %q", password, host)
			}
		}
	}
	return nil
}
