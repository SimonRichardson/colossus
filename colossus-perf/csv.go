package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/SimonRichardson/colossus/typex"
)

type row struct {
	now                     time.Time
	total, count, result    int
	duration                time.Duration
	ciServer, clusterServer bool
}

type CSV struct {
	rows []row
}

func newCSV() *CSV {
	return &CSV{rows: []row{}}
}

func (v *CSV) add(t, c, r int, d time.Duration, ci, cl bool) {
	v.rows = append(v.rows, row{
		now:   time.Now(),
		total: t, count: c, result: r,
		duration: d,
		ciServer: ci, clusterServer: cl,
	})
}

func (v *CSV) write() []byte {
	records := [][]string{
		{"n", "time", "total", "count", "result", "duration", "ci-server", "cluster"},
	}

	for k, v := range v.rows {
		records = append(records, []string{
			fmt.Sprintf("%d", k),
			fmt.Sprintf("%d", v.now.UnixNano()/(int64(time.Millisecond)/int64(time.Nanosecond))),
			fmt.Sprintf("%d", v.total),
			fmt.Sprintf("%d", v.count),
			fmt.Sprintf("%d", v.result),
			fmt.Sprintf("%d", v.duration.Nanoseconds()),
			fmt.Sprintf("%v", v.ciServer),
			fmt.Sprintf("%v", v.clusterServer),
		})
	}

	var (
		b  bytes.Buffer
		bw = bufio.NewWriter(&b)
		w  = csv.NewWriter(bw)
	)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			typex.Fatal("Error writing record to csv:", err)
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		typex.Fatal(err)
	}

	return b.Bytes()
}
