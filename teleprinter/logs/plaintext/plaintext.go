package plaintext

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sort"
	synx "sync"
	"time"

	"strings"

	"github.com/SimonRichardson/colossus/teleprinter/logs"
)

const (
	defaultBufferAmount        = 100
	defaultBufferFlushDuration = time.Second * 1

	defaultHRRune   = '-'
	defaultHRAmount = 100
)

type log struct {
	info, instr, errors logs.Logger
}

func NewAsync(writer io.Writer) logs.Log {
	var (
		duration = defaultBufferFlushDuration
		ticker   = func() <-chan time.Time { return time.Tick(duration) }
	)
	return log{
		info:   newRingBuffer(writer, ticker, logs.Info, defaultBufferAmount),
		instr:  newRingBuffer(writer, ticker, logs.Instr, defaultBufferAmount),
		errors: makeSync(writer, logs.Error),
	}
}

func NewSync(writer io.Writer) logs.Log {
	return log{
		info:   makeSync(writer, logs.Info),
		instr:  makeSync(writer, logs.Instr),
		errors: makeSync(writer, logs.Error),
	}
}

func NewEmojiSync(writer io.Writer) logs.Log {
	return log{
		info:   makeEmojiSync(writer, logs.Info),
		instr:  makeSync(writer, logs.Instr),
		errors: makeEmojiSync(writer, logs.Error),
	}
}

func (l log) Info() logs.Logger  { return l.info }
func (l log) Instr() logs.Logger { return l.instr }
func (l log) Error() logs.Logger { return l.errors }

type sync struct {
	writer io.Writer
	level  logs.LogLevel
}

func makeSync(writer io.Writer, level logs.LogLevel) sync {
	return sync{writer, level}
}

func (l sync) Printf(format string, args ...interface{}) {
	f := fmt.Sprintf("%s [%s] ", formatTime(time.Now()), l.level.String()) + format
	fmt.Fprintf(l.writer, f, args...)
}

func (l sync) Println(args ...interface{}) {
	f := fmt.Sprintf("%s [%s]", formatTime(time.Now()), l.level.String())
	fmt.Fprintln(l.writer, append([]interface{}{f}, args...)...)
}

func (l sync) Write(p []byte) (int, error) {
	f := fmt.Sprintf("%s [%s] %s", formatTime(time.Now()), l.level.String(), p)
	return l.writer.Write([]byte(f))
}

func (l sync) HR() {
	l.Println(strings.Repeat(string(defaultHRRune), defaultHRAmount))
}

func (l sync) Segment() logs.Segment {
	var (
		buffer = new(bytes.Buffer)
		writer = bufio.NewWriter(buffer)
	)
	return closer{
		makeSync(writer, l.level),
		l.writer,
		buffer,
		writer,
	}
}

type emojiSync struct {
	sync sync
}

func makeEmojiSync(writer io.Writer, level logs.LogLevel) emojiSync {
	return emojiSync{makeSync(writer, level)}
}

func (l emojiSync) Printf(format string, args ...interface{}) {
	inject(&args)
	l.sync.Printf(compile(format), args...)
}

func (l emojiSync) Println(args ...interface{}) {
	inject(&args)
	l.sync.Println(args...)
}

func (l emojiSync) Write(p []byte) (int, error) {
	return l.sync.Write(p)
}

func (l emojiSync) HR() {
	l.sync.HR()
}

func (l emojiSync) Segment() logs.Segment {
	return l.sync.Segment()
}

func inject(a *[]interface{}) {
	values := *a
	for k, v := range values {
		if s, ok := v.(string); ok {
			values[k] = compile(s)
		}
	}
}

func compile(str string) string {
	if len(str) < 1 {
		return ""
	}

	var (
		input  = bytes.NewBufferString(str)
		output = bytes.NewBufferString("")
	)

	for {
		r, _, err := input.ReadRune()
		if err != nil {
			break
		}

		if r == ':' {
			output.WriteString(peek(input))
			continue
		}

		output.WriteRune(r)
	}

	return output.String()
}

func peek(input *bytes.Buffer) string {
	res := bytes.NewBufferString(":")
	for {
		r, _, err := input.ReadRune()
		if err != nil {
			break
		}

		res.WriteRune(r)

		if r <= ' ' {
			break
		}

		if r == ':' {
			if emoji, ok := codeMap[res.String()]; ok {
				res.Reset()
				res.WriteString(fmt.Sprintf("%s ", emoji))
				break
			}
		}
	}

	return res.String()
}

