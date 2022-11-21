package main

import (
	"flag"
	"fmt"
	"math"
	"net/http/httptest"
	"os"
	"reflect"
	"sort"
	"testing"
	"testing/quick"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"runtime"

	"strings"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/blist/services/consul/client"
	"github.com/SimonRichardson/colossus/cluster/store"
	"github.com/SimonRichardson/colossus/coordinator"
	"github.com/SimonRichardson/colossus/env"
	"github.com/SimonRichardson/colossus/quatsch"
	b "github.com/SimonRichardson/colossus/quatsch/pool/bson"
	"github.com/SimonRichardson/colossus/schemas/records"
	"github.com/SimonRichardson/colossus/tests"
	"github.com/SimonRichardson/colossus/typex"
	"github.com/garyburd/redigo/redis"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/hashicorp/consul/api"
)

var (
	defaultMaxSize  = 99999
	defaultExpiry   = time.Minute * 5
	defaultUseStubs = false

	defaultEnv     *env.Env
	defaultSession *mgo.Session
	defaultConsul  *api.Client
	defaultRedi    []redis.Conn
)

func TestMain(t *testing.M) {
	var flagStubs bool
	flag.BoolVar(&flagStubs, "stubs", false, "enable stubs testing")
	flag.Parse()

	defaultUseStubs = flagStubs
	defaultEnv = env.New(nil)
	if !defaultUseStubs {
		var err error
		if defaultSession == nil {
			if defaultSession, err = mgo.Dial(defaultEnv.MongoInstances); err != nil {
				typex.Fatal("Intergration Testing Mongo: ", defaultEnv.MongoInstances, " : ", err)
			}
		}
		if defaultConsul == nil {
			config := api.DefaultConfig()
			config.Address = defaultEnv.ConsulInstances

			if defaultConsul, err = api.NewClient(config); err != nil {
				typex.Fatal("Intergration Testing Consul: ", defaultEnv.ConsulInstances, " : ", err)
			}
		}
		if defaultRedi == nil {
			var (
				redi  []redis.Conn
				banks = strings.Split(defaultEnv.StoreInstances, ";")
			)

			for _, address := range strings.Split(banks[0], ",") {
				if conn, err := redis.Dial("tcp", strings.Replace(address, "tcp://", "", 1)); err != nil {
					typex.Fatal("Intergration Testing Store (redis): ", defaultEnv.StoreInstances, " : ", err)
				} else {
					redi = append(redi, conn)
				}
			}

			defaultRedi = redi
		}
	}

	os.Exit(t.Run())
}

func setup(e *env.Env) (*httptest.Server, *coordinator.Coordinator) {
	e.Logs = "Noop"
	e.Instrumentation = "Noop"

	server := newServer(e)
	return httptest.NewServer(server.Handler), server.co
}

func tear(ts *httptest.Server) {
	ts.Close()
}

func getIdentPool() quatsch.Pool {
	var (
		maxBuffer               = 99999
		maxInsertionPerDuration = int64(1000000)
	)

	return quatsch.New(b.New(maxBuffer, time.Second, maxInsertionPerDuration))
}

// Test Version

func testVersion(url string,
	co *coordinator.Coordinator,
) (func() string, func() string) {
	var (
		f = func() string {
			body := tests.Get(fmt.Sprintf("%s/http/version", url))

			s := &records.OKVersion{}
			s.Read(body)

			return s.Records.Version
		}
		g = func() string {
			return defaultEnv.Version
		}
	)
	return f, g
}

func TestVersion(t *testing.T) {
	e := env.New(nil)

	ts, co := setup(e)
	defer tear(ts)

	f, g := testVersion(ts.URL, co)

	if err := quick.CheckEqual(f, g, tests.Config()); err != nil {
		t.Error(err)
	}
}

// Test Post

func testPost(url string,
	maxSize int,
	co *coordinator.Coordinator,
) func(tests.PostBody) bool {
	pool := getIdentPool()
	return func(values tests.PostBody) bool {
		key, err := b.Bson(pool.Get())
		if err != nil {
			typex.Fatal(err)
		}

		return benchPost(url, key, maxSize)(values) == len(values)
	}
}

func TestPost_InsertAllReadAll(t *testing.T) {
	e := env.New(nil)
	e.StoreInsertStrategy = "InsertAllReadAll"

	ts, co := setup(e)
	defer tear(ts)

	f := testPost(ts.URL, defaultMaxSize, co)

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}

