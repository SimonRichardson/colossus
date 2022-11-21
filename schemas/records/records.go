package records

import (
	"encoding/json"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/errors"
	"github.com/SimonRichardson/colossus/schemas/schema"
	"github.com/SimonRichardson/colossus/selectors"
	"github.com/SimonRichardson/colossus/typex"
	flatbuffers "github.com/google/flatbuffers/go"
)

var (
	ErrInvalidIdHex = func(index int) error {
		return typex.Errorf(errors.Source, errors.UnexpectedArgument,
			"Invalid Id Hex: (%d)", index)
	}
	ErrInvalidLength = func(index int) error {
		return typex.Errorf(errors.Source, errors.UnexpectedArgument,
			"Invalid length: (%d)", index)
	}
	ErrInvalidRecord = func(index int) error {
		return typex.Errorf(errors.Source, errors.UnexpectedArgument,
			"Invalid record: (%d)", index)
	}
)

// Write defines a common way to write records to a slice bytes
type Write interface {
	Write(fb *flatbuffers.Builder) ([]byte, error)
}

// WriteSub defines a common way to write a record to a builder but return its
// new position.
type WriteSub interface {
	WriteSub(fb *flatbuffers.Builder) (flatbuffers.UOffsetT, error)
}

// Read defines a common way to read a slice of bytes on into a record
type Read interface {
	Read([]byte) error
}

// Id defines a struct that represents a ObjectId
type Id struct {
	Hex string
}

// MakeId creates an id from a hex string
func MakeId(hex string) Id {
	return Id{hex}
}

func (i Id) WriteSub(fb *flatbuffers.Builder) (flatbuffers.UOffsetT, error) {
	if !bson.IsObjectIdHex(i.Hex) {
		return 0, ErrInvalidIdHex(3)
	}

	position := fb.CreateString(i.Hex)

	schema.IdStart(fb)
	schema.IdAddHex(fb, position)

	return schema.IdEnd(fb), nil
}

// Version defines a struct that represents a Version
type Version struct {
	Version string
}

func (c Version) WriteSub(fb *flatbuffers.Builder) (flatbuffers.UOffsetT, error) {
	position := fb.CreateString(c.Version)

	schema.VersionStart(fb)
	schema.VersionAddVersion(fb, position)

	return schema.VersionEnd(fb), nil
}

// Cost defines a struct that represents a currency and price as a tuple
type Cost struct {
	Currency string
	Price    uint64
}

func (c Cost) WriteSub(fb *flatbuffers.Builder) (flatbuffers.UOffsetT, error) {
	if len(c.Currency) < 1 {
		return 0, ErrInvalidLength(15)
	}

	position := fb.CreateString(c.Currency)

	schema.CostStart(fb)
	schema.CostAddCurrency(fb, position)
	schema.CostAddPrice(fb, c.Price)

	return schema.CostEnd(fb), nil
}

type Dates struct {
	Start, End uint64
}

func (r Dates) WriteSub(fb *flatbuffers.Builder) (flatbuffers.UOffsetT, error) {
	schema.DatesStart(fb)
	schema.DatesAddStart(fb, r.Start)
	schema.DatesAddEnd(fb, r.End)

	return schema.DatesEnd(fb), nil
}

type Codes struct {
	BarcodeType, BarcodeOrigin, BarcodeSource, QRCode string
}

func (c Codes) WriteSub(fb *flatbuffers.Builder) (flatbuffers.UOffsetT, error) {
	if len(c.BarcodeType) < 1 || len(c.BarcodeOrigin) < 1 ||
		len(c.BarcodeSource) < 1 || len(c.QRCode) < 1 {
		return 0, ErrInvalidLength(16)
	}

	var (
		position0 = fb.CreateString(c.BarcodeType)
		position1 = fb.CreateString(c.BarcodeOrigin)
		position2 = fb.CreateString(c.BarcodeSource)
		position3 = fb.CreateString(c.QRCode)
	)

	schema.CodesStart(fb)
	schema.CodesAddBarcodeType(fb, position0)
	schema.CodesAddBarcodeOrigin(fb, position1)
	schema.CodesAddBarcodeSource(fb, position2)
	schema.CodesAddQrcode(fb, position3)

	return schema.CodesEnd(fb), nil
}

type KeyFieldScoreTxnValues []KeyFieldScoreTxnValue

func (k KeyFieldScoreTxnValues) KeyFieldValues() []selectors.KeyFieldTxnValue {
	res := make([]selectors.KeyFieldTxnValue, 0, len(k))
	for _, v := range k {
		res = append(res, v.KeyFieldTxnValue())
	}
	return res
}

type KeyFieldScoreTxnValue struct {
	Key, Field bs.Key
	Score      float64
	Txn        bs.Key
	Value      string
}