type closer struct {
	sync
	writer         io.Writer
	buffer         *bytes.Buffer
	bufferedWriter *bufio.Writer
}

func (l closer) Flush() {
	l.bufferedWriter.Flush()
	l.writer.Write(l.buffer.Bytes())
}

type ringBuffer struct {
	writer io.Writer
	ticker func() <-chan time.Time
	level  logs.LogLevel
	ring   *Ring
}

func newRingBuffer(writer io.Writer,
	ticker func() <-chan time.Time,
	level logs.LogLevel,
	amount int,
) *ringBuffer {
	r := &ringBuffer{writer, ticker, level, newRing(amount)}
	go r.run()
	return r
}

func (l *ringBuffer) run() {
	for range l.ticker() {
		go flush(l.writer, l.ring.Flush())
	}
}

func (l *ringBuffer) Printf(format string, args ...interface{}) {
	l.ring.Update(printf, l.level, time.Now(), format, args)
}

func (l *ringBuffer) Println(args ...interface{}) {
	l.ring.Update(println, l.level, time.Now(), "", args)
}

func (l *ringBuffer) Write(p []byte) (int, error) {
	l.ring.Update(printf, l.level, time.Now(), string(p), []interface{}{})
	return 0, nil
}

func (l *ringBuffer) HR() {
	l.Println(strings.Repeat(string(defaultHRRune), defaultHRAmount))
}

func (l *ringBuffer) Segment() logs.Segment {
	var (
		buffer = new(bytes.Buffer)
		writer = bufio.NewWriter(buffer)
	)
	return closer{
		makeSync(writer, l.level),
		l.writer,
		buffer,
		writer,
	}
}

func flush(writer io.Writer, list []*RingNode) {
	if len(list) > 0 {
		// Sort the list according to the time.
		sort.Sort(ByTime(list))

		// Dump out the list.
		var buffer bytes.Buffer
		for _, v := range list {
			buffer.WriteString(v.String())
		}

		fmt.Fprintf(writer, buffer.String())
	}
}

type printType int

const (
	printf printType = iota
	println
)

type Ring struct {
	mutex   *synx.Mutex
	list    []*RingNode
	current int
}

func newRing(amount int) *Ring {
	ring := &Ring{&synx.Mutex{}, make([]*RingNode, amount, amount), 0}
	for i := 0; i < amount; i++ {
		ring.list[i] = &RingNode{filled: false}
	}
	return ring
}

func (r *Ring) Update(p printType, l logs.LogLevel, t time.Time, f string, a []interface{}) {
	r.mutex.Lock()

	node := r.list[r.current%cap(r.list)]
	node.filled = true
	node.printType = p
	node.level = l
	node.time = t
	node.format = f
	node.args = a
	r.current++

	r.mutex.Unlock()
}

func (r *Ring) Flush() []*RingNode {
	list := []*RingNode{}
	for _, v := range r.list {
		if v.filled {
			list = append(list, v)
		}
		v.filled = false
	}
	return list
}

type RingNode struct {
	filled    bool
	printType printType
	level     logs.LogLevel
	time      time.Time
	format    string
	args      []interface{}
}

func (n *RingNode) String() string {
	switch n.printType {
	case printf:
		f := fmt.Sprintf("%s [%s] ", formatTime(n.time), n.level.String()) + n.format
		return fmt.Sprintf(f, n.args...)
	case println:
		f := fmt.Sprintf("%s [%s]", formatTime(n.time), n.level.String())
		return fmt.Sprintln(append([]interface{}{f}, n.args...)...)
	}
	return ""
}

type ByTime []*RingNode

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].time.UnixNano() < a[j].time.UnixNano() }

func formatTime(t time.Time) string {
	var buffer bytes.Buffer
	year, month, day := t.Date()
	buffer.WriteString(fmt.Sprintf("%d/%d/%d ", year, month, day))
	hour, min, sec := t.Clock()
	buffer.WriteString(fmt.Sprintf("%02d:%02d:%02d", hour, min, sec))
	ns := t.Nanosecond() / 1e3
	buffer.WriteString(fmt.Sprintf(".%06d", ns))
	return buffer.String()
}
