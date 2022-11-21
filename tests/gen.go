package tests

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/SimonRichardson/colossus/schemas/records"
	"github.com/SimonRichardson/colossus/typex"
)

type RandString []string

func (s RandString) Generate(r *rand.Rand, size int) reflect.Value {
	var (
		chars        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		randomString = func(r *rand.Rand, size int) string {
			result := make([]byte, size)
			for i := 0; i < size; i++ {
				result[i] = chars[r.Intn(len(chars)-1)]
			}
			return string(result)
		}
		m = []string{randomString(r, size)}
	)

	return reflect.ValueOf(m)
}

func (s RandString) String() string {
	return s[0]
}

func generateString(rand *rand.Rand, size int) string {
	if size < 1 {
		size = 0
	}
	return RandString(RandString{}.Generate(rand, size).Interface().([]string)).String()
}

type PostBody []records.PostRecord

func (b PostBody) Make(rand *rand.Rand, size int) PostBody {
	if size < 1 {
		size = 1
	}

	var (
		now      = time.Now()
		owner_id = bson.NewObjectId()
		m        = make([]records.PostRecord, 0, size)
	)
	for i := 0; i < size; i++ {
		m = append(m, records.PostRecord{
			Id:       bson.NewObjectId(),
			Updated:  now,
			Reserved: now,
			Expiry:   now.Add(time.Minute * time.Duration(int64((rand.Intn(5) + 5)))),
			Cost: records.Cost{
				Currency: "GBP",
				Price:    uint64(rand.Intn(100000)),
			},
			OwnerId:       owner_id,
			TransactionId: bson.NewObjectId(),
		})
	}

	return m
}

func (b PostBody) Generate(rand *rand.Rand, size int) reflect.Value {
	body := PostBody(nil)
	m := body.Make(rand, randomSize(rand, size))
	return reflect.ValueOf(m)
}

func (b PostBody) PutBody() PutBody {
	var (
		now = time.Now()
		m   = make([]records.PutRecord, 0, len(b))
	)
	for _, v := range b {
		r := records.PutRecord{
			Id:        v.Id,
			Updated:   v.Updated,
			Purchased: now,
			OwnerId:   v.OwnerId,
			EventCost: records.Cost{
				Currency: "GBP",
				Price:    1,
			},
			EventDates: records.Dates{
				Start: uint64(now.UnixNano()),
				End:   uint64(now.UnixNano()),
			},
			TransactionId: v.TransactionId,
			Codes: records.Codes{
				BarcodeType:   "barcode-39",
				BarcodeOrigin: "xxxx",
				BarcodeSource: "xxxx",
			},
		}
		m = append(m, r)
	}
	return m
}

func (b PostBody) DeleteBody() DeleteBody {
	m := make([]records.DeleteRecord, 0, len(b))
	for _, v := range b {
		r := records.DeleteRecord{
			Id:            v.Id,
			Updated:       v.Updated,
			OwnerId:       v.OwnerId,
			TransactionId: v.TransactionId,
		}
		m = append(m, r)
	}
	return m
}

func (b PostBody) RollbackBody() RollbackBody {
	m := make([]records.RollbackRecord, 0, len(b))
	for _, v := range b {
		r := records.RollbackRecord{
			Id:            v.Id,
			Updated:       v.Updated,
			OwnerId:       v.OwnerId,
			TransactionId: v.TransactionId,
		}
		m = append(m, r)
	}
	return m
}

func (b PostBody) GetOwnerId() bson.ObjectId {
	if len(b) < 1 {
		typex.Fatal(fmt.Errorf("No records"))
	}
	return b[0].OwnerId
}

func (b PostBody) GetFirstFieldId() bson.ObjectId {
	if len(b) < 1 {
		typex.Fatal(fmt.Errorf("No records"))
	}
	return b[0].Id
}

func (b PostBody) GetAllFieldIds() []bson.ObjectId {
	if len(b) < 1 {
		typex.Fatal(fmt.Errorf("No records"))
	}
	var res []bson.ObjectId
	for _, v := range b {
		res = append(res, v.Id)
	}
	return res
}

func (b PostBody) ContainsFieldId(id bson.ObjectId) bool {
	for _, v := range b {
		if v.Id == id {
			return true
		}
	}
	return false
}

type PutBody []records.PutRecord

type DeleteBody []records.DeleteRecord

type RollbackBody []records.RollbackRecord

func randomSize(rand *rand.Rand, size int) int {
	if size < 1 {
		return 1
	}
	return rand.Intn(size) + size
}