// Test Put

func testPut(url string,
	maxSize int,
	co *coordinator.Coordinator,
) func(tests.PostBody) bool {
	pool := getIdentPool()
	return func(values tests.PostBody) bool {
		key, err := b.Bson(pool.Get())
		if err != nil {
			typex.Fatal(err)
		}
		// Insert
		inserter := benchPost(url, key, maxSize)
		inserter(values)

		// Modify
		record := records.PutRecords{
			Key:     bs.Key(key.Hex()),
			Records: values.PutBody(),
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

		return s.Records == len(values)
	}

}

func TestPut_ModifyAllReadAll(t *testing.T) {
	e := env.New(nil)

	ts, co := setup(e)
	defer tear(ts)

	f := testPut(ts.URL, defaultMaxSize, co)

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}

func testPutPersistency(url string,
	maxSize int,
	co *coordinator.Coordinator,
) func(tests.PostBody) bool {
	pool := getIdentPool()
	return func(values tests.PostBody) bool {
		key, err := b.Bson(pool.Get())
		if err != nil {
			typex.Fatal(err)
		}
		// Insert
		inserter := benchPost(url, key, maxSize)
		inserter(values)

		// Modify
		record := records.PutRecords{
			Key:     bs.Key(key.Hex()),
			Records: values.PutBody(),
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

		col := defaultSession.DB("db").C(fmt.Sprintf("tickets_%s", key.Hex()))

		var docs []map[string]interface{}
		if err := col.Find(bson.M{
			"_id": bson.M{"$in": values.GetAllFieldIds()},
		}).All(&docs); err != nil {
			typex.Fatal(err)
		}

		if num := len(docs); num != s.Records || num != len(values) {
			return false
		}

		for _, v := range docs {
			if !values.ContainsFieldId(v["_id"].(bson.ObjectId)) {
				return false
			}
		}

		return true
	}
}

func TestPut_ModifyAllReadAll_Persists(t *testing.T) {
	if defaultUseStubs {
		t.Skip("Requires db access for this atm.")
	}

	e := env.New(nil)

	ts, co := setup(e)
	defer tear(ts)

	f := testPutPersistency(ts.URL, defaultMaxSize, co)

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}

// Test Patch

func testPatch(url string,
	maxSize int,
	co *coordinator.Coordinator,
) func(tests.PostBody) bool {
	pool := getIdentPool()
	return func(values tests.PostBody) bool {
		key, err := b.Bson(pool.Get())
		if err != nil {
			typex.Fatal(err)
		}
		// Insert
		inserter := benchPost(url, key, maxSize)
		inserter(values)

		// Patch
		newOwnerId := bson.NewObjectId().Hex()
		record := records.PatchRecords{
			Operations: []records.Operation{
				records.Operation{
					Op:    "match",
					Path:  "/owner_id",
					Value: values.GetOwnerId().Hex(),
				},
				records.Operation{
					Op:    "replace",
					Path:  "/owner_id",
					Value: newOwnerId,
				},
				records.Operation{
					Op:    "match",
					Path:  "/owner_id",
					Value: newOwnerId,
				},
			},
			Score:   2,
			MaxSize: int64(maxSize),
			Expiry:  defaultExpiry,
		}
		bytes, err := record.Write(flatbuffers.NewBuilder(0))
		if err != nil {
			typex.Fatal(err)
		}

		var (
			field = values.GetFirstFieldId().Hex()
			body  = tests.Patch(fmt.Sprintf("%s/http/v1/%s/%s", url, key.Hex(), field), bytes)
		)

		s := &records.OKInt{}
		s.Read(body)

		return s.Records == 2
	}
}

func TestPatch(t *testing.T) {
	e := env.New(nil)

	ts, co := setup(e)
	defer tear(ts)

	f := testPatch(ts.URL, defaultMaxSize, co)

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}

// Test Select

type Header struct {
	Key     bs.Key
	Field   bs.Key
	OwnerId bs.Key
}

type HeadersByField []Header

func (s HeadersByField) Len() int {
	return len(s)
}

func (s HeadersByField) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s HeadersByField) Less(i, j int) bool {
	return s[i].Field.String() < s[j].Field.String()
}

func testSelect(url string,
	co *coordinator.Coordinator,
) func(tests.PostBody) bool {
	var (
		match = func(key bson.ObjectId, values tests.PostBody, other []Header) bool {
			res := make([]Header, 0, len(values))
			for _, v := range values {
				res = append(res, Header{
					Key:     bs.Key(key.Hex()),
					Field:   bs.Key(v.Id.Hex()),
					OwnerId: bs.Key(v.OwnerId.Hex()),
				})
			}

			sort.Sort(HeadersByField(res))
			return reflect.DeepEqual(res, other)
		}
		pool = getIdentPool()
	)
	return func(values tests.PostBody) bool {
		key, err := b.Bson(pool.Get())
		if err != nil {
			typex.Fatal(err)
		}
		// Insert
		inserter := benchPost(url, key, defaultMaxSize)
		inserter(values)

		// Select
		body := tests.Get(fmt.Sprintf("%s/http/v1/%s?limit=99999&size=%d&expiry=%d",
			url,
			key.Hex(),
			defaultMaxSize,
			defaultExpiry.Nanoseconds(),
		))

		s := &records.OKKeyFieldScoreTxnValues{}
		s.Read(body)

		res := make([]Header, 0, len(s.Records))
		for _, v := range s.Records {

			m := &records.Header{}
			m.Read([]byte(v.Value[1:]))

			res = append(res, Header{
				Key:     v.Key,
				Field:   v.Field,
				OwnerId: m.OwnerId,
			})
		}

		sort.Sort(HeadersByField(res))
		return match(key, values, res)
	}
}

func TestSelect_SelectOneReadOne(t *testing.T) {
	e := env.New(nil)
	e.StoreSelectStrategy = "SelectOneReadOne"

	ts, co := setup(e)
	defer tear(ts)

	f := testSelect(ts.URL, co)

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}

func TestSelect_SelectQuorumReadAll(t *testing.T) {
	e := env.New(nil)
	e.StoreSelectStrategy = "SelectQuorumReadAll"

	ts, co := setup(e)
	defer tear(ts)

	f := testSelect(ts.URL, co)

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}

func TestSelect_SelectAllReadAll(t *testing.T) {
	e := env.New(nil)
	e.StoreSelectStrategy = "SelectAllReadAll"

	ts, co := setup(e)
	defer tear(ts)

	f := testSelect(ts.URL, co)

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}

// Test Get

func testGet(url string,
	maxSize int,
	co *coordinator.Coordinator,
) func(tests.PostBody) bool {
	pool := getIdentPool()
	return func(values tests.PostBody) bool {
		key, err := b.Bson(pool.Get())
		if err != nil {
			typex.Fatal(err)
		}
		// Insert
		inserter := benchPost(url, key, maxSize)
		inserter(values)

		// Get
		var (
			field = values.GetFirstFieldId().Hex()
			body  = tests.Get(fmt.Sprintf("%s/http/v1/%s/%s", url, key.Hex(), field))
		)

		s := &records.OKKeyFieldScoreTxnValue{}
		s.Read(body)

		return s.Records.Field.String() == values[0].Id.Hex()
	}
}

func TestGet(t *testing.T) {
	e := env.New(nil)

	ts, co := setup(e)
	defer tear(ts)

	f := testGet(ts.URL, defaultMaxSize, co)

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}

// Test Delete

type deleteOpts struct {
	amountToDelete float64
}

func deleteAmount(opts *deleteOpts, values tests.PostBody) int {
	amount := len(values)
	if opts != nil {
		amount = int(math.Ceil(float64(amount) * opts.amountToDelete))
	}
	return amount
}

func testDelete(url string,
	co *coordinator.Coordinator,
	opts *deleteOpts,
) func(tests.PostBody) int {
	pool := getIdentPool()
	return func(values tests.PostBody) int {
		key, err := b.Bson(pool.Get())
		if err != nil {
			typex.Fatal(err)
		}

		return benchDelete(url, key, opts)(values)
	}
}

func TestDelete_DeleteAllReadAll(t *testing.T) {
	e := env.New(nil)
	e.StoreDeleteStrategy = "DeleteAllReadAll"

	ts, co := setup(e)
	defer tear(ts)

	var (
		f = testDelete(ts.URL, co, nil)
		g = func(values tests.PostBody) bool {
			return f(values) == deleteAmount(nil, values)
		}
	)

	if err := quick.Check(g, tests.Config()); err != nil {
		t.Error(err)
	}
}

func testDeleteThenSelect(url string,
	co *coordinator.Coordinator,
	opts *deleteOpts,
) func(tests.PostBody) bool {
	var (
		validator = func(values tests.PostBody, opts *deleteOpts) int {
			return deleteAmount(opts, values)
		}
		pool = getIdentPool()
	)
	return func(values tests.PostBody) bool {
		key, err := b.Bson(pool.Get())
		if err != nil {
			typex.Fatal(err)
		}
		// Insert + Delete
		deleter := benchDelete(url, key, opts)
		if res := deleter(values); res != validator(values, opts) {
			typex.Fatalf("Unsucessful delete action (%d, %d)", res, deleteAmount(opts, values))
		}

		// Select
		body := tests.Get(fmt.Sprintf("%s/http/v1/%s?limit=99999&size=%d&expiry=%d",
			url,
			key.Hex(),
			defaultMaxSize,
			defaultExpiry.Nanoseconds(),
		))

		s := &records.OKKeyFieldScoreTxnValues{}
		if err := s.Read(body); err != nil {
			typex.Fatal(err)
		}

		return len(s.Records) == (len(values) - deleteAmount(opts, values))
	}
}

func TestDeleteAll_DeleteAllReadAll(t *testing.T) {
	e := env.New(nil)
	e.StoreDeleteStrategy = "DeleteAllReadAll"

	ts, co := setup(e)
	defer tear(ts)

	f := testDeleteThenSelect(ts.URL, co, nil)

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}

func TestDeleteHalf_DeleteAllReadAll(t *testing.T) {
	e := env.New(nil)
	e.StoreDeleteStrategy = "DeleteAllReadAll"

	ts, co := setup(e)
	defer tear(ts)

	f := testDeleteThenSelect(ts.URL, co, &deleteOpts{0.5})

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}

// Test Rollback

type rollbackOpts struct {
	amountToRollback float64
}

func rollbackAmount(opts *rollbackOpts, values tests.PostBody) int {
	amount := len(values)
	if opts != nil {
		amount = int(math.Ceil(float64(amount) * opts.amountToRollback))
	}
	return amount
}

func testRollback(url string,
	co *coordinator.Coordinator,
	opts *rollbackOpts,
) func(tests.PostBody) int {
	pool := getIdentPool()
	return func(values tests.PostBody) int {
		key, err := b.Bson(pool.Get())
		if err != nil {
			typex.Fatal(err)
		}

		return benchRollback(url, key, opts)(values)
	}
}

func TestRollback_DeleteAllReadAll(t *testing.T) {
	e := env.New(nil)

	ts, co := setup(e)
	defer tear(ts)

	var (
		f = testRollback(ts.URL, co, nil)
		g = func(values tests.PostBody) bool {
			return f(values) > 0
		}
	)

	if err := quick.Check(g, tests.Config()); err != nil {
		t.Error(err)
	}
}

func testRollbackThenSelect(url string,
	co *coordinator.Coordinator,
	opts *rollbackOpts,
) func(tests.PostBody) bool {
	var (
		validator = func(values tests.PostBody) int {
			return rollbackAmount(opts, values)
		}
		pool = getIdentPool()
	)
	return func(values tests.PostBody) bool {
		key, err := b.Bson(pool.Get())
		if err != nil {
			typex.Fatal(err)
		}
		// Insert + Rollback
		rollbackr := benchRollback(url, key, opts)
		if res := rollbackr(values); res != validator(values) {
			typex.Fatalf("Unsucessful rollback action (%d, %d)", res, rollbackAmount(opts, values))
		}

		// Select
		body := tests.Get(fmt.Sprintf("%s/http/v1/%s?limit=99999&size=%d&expiry=%d",
			url,
			key.Hex(),
			defaultMaxSize,
			defaultExpiry.Nanoseconds(),
		))

		s := &records.OKKeyFieldScoreTxnValues{}
		s.Read(body)

		return len(s.Records) == (len(values) - rollbackAmount(opts, values))
	}
}

func TestRollbackAll_DeleteAllReadAll(t *testing.T) {
	e := env.New(nil)

	ts, co := setup(e)
	defer tear(ts)

	f := testRollbackThenSelect(ts.URL, co, nil)

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}

func TestRollbackHalf_DeleteAllReadAll(t *testing.T) {
	e := env.New(nil)

	ts, co := setup(e)
	defer tear(ts)

	f := testRollbackThenSelect(ts.URL, co, &rollbackOpts{0.5})

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}

// Test Query

func testQuery(url string,
	maxSize int,
	co *coordinator.Coordinator,
) func(tests.PostBody) bool {
	pool := getIdentPool()
	return func(values tests.PostBody) bool {
		key, err := b.Bson(pool.Get())
		if err != nil {
			typex.Fatal(err)
		}
		// Insert
		inserter := benchPost(url, key, maxSize)
		inserter(values)

		// Query
		var (
			body = tests.Get(fmt.Sprintf("%s/http/v1/%s/query?owner_id=%s&size=%d&expiry=%d",
				url,
				key.Hex(),
				values.GetOwnerId().Hex(),
				defaultMaxSize,
				defaultExpiry,
			))
		)

		s := &records.OKQuery{}
		s.Read(body)

		return len(s.Records) == len(values)
	}
}

func TestQuery(t *testing.T) {
	e := env.New(nil)

	ts, co := setup(e)
	defer tear(ts)

	f := testQuery(ts.URL, defaultMaxSize, co)

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}
}

