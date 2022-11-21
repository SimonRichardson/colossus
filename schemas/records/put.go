package records

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/schemas/schema"
	flatbuffers "github.com/google/flatbuffers/go"
)

type PutRecords struct {
	Key     bs.Key
	Records []PutRecord
	Score   float64
	MaxSize int64
	Expiry  time.Duration
}

func (r PutRecords) Write(fb *flatbuffers.Builder) ([]byte, error) {
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

	schema.PutRequestStartRecordsVector(fb, num)

	for _, v := range positions {
		fb.PrependUOffsetT(v)
	}

	vector := fb.EndVector(num)

	schema.PutRequestStart(fb)
	schema.PutRequestAddScore(fb, r.Score)
	schema.PutRequestAddMaxSize(fb, uint64(r.MaxSize))
	schema.PutRequestAddExpiry(fb, uint64(r.Expiry.Nanoseconds()))
	schema.PutRequestAddRecords(fb, vector)
	position := schema.PutRequestEnd(fb)

	fb.Finish(position)
	return fb.FinishedBytes(), nil
}

type PutRecord struct {
	Id                 bson.ObjectId
	Updated, Purchased time.Time
	EventCost          Cost
	EventDates         Dates
	OwnerId            bson.ObjectId
	TransactionId      bson.ObjectId
	Codes              Codes
}

// WritePutRecord represents away of writing a PutRecord to a byte buffer
func (r PutRecord) Write(fb *flatbuffers.Builder) ([]byte, error) {
	position, err := r.WriteSub(fb)
	if err != nil {
		return nil, err
	}

	fb.Finish(position)
	return fb.FinishedBytes(), nil
}

func (r PutRecord) WriteSub(fb *flatbuffers.Builder) (flatbuffers.UOffsetT, error) {
	var (
		idPosition, ownerIdPosition, transactionIdPosition flatbuffers.UOffsetT
		costPosition, datesPosition                        flatbuffers.UOffsetT
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
	if costPosition, err = r.EventCost.WriteSub(fb); err != nil {
		return 0, err
	}
	if datesPosition, err = r.EventDates.WriteSub(fb); err != nil {
		return 0, err
	}

	now := time.Now()

	schema.PutRecordStart(fb)
	schema.PutRecordAddTyp(fb, schema.TypePut)
	schema.PutRecordAddId(fb, idPosition)
	schema.PutRecordAddUpdated(fb, uint64(now.UnixNano()))
	schema.PutRecordAddPurchased(fb, uint64(now.UnixNano()))
	schema.PutRecordAddOwnerId(fb, ownerIdPosition)
	schema.PutRecordAddEventCost(fb, costPosition)
	schema.PutRecordAddEventDates(fb, datesPosition)
	schema.PutRecordAddTransactionId(fb, transactionIdPosition)

	return schema.PutRecordEnd(fb), nil
}

func (r *PutRecord) Read(bytes []byte) error {
	if len(bytes) < 1 {
		return ErrInvalidLength(11)
	}

	obj, err := GetRootAsPutRecord(bytes, 0)
	if err != nil {
		return err
	}

	var (
		id            = string(obj.Id(nil).Hex())
		ownerId       = string(obj.OwnerId(nil).Hex())
		transactionId = string(obj.TransactionId(nil).Hex())
	)

	if !bson.IsObjectIdHex(id) || !bson.IsObjectIdHex(ownerId) || !bson.IsObjectIdHex(transactionId) {
		return ErrInvalidIdHex(2)
	}

	r.Id = bson.ObjectIdHex(id)
	r.OwnerId = bson.ObjectIdHex(ownerId)
	r.Updated = time.Unix(0, int64(obj.Updated()))
	r.Purchased = time.Unix(0, int64(obj.Purchased()))
	r.TransactionId = bson.ObjectIdHex(transactionId)

	var (
		cost     = obj.EventCost(nil)
		currency = string(cost.Currency())
	)

	if len(currency) < 1 {
		return ErrInvalidLength(12)
	}

	r.EventCost = Cost{
		Currency: string(cost.Currency()),
		Price:    cost.Price(),
	}

	dates := obj.EventDates(nil)
	r.EventDates = Dates{
		Start: dates.Start(),
		End:   dates.End(),
	}

	return nil
}

func PutRecordFromSchemaToByte(fb *flatbuffers.Builder, record *schema.PutRecord) ([]byte, error) {
	var (
		id            = record.Id(nil).Hex()
		ownerId       = record.OwnerId(nil).Hex()
		cost          = record.EventCost(nil)
		dates         = record.EventDates(nil)
		transactionId = record.TransactionId(nil).Hex()
	)

	value := PutRecord{
		Id:        bson.ObjectIdHex(string(id)),
		Updated:   time.Unix(0, int64(record.Updated())),
		Purchased: time.Unix(0, int64(record.Purchased())),
		OwnerId:   bson.ObjectIdHex(string(ownerId)),
		EventCost: Cost{
			Currency: string(cost.Currency()),
			Price:    cost.Price(),
		},
		EventDates: Dates{
			Start: dates.Start(),
			End:   dates.End(),
		},
		TransactionId: bson.ObjectIdHex(string(transactionId)),
	}

	return value.Write(fb)
}

func GetRootAsPutRecord(buf []byte, offset flatbuffers.UOffsetT) (*schema.PutRecord, error) {
	if len(buf) <= int(offset) {
		return nil, ErrInvalidLength(13)
	}

	var (
		n = flatbuffers.GetUOffsetT(buf[offset:])
		x = &schema.PutRecord{}
	)

	if len(buf) < int(n+offset) {
		return nil, ErrInvalidLength(14)
	}

	x.Init(buf, n+offset)
	return x, nil
}

func PackagePutRecord(buf []byte) string {
	return fmt.Sprintf("%d", schema.TypePut) + string(buf)
}