func (k KeyFieldScoreTxnValue) WriteSub(fb *flatbuffers.Builder) (flatbuffers.UOffsetT, error) {
	var (
		key         = k.Key.String()
		field       = k.Field.String()
		transaction = k.Txn.String()
	)

	if !bson.IsObjectIdHex(key) || !bson.IsObjectIdHex(field) || !bson.IsObjectIdHex(transaction) {
		return 0, ErrInvalidIdHex(4)
	}

	if len(k.Value) < 1 {
		return 0, ErrInvalidLength(17)
	}

	var (
		position0 = fb.CreateString(key)
		position1 = fb.CreateString(field)
		position2 = fb.CreateString(transaction)
		position3 = fb.CreateString(k.Value)
	)

	schema.KeyFieldScoreTxnValueStart(fb)
	schema.KeyFieldScoreTxnValueAddKey(fb, position0)
	schema.KeyFieldScoreTxnValueAddField(fb, position1)
	schema.KeyFieldScoreTxnValueAddScore(fb, k.Score)
	schema.KeyFieldScoreTxnValueAddTxn(fb, position2)
	schema.KeyFieldScoreTxnValueAddValue(fb, position3)

	return schema.KeyFieldScoreTxnValueEnd(fb), nil
}

func (k KeyFieldScoreTxnValue) KeyFieldTxnValue() selectors.KeyFieldTxnValue {
	return selectors.KeyFieldTxnValue{
		Key:   k.Key,
		Field: k.Field,
		Txn:   k.Txn,
		Value: k.Value,
	}
}

func FromKeyFieldScoreTxnValues(values []selectors.KeyFieldScoreTxnValue) []KeyFieldScoreTxnValue {
	res := make([]KeyFieldScoreTxnValue, 0, len(values))
	for _, v := range values {
		res = append(res, KeyFieldScoreTxnValue{
			Key:   v.Key,
			Field: v.Field,
			Score: v.Score,
			Txn:   v.Txn,
			Value: v.Value,
		})
	}
	return res
}

func FromKeyFieldScoreTxnValue(value selectors.KeyFieldScoreTxnValue) KeyFieldScoreTxnValue {
	return KeyFieldScoreTxnValue{
		Key:   value.Key,
		Field: value.Field,
		Score: value.Score,
		Txn:   value.Txn,
		Value: value.Value,
	}
}

type Header struct {
	Type           int8
	Field, OwnerId bs.Key
	Updated        int64
}

func (h *Header) Read(bytes []byte) error {
	if len(bytes) < 1 {
		return ErrInvalidLength(18)
	}

	var (
		record = schema.GetRootAsHeader(bytes, 0)

		field   = string(record.Id(nil).Hex())
		ownerId = string(record.OwnerId(nil).Hex())
	)

	if !bson.IsObjectIdHex(field) || !bson.IsObjectIdHex(ownerId) {
		return ErrInvalidIdHex(5)
	}

	h.Type = record.Typ()
	h.Field = bs.Key(field)
	h.OwnerId = bs.Key(ownerId)
	h.Updated = int64(record.Updated())

	return nil
}

// ReadType returns the type of value stored in the value
func ReadType(value string) (int, error) {
	x, err := strconv.ParseInt(value[0:1], 10, 32)
	if err != nil {
		return 0, err
	}
	return int(x), nil
}

// ReadBody returns the body of the value, minus the header
func ReadBody(value string) ([]byte, error) {
	if len(value) < 2 {
		return nil, ErrInvalidLength(19)
	}
	return []byte(value[1:]), nil
}

type KeyField struct {
	Key, Field bs.Key
}

func (k KeyField) WriteSub(fb *flatbuffers.Builder) (flatbuffers.UOffsetT, error) {
	var (
		key   = k.Key.String()
		field = k.Field.String()
	)

	if !bson.IsObjectIdHex(key) || !bson.IsObjectIdHex(field) {
		return 0, ErrInvalidIdHex(6)
	}

	var (
		position0 = fb.CreateString(key)
		position1 = fb.CreateString(field)
	)

	schema.KeyFieldStart(fb)
	schema.KeyFieldAddKey(fb, position0)
	schema.KeyFieldAddField(fb, position1)

	return schema.KeyFieldEnd(fb), nil
}

type KeyFieldSizeExpiry struct {
	Key, Field bs.Key
	Size       int64
	Expiry     time.Duration
}

func (k KeyFieldSizeExpiry) Write(fb *flatbuffers.Builder) ([]byte, error) {
	position, err := k.WriteSub(fb)
	if err != nil {
		return nil, err
	}

	fb.Finish(position)
	return fb.FinishedBytes(), nil
}

func (k KeyFieldSizeExpiry) WriteSub(fb *flatbuffers.Builder) (flatbuffers.UOffsetT, error) {
	var (
		key   = k.Key.String()
		field = k.Field.String()
	)

	if !bson.IsObjectIdHex(key) || !bson.IsObjectIdHex(field) {
		return 0, ErrInvalidIdHex(7)
	}

	var (
		position0 = fb.CreateString(key)
		position1 = fb.CreateString(field)
	)

	schema.KeyFieldSizeExpiryStart(fb)
	schema.KeyFieldSizeExpiryAddKey(fb, position0)
	schema.KeyFieldSizeExpiryAddField(fb, position1)
	schema.KeyFieldSizeExpiryAddSize(fb, k.Size)
	schema.KeyFieldSizeExpiryAddExpiry(fb, k.Expiry.Nanoseconds())

	return schema.KeyFieldSizeExpiryEnd(fb), nil
}

