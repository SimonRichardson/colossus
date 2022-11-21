package records

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/schemas/schema"
	flatbuffers "github.com/google/flatbuffers/go"
)

type DeleteRecords struct {
	Key     bs.Key
	Records []DeleteRecord
	Score   float64
	MaxSize int64
	Expiry  time.Duration
}

func (r DeleteRecords) Write(fb *flatbuffers.Builder) ([]byte, error) {
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

	schema.DeleteRequestStartRecordsVector(fb, num)

	for _, v := range positions {
		fb.PrependUOffsetT(v)
	}

	vector := fb.EndVector(num)

	schema.DeleteRequestStart(fb)
	schema.DeleteRequestAddScore(fb, r.Score)
	schema.DeleteRequestAddMaxSize(fb, uint64(r.MaxSize))
	schema.DeleteRequestAddExpiry(fb, uint64(r.Expiry.Nanoseconds()))
	schema.DeleteRequestAddRecords(fb, vector)
	position := schema.DeleteRequestEnd(fb)

	fb.Finish(position)
	return fb.FinishedBytes(), nil
}

type DeleteRecord struct {
	Id                     bson.ObjectId
	Updated                time.Time
	OwnerId, TransactionId bson.ObjectId
}

// WriteDeleteRecord represents away of writing a DeleteRecord to a byte buffer
func (r DeleteRecord) Write(fb *flatbuffers.Builder) ([]byte, error) {
	position, err := r.WriteSub(fb)
	if err != nil {
		return nil, err
	}

	fb.Finish(position)
	return fb.FinishedBytes(), nil
}

func (r DeleteRecord) WriteSub(fb *flatbuffers.Builder) (flatbuffers.UOffsetT, error) {
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

	schema.DeleteRecordStart(fb)
	schema.DeleteRecordAddTyp(fb, schema.TypeDelete)
	schema.DeleteRecordAddId(fb, idPosition)
	schema.DeleteRecordAddUpdated(fb, uint64(now.UnixNano()))
	schema.DeleteRecordAddOwnerId(fb, ownerIdPosition)
	schema.DeleteRecordAddTransactionId(fb, transactionIdPosition)

	return schema.DeleteRecordEnd(fb), nil
}

func (r *DeleteRecord) Read(bytes []byte) error {
	if len(bytes) < 1 {
		return ErrInvalidLength(0)
	}

	obj, err := GetRootAsDeleteRecord(bytes, 0)
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

func DeleteRecordFromSchemaToByte(fb *flatbuffers.Builder, record *schema.DeleteRecord) ([]byte, error) {
	var (
		id            = record.Id(nil).Hex()
		ownerId       = record.OwnerId(nil).Hex()
		transactionId = record.TransactionId(nil).Hex()
	)

	value := DeleteRecord{
		Id:            bson.ObjectIdHex(string(id)),
		Updated:       time.Unix(0, int64(record.Updated())),
		OwnerId:       bson.ObjectIdHex(string(ownerId)),
		TransactionId: bson.ObjectIdHex(string(transactionId)),
	}

	return value.Write(fb)
}

func GetRootAsDeleteRecord(buf []byte, offset flatbuffers.UOffsetT) (*schema.DeleteRecord, error) {
	if len(buf) <= int(offset) {
		return nil, ErrInvalidLength(1)
	}

	var (
		n = flatbuffers.GetUOffsetT(buf[offset:])
		x = &schema.DeleteRecord{}
	)

	if len(buf) < int(n+offset) {
		return nil, ErrInvalidLength(2)
	}

	x.Init(buf, n+offset)
	return x, nil
}

func PackageDeleteRecord(buf []byte) string {
	return fmt.Sprintf("%d", schema.TypeDelete) + string(buf)
}
