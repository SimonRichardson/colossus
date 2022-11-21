package plaintext

import (
	"flag"
	"os"
	"runtime"
	"strings"
	synx "sync"
	"testing"
	"testing/quick"
	"time"

	"github.com/SimonRichardson/colossus/teleprinter/common"
	"github.com/SimonRichardson/colossus/teleprinter/logs"
)

var (
	defaultUseStubs = false
)

func TestMain(t *testing.M) {
	var flagStubs bool
	flag.BoolVar(&flagStubs, "stubs", false, "enable stubs testing")
	flag.Parse()

	defaultUseStubs = flagStubs

	os.Exit(t.Run())
}

func config() *quick.Config {
	if testing.Short() {
		return &quick.Config{
			MaxCount:      10,
			MaxCountScale: 10,
		}
	}
	return nil
}

type mockWriter struct {
	mutex  *synx.Mutex
	amount int
}

func newWriter() *mockWriter {
	return &mockWriter{&synx.Mutex{}, 0}
}

func (w *mockWriter) Write(p []byte) (int, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.amount += len(strings.Split(string(p), "\n")) - 1
	return 0, nil
}

// RingBuffer

func testRingBuffer(amount, total int) int {
	var (
		writer = newWriter()
		ch     = make(chan time.Time)
		buffer = newRingBuffer(writer, func() <-chan time.Time {
			return ch
		}, logs.Info, total)
	)

	for i := 0; i < amount; i++ {
		buffer.Printf("%d\n", i)
	}

	ch <- time.Now()

	runtime.Gosched()
	time.Sleep(time.Millisecond * 20)

	return writer.amount
}

func TestRingBuffer_SameAmount(t *testing.T) {
	f := func(amount int) bool {
		num := common.Abs(amount) % 100
		return testRingBuffer(num, num) == num
	}

	if err := quick.Check(f, config()); err != nil {
		t.Error(err)
	}
}

func TestRingBuffer_LowerAmount(t *testing.T) {
	f := func(amount int) bool {
		num := common.Abs(amount) % 100
		res := testRingBuffer(num, num*2)
		return res == num
	}

	if err := quick.Check(f, config()); err != nil {
		t.Error(err)
	}
}

func TestRingBuffer_OverAmount(t *testing.T) {
	f := func(amount int) bool {
		num := common.MaxInt(common.Abs(amount)%100, 4)
		return testRingBuffer(num, num/2) == num/2
	}

	if err := quick.Check(f, config()); err != nil {
		t.Error(err)
	}
}
