package noop

import "github.com/SimonRichardson/colossus/teleprinter/logs"

type log struct{}

func New() logs.Log {
	return log{}
}

func (log) Info() logs.Logger  { return logger{} }
func (log) Instr() logs.Logger { return logger{} }
func (log) Error() logs.Logger { return logger{} }

type logger struct{}

func (logger) Printf(string, ...interface{}) {}
func (logger) Println(...interface{})        {}
func (logger) Write([]byte) (int, error)     { return 0, nil }
func (logger) HR()                           {}
func (logger) Segment() logs.Segment {
	return segment{}
}

type segment struct {
	logger
}

func (s segment) Flush() {}