// Test KV store

func testKVStore(url string,
	maxSize int,
	setup func(bson.ObjectId),
	equality func(bson.ObjectId, tests.PostBody) bool,
) func(tests.PostBody) bool {
	pool := getIdentPool()
	return func(values tests.PostBody) bool {
		key, err := b.Bson(pool.Get())
		if err != nil {
			typex.Fatal(err)
		}
		setup(key)

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

		return equality(key, values)
	}
}

func testKVStoreWithBank(t *testing.T, node int) {
	if defaultUseStubs {
		t.Skip("Skipping...")
	}

	e := env.New(nil)
	e.StoreInsertStrategy = "InsertAllReadAll"
	e.StorePoolRoutingStrategy = "KeyStore"

	ts, co := setup(e)
	defer tear(ts)

	f := testKVStore(ts.URL, defaultMaxSize,
		func(key bson.ObjectId) {
			var (
				kvKey   = fmt.Sprintf("%s_%s", e.StoreKeyStorePrefix, key.Hex())
				kvBytes = tests.MustMarshal(client.KeyValueStoreElement{
					Key:  key.Hex(),
					Node: node,
				})
			)

			kvs := defaultConsul.KV()
			if _, err := kvs.Put(&api.KVPair{Key: kvKey, Value: kvBytes}, nil); err != nil {
				typex.Fatal(err)
			}

			// Access things which shouldn't be accessed!
			accessor := coordinator.NewCoordinatorAccessor(co)
			accessor.StoreOpts().KeyStoreTicker <- struct{}{}

			// Sleep for the request
			time.Sleep(time.Millisecond * 40)
			runtime.Gosched()
		},
		func(key bson.ObjectId, values tests.PostBody) bool {
			//	1. Check the redis farm for it!
			res := true

			for _, v := range values {
				var (
					hash  = fmt.Sprintf("s:%s+", key.Hex())
					field = v.Id.Hex()
				)
				bytes, err := redis.Bytes(defaultRedi[node].Do("HGET", hash, field))
				if err != nil {
					alt := 0
					if node == 0 {
						alt = 1
					}
					b, e := redis.Bytes(defaultRedi[alt].Do("HGET", hash, field))
					if e == nil && b != nil {
						fmt.Printf("Hash and Field stored in other bank!! %d - %s : %s\n", alt, hash, field)
					}
					typex.Fatalf("Redis error: %s", err.Error())
				}

				score, _, _, value, err := store.ExtractScoreTxnExpiryValue(string(bytes))
				if err != nil {
					typex.Fatalf("Extract error: %s", err.Error())
				}

				record := records.PostRecord{}
				body, err := records.ReadBody(value)
				if err != nil {
					typex.Fatalf("Read body error: %s", err.Error())
				}
				if err := record.Read(body); err != nil {
					typex.Fatalf("Record body error: %s", err.Error())
				}

				res = res && (score == 1) && (record.OwnerId.Hex() == v.OwnerId.Hex())
			}

			return res
		},
	)

	if err := quick.Check(f, tests.Config()); err != nil {
		t.Error(err)
	}

}

