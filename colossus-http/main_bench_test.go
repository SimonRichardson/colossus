package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"gopkg.in/mgo.v2/bson"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/env"
	"github.com/SimonRichardson/colossus/schemas/records"
	"github.com/SimonRichardson/colossus/tests"
	"github.com/SimonRichardson/colossus/typex"
	flatbuffers "github.com/google/flatbuffers/go"
)

var result int

func benchmarkRequest(a int, b *testing.B) {
	// Turn off logging for benchmarking
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

	e := env.New(nil)
	e.Logs = "Noop"
	e.Instrumentation = "Noop"

	var (
		ts, _ = setup(e)

		key = bson.NewObjectId()

		inserter = benchPost(ts.URL, key, 999999)

		bodies = make([]tests.PostBody, b.N)
	)

	for i := 0; i < b.N; i++ {
		bodies[i] = generatePostBody(a)
	}

	defer tear(ts)

	b.ResetTimer()

	// Do this so we don't get optimized out!
	var response int
	for i := 0; i < b.N; i++ {
		response = inserter(bodies[i])
	}

	result = response
}

func BenchmarkRequest1(b *testing.B) { benchmarkRequest(1, b) }
func BenchmarkRequest2(b *testing.B) { benchmarkRequest(2, b) }
func BenchmarkRequest3(b *testing.B) { benchmarkRequest(3, b) }
func BenchmarkRequest4(b *testing.B) { benchmarkRequest(4, b) }
func BenchmarkRequest5(b *testing.B) { benchmarkRequest(5, b) }
func BenchmarkRequest6(b *testing.B) { benchmarkRequest(6, b) }

func benchmarkRequestAndAssign(a int, b *testing.B) {
	// Turn off logging for benchmarking
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

	e := env.New(nil)
	e.Logs = "Noop"
	e.Instrumentation = "Noop"

	var (
		maxSize = 999999

		ts, _ = setup(e)

		key = bson.NewObjectId()

		inserter = benchPost(ts.URL, key, maxSize)

		postBodies = make([]tests.PostBody, b.N)
		putBodies  = make([]tests.PutBody, b.N)

		modify = func(url string, key bson.ObjectId, values tests.PutBody) int {
			// Modify
			record := records.PutRecords{
				Key:     bs.Key(key.Hex()),
				Records: values,
				Score:   2,
				MaxSize: int64(maxSize),
				Expiry:  defaultExpiry,
			}
			bytes, err := record.Write(flatbuffers.NewBuilder(0))
			if err != nil {
				typex.Fatal(err)
			}

			body := tests.Put(fmt.Sprintf("%s/http/v1/%s", url, key.Hex()), bytes)

			s := &records.OKInt{}
			s.Read(body)

			return s.Records
		}
	)

	for i := 0; i < b.N; i++ {
		post := generatePostBody(a)

		postBodies[i] = post
		putBodies[i] = post.PutBody()
	}

	defer tear(ts)

	b.ResetTimer()

	// Do this so we don't get optimized out!
	var response int
	for i := 0; i < b.N; i++ {
		response = inserter(postBodies[i])
		response = modify(ts.URL, key, putBodies[i])
	}

	result = response
}

func BenchmarkRequestAndAssign1(b *testing.B) { benchmarkRequestAndAssign(1, b) }
func BenchmarkRequestAndAssign2(b *testing.B) { benchmarkRequestAndAssign(2, b) }
func BenchmarkRequestAndAssign3(b *testing.B) { benchmarkRequestAndAssign(3, b) }
func BenchmarkRequestAndAssign4(b *testing.B) { benchmarkRequestAndAssign(4, b) }
func BenchmarkRequestAndAssign5(b *testing.B) { benchmarkRequestAndAssign(5, b) }
func BenchmarkRequestAndAssign6(b *testing.B) { benchmarkRequestAndAssign(6, b) }

func benchPost(url string, key bson.ObjectId, maxSize int) func(tests.PostBody) int {
	return func(values tests.PostBody) int {
		record := records.PostRecords{
			Records: values,
			Score:   1,
			MaxSize: int64(maxSize),
			Expiry:  defaultExpiry,
		}
		bytes, err := record.Write(flatbuffers.NewBuilder(0))
		if err != nil {
			typex.Fatal(err)
		}

		body := tests.Post(fmt.Sprintf("%s/http/v1/%s", url, key.Hex()), bytes)

		s := &records.OKInt{}
		s.Read(body)

		return s.Records
	}
}

func benchDelete(url string, key bson.ObjectId, opts *deleteOpts) func(tests.PostBody) int {
	return func(values tests.PostBody) int {
		// Insert
		inserter := benchPost(url, key, defaultMaxSize)
		inserter(values)

		// Delete
		slice := tests.PostBody(values[:deleteAmount(opts, values)])
		record := records.DeleteRecords{
			Key:     bs.Key(key.Hex()),
			Records: slice.DeleteBody(),
			Score:   2,
			MaxSize: int64(defaultMaxSize),
			Expiry:  defaultExpiry,
		}
		bytes, err := record.Write(flatbuffers.NewBuilder(0))
		if err != nil {
			typex.Fatal(err)
		}

		body := tests.Del(fmt.Sprintf("%s/http/v1/%s", url, key.Hex()), bytes)

		s := &records.OKInt{}
		s.Read(body)

		return s.Records
	}
}

func benchRollback(url string, key bson.ObjectId, opts *rollbackOpts) func(tests.PostBody) int {
	return func(values tests.PostBody) int {
		// Insert
		inserter := benchPost(url, key, defaultMaxSize)
		inserter(values)

		// Rollback
		slice := tests.PostBody(values[:rollbackAmount(opts, values)])
		record := records.RollbackRecords{
			Key:     bs.Key(key.Hex()),
			Records: slice.RollbackBody(),
			Score:   2,
			MaxSize: int64(defaultMaxSize),
			Expiry:  defaultExpiry,
		}
		bytes, err := record.Write(flatbuffers.NewBuilder(0))
		if err != nil {
			typex.Fatal(err)
		}

		tests.Del(fmt.Sprintf("%s/http/v1/%s/rollback", url, key.Hex()), bytes)

		return rollbackAmount(opts, values)
	}
}