func (k *KeyFieldSizeExpiry) Read(bytes []byte) error {
	if len(bytes) < 1 {
		return ErrInvalidLength(20)
	}

	var (
		record = getRootAsKeyFieldSizeExpiry(bytes, 0)

		key   = string(record.Key())
		field = string(record.Field())
	)

	if !bson.IsObjectIdHex(key) || !bson.IsObjectIdHex(field) {
		return ErrInvalidIdHex(8)
	}

	k.Key = bs.Key(key)
	k.Field = bs.Key(field)
	k.Size = record.Size()
	k.Expiry = time.Duration(record.Expiry())

	return nil
}

// polyfill until flatbuffers library generate this for all instances.
func getRootAsKeyFieldSizeExpiry(buf []byte, offset flatbuffers.UOffsetT) *schema.KeyFieldSizeExpiry {
	var (
		n = flatbuffers.GetUOffsetT(buf[offset:])
		x = &schema.KeyFieldSizeExpiry{}
	)
	x.Init(buf, n+offset)
	return x
}

type KeyFieldScoreSizeExpiry struct {
	Key, Field bs.Key
	Score      float64
	Size       int64
	Expiry     time.Duration
}

func (k KeyFieldScoreSizeExpiry) Write(fb *flatbuffers.Builder) ([]byte, error) {
	position, err := k.WriteSub(fb)
	if err != nil {
		return nil, err
	}

	fb.Finish(position)
	return fb.FinishedBytes(), nil
}

func (k KeyFieldScoreSizeExpiry) WriteSub(fb *flatbuffers.Builder) (flatbuffers.UOffsetT, error) {
	var (
		key   = k.Key.String()
		field = k.Field.String()
	)

	if !bson.IsObjectIdHex(key) || !bson.IsObjectIdHex(field) {
		return 0, ErrInvalidIdHex(9)
	}

	var (
		position0 = fb.CreateString(key)
		position1 = fb.CreateString(field)
	)

	schema.KeyFieldScoreSizeExpiryStart(fb)
	schema.KeyFieldScoreSizeExpiryAddKey(fb, position0)
	schema.KeyFieldScoreSizeExpiryAddField(fb, position1)
	schema.KeyFieldScoreSizeExpiryAddScore(fb, k.Score)
	schema.KeyFieldScoreSizeExpiryAddSize(fb, k.Size)
	schema.KeyFieldScoreSizeExpiryAddExpiry(fb, k.Expiry.Nanoseconds())

	return schema.KeyFieldScoreSizeExpiryEnd(fb), nil
}

func (k *KeyFieldScoreSizeExpiry) Read(bytes []byte) error {
	if len(bytes) < 1 {
		return ErrInvalidLength(21)
	}

	var (
		record = schema.GetRootAsKeyFieldScoreSizeExpiry(bytes, 0)

		key   = string(record.Key())
		field = string(record.Field())
	)

	if !bson.IsObjectIdHex(key) || !bson.IsObjectIdHex(field) {
		return ErrInvalidIdHex(10)
	}

	k.Key = bs.Key(key)
	k.Field = bs.Key(field)
	k.Score = record.Score()
	k.Size = record.Size()
	k.Expiry = time.Duration(record.Expiry())

	return nil
}

type QueryRecord struct {
	Key, Field bs.Key
	Record     string
}

func (k QueryRecord) WriteSub(fb *flatbuffers.Builder) (flatbuffers.UOffsetT, error) {
	var (
		key   = k.Key.String()
		field = k.Field.String()
	)

	if !bson.IsObjectIdHex(key) || !bson.IsObjectIdHex(field) {
		return 0, ErrInvalidIdHex(11)
	}

	if len(k.Record) < 1 {
		return 0, ErrInvalidLength(22)
	}

	var (
		position0 = fb.CreateString(key)
		position1 = fb.CreateString(field)
		position2 = fb.CreateString(k.Record)
	)

	schema.QueryRecordStart(fb)
	schema.QueryRecordAddKey(fb, position0)
	schema.QueryRecordAddField(fb, position1)
	schema.QueryRecordAddRecord(fb, position2)

	return schema.QueryRecordEnd(fb), nil
}

func FromQueryRecords(values []selectors.QueryRecord) ([]QueryRecord, error) {
	res := make([]QueryRecord, 0, len(values))
	for _, v := range values {

		data, err := json.Marshal(v.Record)
		if err != nil {
			return nil, err
		}

		res = append(res, QueryRecord{
			Key:    v.Key,
			Field:  v.Field,
			Record: string(data),
		})
	}
	return res, nil
}