func TestKVStore_InsertAllReadAll_StoresBank0(t *testing.T) {
	testKVStoreWithBank(t, 0)
}

func TestKVStore_InsertAllReadAll_StoresBank1(t *testing.T) {
	testKVStoreWithBank(t, 1)
}

// Test End-to-end

func testConsumption(url string,
	maxSize int,
	co *coordinator.Coordinator,
) func(bson.ObjectId, int) {
	var (
		request = func(url string, key bson.ObjectId, amount, maxSize int) []records.PostRecord {
			values := generatePostBody(amount)

			// Insert
			inserter := benchPost(url, key, maxSize)
			res := inserter(values)

			if num := len(values); res != num {
				typex.Fatal("Request amount does not match response amount: ", num, res)
			}

			return values
		}
		assign = func(url string, key bson.ObjectId, body []records.PostRecord, maxSize int) []records.PutRecord {
			values := tests.PostBody(body).PutBody()

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

			response := tests.Put(fmt.Sprintf("%s/http/v1/%s", url, key.Hex()), bytes)

			s := &records.OKInt{}
			s.Read(response)

			if num := len(values); s.Records != num {
				typex.Fatal("Assign amount does not match response amount: ", num, s.Records)
			}

			return values
		}
	)

	return func(key bson.ObjectId, amount int) {
		var (
			a = request(url, key, amount, maxSize)
			b = assign(url, key, a, maxSize)

			aIds = make([]string, 0, len(a))
			bIds = make([]string, 0, len(b))
		)

		for _, v := range a {
			aIds = append(aIds, v.Id.Hex())
		}

		for _, v := range b {
			bIds = append(bIds, v.Id.Hex())
		}

		sort.Strings(aIds)
		sort.Strings(bIds)

		if !reflect.DeepEqual(aIds, bIds) {
			typex.Fatal("Request and Assign ids don't match!", aIds, bIds)
		}
	}
}

