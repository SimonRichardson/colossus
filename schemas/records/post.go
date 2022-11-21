package records

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/SimonRichardson/colossus/schemas/schema"
	flatbuffers "github.com/google/flatbuffers/go"
)

type PostRecords struct {
	Records []PostRecord
	Score   float64
	MaxSize int64
	Expiry  time.Duration
}

func (r PostRecords) Write(fb *flatbuffers.Builder) ([]byte, error) {
	var (
		num       = len(r.Records)
		positions = make([]flatbuffers.UOffsetT, 0, num)
	)

	for _, v := range r.Records {
		position, err := v.WriteSub(fb)
		if err != nil {
			return nil, err
		}
		positions = append(positions, position)
	}

	schema.PostRequestStartRecordsVector(fb, num)

	for _, v := range positions {
		fb.PrependUOffsetT(v)
	}

	vector := fb.EndVector(num)

	schema.PostRequestStart(fb)
	schema.PostRequestAddScore(fb, r.Score)
	schema.PostRequestAddMaxSize(fb, uint64(r.MaxSize))
	schema.PostRequestAddExpiry(fb, uint64(r.Expiry.Nanoseconds()))
	schema.PostRequestAddRecords(fb, vector)
	position := schema.PostRequestEnd(fb)

	fb.Finish(position)
	return fb.FinishedBytes(), nil
}

type PostRecord struct {
	Id                        bson.ObjectId
	Updated, Reserved, Expiry time.Time
	Cost                      Cost
	OwnerId, TransactionId    bson.ObjectId
}

// WritePostRecord represents away of writing a PostRecord to a byte buffer
func (r PostRecord) Write(fb *flatbuffers.Builder) ([]byte, error) {
	position, err := r.WriteSub(fb)
	if err != nil {
		return nil, err
	}

	fb.Finish(position)
	return fb.FinishedBytes(), nil
}

func (r PostRecord) WriteSub(fb *flatbuffers.Builder) (flatbuffers.UOffsetT, error) {
	var (
		idPosition, ownerIdPosition, transactionIdPosition flatbuffers.UOffsetT
		costPosition                                       flatbuffers.UOffsetT
		err                                                error
	)

	if idPosition, err = MakeId(r.Id.Hex()).WriteSub(fb); err != nil {
		return 0, err
	}
	if ownerIdPosition, err = MakeId(r.OwnerId.Hex()).WriteSub(fb); err != nil {
		return 0, err
	}
	if transactionIdPosition, err = MakeId(r.TransactionId.Hex()).WriteSub(fb); err != nil {
		return 0, err
	}
	if costPosition, err = r.Cost.WriteSub(fb); err != nil {
		return 0, err
	}

	now := time.Now()

	schema.PostRecordStart(fb)
	schema.PostRecordAddTyp(fb, schema.TypePost)
	schema.PostRecordAddId(fb, idPosition)
	schema.PostRecordAddUpdated(fb, uint64(now.UnixNano()))
	schema.PostRecordAddExpiry(fb, uint64(r.Expiry.UnixNano()))
	schema.PostRecordAddReserved(fb, uint64(now.UnixNano()))
	schema.PostRecordAddCost(fb, costPosition)
	schema.PostRecordAddOwnerId(fb, ownerIdPosition)
	schema.PostRecordAddTransactionId(fb, transactionIdPosition)

	return schema.PostRecordEnd(fb), nil
}

func (r *PostRecord) Read(bytes []byte) error {
	if len(bytes) < 1 {
		return ErrInvalidLength(7)
	}

	obj, err := GetRootAsPostRecord(bytes, 0)
	if err != nil {
		return err
	}

	var (
		id            = string(obj.Id(nil).Hex())
		ownerId       = string(obj.OwnerId(nil).Hex())
		transactionId = string(obj.TransactionId(nil).Hex())
	)

	if !bson.IsObjectIdHex(id) || !bson.IsObjectIdHex(ownerId) || !bson.IsObjectIdHex(transactionId) {
		return ErrInvalidIdHex(1)
	}

	r.Id = bson.ObjectIdHex(id)
	r.OwnerId = bson.ObjectIdHex(ownerId)
	r.Updated = time.Unix(0, int64(obj.Updated()))
	r.Reserved = time.Unix(0, int64(obj.Reserved()))
	r.Expiry = time.Unix(0, int64(obj.Expiry()))
	r.TransactionId = bson.ObjectIdHex(transactionId)

	var (
		cost     = obj.Cost(nil)
		currency = string(cost.Currency())
	)

	if len(currency) < 1 {
		return ErrInvalidLength(8)
	}

	r.Cost = Cost{
		Currency: currency,
		Price:    cost.Price(),
	}

	return nil
}

func PostRecordFromSchemaToByte(fb *flatbuffers.Builder, record *schema.PostRecord) ([]byte, error) {
	var (
		id            = record.Id(nil).Hex()
		ownerId       = record.OwnerId(nil).Hex()
		transactionId = record.TransactionId(nil).Hex()
		cost          = record.Cost(nil)
	)

	value := PostRecord{
		Id:       bson.ObjectIdHex(string(id)),
		Updated:  time.Unix(0, int64(record.Updated())),
		Reserved: time.Unix(0, int64(record.Reserved())),
		Expiry:   time.Unix(0, int64(record.Expiry())),
		Cost: Cost{
			Currency: string(cost.Currency()),
			Price:    cost.Price(),
		},
		OwnerId:       bson.ObjectIdHex(string(ownerId)),
		TransactionId: bson.ObjectIdHex(string(transactionId)),
	}

	return value.Write(fb)
}

func GetRootAsPostRecord(buf []byte, offset flatbuffers.UOffsetT) (*schema.PostRecord, error) {
	if len(buf) <= int(offset) {
		return nil, ErrInvalidLength(9)
	}

	var (
		n = flatbuffers.GetUOffsetT(buf[offset:])
		x = &schema.PostRecord{}
	)

	if len(buf) < int(n+offset) {
		return nil, ErrInvalidLength(10)
	}

	x.Init(buf, n+offset)
	return x, nil
}

func PackagePostRecord(buf []byte) string {
	return fmt.Sprintf("%d", schema.TypePost) + string(buf)
}
