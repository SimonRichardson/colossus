package records

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/schemas/schema"
	flatbuffers "github.com/google/flatbuffers/go"
)

type RollbackRecords struct {
	Key     bs.Key
	Records []RollbackRecord
	Score   float64
	MaxSize int64
	Expiry  time.Duration
}

func (r RollbackRecords) Write(fb *flatbuffers.Builder) ([]byte, error) {
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

	schema.RollbackRequestStartRecordsVector(fb, num)

	for _, v := range positions {
		fb.PrependUOffsetT(v)
	}

	vector := fb.EndVector(num)

	schema.RollbackRequestStart(fb)
	schema.RollbackRequestAddScore(fb, r.Score)
	schema.RollbackRequestAddMaxSize(fb, uint64(r.MaxSize))
	schema.RollbackRequestAddExpiry(fb, uint64(r.Expiry.Nanoseconds()))
	schema.RollbackRequestAddRecords(fb, vector)
	position := schema.RollbackRequestEnd(fb)

	fb.Finish(position)
	return fb.FinishedBytes(), nil
}

type RollbackRecord struct {
	Id                     bson.ObjectId
	Updated                time.Time
	OwnerId, TransactionId bson.ObjectId
}

// WriteRollbackRecord represents away of writing a RollbackRecord to a byte buffer
func (r RollbackRecord) Write(fb *flatbuffers.Builder) ([]byte, error) {
	position, err := r.WriteSub(fb)
	if err != nil {
		return nil, err
	}

	fb.Finish(position)
	return fb.FinishedBytes(), nil
}

func (r RollbackRecord) WriteSub(fb *flatbuffers.Builder) (flatbuffers.UOffsetT, error) {
	var (
		idPosition, ownerIdPosition, transactionIdPosition flatbuffers.UOffsetT
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

	now := time.Now()

	schema.RollbackRecordStart(fb)
	schema.RollbackRecordAddTyp(fb, schema.TypeRollback)
	schema.RollbackRecordAddId(fb, idPosition)
	schema.RollbackRecordAddUpdated(fb, uint64(now.UnixNano()))
	schema.RollbackRecordAddOwnerId(fb, ownerIdPosition)
	schema.RollbackRecordAddTransactionId(fb, transactionIdPosition)

	return schema.RollbackRecordEnd(fb), nil
}

func (r *RollbackRecord) Read(bytes []byte) error {
	if len(bytes) < 1 {
		return ErrInvalidLength(0)
	}

	obj, err := GetRootAsRollbackRecord(bytes, 0)
	if err != nil {
		return err
	}

	var (
		id            = string(obj.Id(nil).Hex())
		ownerId       = string(obj.OwnerId(nil).Hex())
		transactionId = string(obj.TransactionId(nil).Hex())
	)

	if !bson.IsObjectIdHex(id) || !bson.IsObjectIdHex(ownerId) {
		return ErrInvalidIdHex(0)
	}

	r.Id = bson.ObjectIdHex(id)
	r.OwnerId = bson.ObjectIdHex(ownerId)
	r.Updated = time.Unix(0, int64(obj.Updated()))
	r.TransactionId = bson.ObjectIdHex(transactionId)

	return nil
}

func RollbackRecordFromSchemaToByte(fb *flatbuffers.Builder, record *schema.RollbackRecord) ([]byte, error) {
	var (
		id            = record.Id(nil).Hex()
		ownerId       = record.OwnerId(nil).Hex()
		transactionId = record.TransactionId(nil).Hex()
	)

	value := RollbackRecord{
		Id:            bson.ObjectIdHex(string(id)),
		Updated:       time.Unix(0, int64(record.Updated())),
		OwnerId:       bson.ObjectIdHex(string(ownerId)),
		TransactionId: bson.ObjectIdHex(string(transactionId)),
	}

	return value.Write(fb)
}

func GetRootAsRollbackRecord(buf []byte, offset flatbuffers.UOffsetT) (*schema.RollbackRecord, error) {
	if len(buf) <= int(offset) {
		return nil, ErrInvalidLength(1)
	}

	var (
		n = flatbuffers.GetUOffsetT(buf[offset:])
		x = &schema.RollbackRecord{}
	)

	if len(buf) < int(n+offset) {
		return nil, ErrInvalidLength(2)
	}

	x.Init(buf, n+offset)
	return x, nil
}

func PackageRollbackRecord(buf []byte) string {
	return fmt.Sprintf("%d", schema.TypeRollback) + string(buf)
}
