package parse

import (
	"os"
	"strings"

	"github.com/SimonRichardson/colossus/teleprinter/common"
	"github.com/SimonRichardson/colossus/teleprinter/logs"
	"github.com/SimonRichardson/colossus/teleprinter/logs/noop"
	"github.com/SimonRichardson/colossus/teleprinter/logs/plaintext"
	"github.com/SimonRichardson/colossus/typex"
)

const (
	Teleprinter typex.ErrorSource = "Teleprinter"
)

var (
	NoCaseFound = typex.InternalServerError.With("No Case Found")
)

func ParseString(value string) (logs.Log, error) {
	parts := strings.Split(value, ";")
	switch common.Normalise(parts[0]) {
	case "noop":
		return noop.New(), nil
	case "plaintext":
		return plaintext.NewSync(os.Stdout), nil
	case "plaintext-buffered":
		return plaintext.NewAsync(os.Stdout), nil
	case "emoji":
		return plaintext.NewEmojiSync(os.Stdout), nil
	}
	return noop.New(), typex.Errorf(Teleprinter, NoCaseFound,
		"Invalid logs %q", value)
}