func generatePostBody(amount int) tests.PostBody {
	return tests.PostBody(nil).Make(tests.Random, amount)
}

func testConsumptionWithAmount(t *testing.T,
	ts *httptest.Server,
	co *coordinator.Coordinator,
	total, amount, assertLeft, overflow int,
) {
	var (
		f   = testConsumption(ts.URL, total, co)
		key = bson.NewObjectId()
	)

	for i := 0; i < (total/amount)+overflow; i++ {
		f(key, amount)
	}

	size, err := co.Size(bs.Key(key.Hex()))
	if err != nil {
		typex.Fatal(err)
	}

	if size != total-assertLeft {
		t.Errorf("Size should be %d, but was %d", total-assertLeft, size)
	}
}

func TestConsumption_10_1(t *testing.T) {
	t.Parallel()

	e := env.New(nil)

	ts, co := setup(e)
	defer tear(ts)

	testConsumptionWithAmount(t, ts, co, 10, 1, 0, 0)
}

func TestConsumption_100_1(t *testing.T) {
	t.Parallel()

	e := env.New(nil)

	ts, co := setup(e)
	defer tear(ts)

	testConsumptionWithAmount(t, ts, co, 100, 1, 0, 0)
}

func TestConsumption_1000_1(t *testing.T) {
	t.Parallel()

	e := env.New(nil)

	ts, co := setup(e)
	defer tear(ts)

	testConsumptionWithAmount(t, ts, co, 1000, 1, 0, 0)
}

func TestConsumption_10_2(t *testing.T) {
	t.Parallel()

	e := env.New(nil)

	ts, co := setup(e)
	defer tear(ts)

	testConsumptionWithAmount(t, ts, co, 10, 2, 0, 0)
}

func TestConsumption_100_2(t *testing.T) {
	t.Parallel()

	e := env.New(nil)

	ts, co := setup(e)
	defer tear(ts)

	testConsumptionWithAmount(t, ts, co, 100, 2, 0, 0)
}

func TestConsumption_1000_2(t *testing.T) {
	t.Parallel()

	e := env.New(nil)

	ts, co := setup(e)
	defer tear(ts)

	testConsumptionWithAmount(t, ts, co, 1000, 2, 0, 0)
}

func TestConsumption_10_3(t *testing.T) {
	t.Parallel()

	e := env.New(nil)

	ts, co := setup(e)
	defer tear(ts)

	testConsumptionWithAmount(t, ts, co, 10, 3, 1, 0)
}
