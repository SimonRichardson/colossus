package teleprinter

import (
	"os"

	"github.com/SimonRichardson/colossus/teleprinter/logs"
	"github.com/SimonRichardson/colossus/teleprinter/logs/plaintext"
)

func DefaultLog() logs.Log {
	L = plaintext.NewSync(os.Stdout)
	return L
}

var (
	L logs.Log
)
