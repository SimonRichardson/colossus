package logs

type LogLevel int

const (
	Info LogLevel = iota
	Instr
	Error
)

func (l LogLevel) String() string {
	switch l {
	case Info:
		return "INFO"
	case Instr:
		return "INSTR"
	case Error:
		return "ERROR"
	}
	return "INVALID"
}

type Log interface {
	Info() Logger
	Instr() Logger
	Error() Logger
}

type Logger interface {
	Printf(string, ...interface{})
	Println(...interface{})
	Write([]byte) (int, error)
	HR()
	Segment() Segment
}

type Segment interface {
	Logger
	Flush()
}
