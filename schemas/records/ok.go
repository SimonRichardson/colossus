package records

import (
	"time"

	bs "github.com/SimonRichardson/colossus/blist/selectors"
	"github.com/SimonRichardson/colossus/schemas/schema"
	flatbuffers "github.com/google/flatbuffers/go"
)

type OKVersion struct {
	Duration time.Duration
	Records  Version
}

func (o OKVersion) Write(fb *flatbuffers.Builder) ([]byte, error) {
	position, err := o.Records.WriteSub(fb)
	if err != nil {
		return nil, err
	}

	schema.OKVersionStart(fb)
	schema.OKVersionAddDuration(fb, int64(o.Duration))
	schema.OKVersionAddRecords(fb, position)

	position = schema.OKVersionEnd(fb)

	fb.Finish(position)
	return fb.FinishedBytes(), nil
}

func (o *OKVersion) Read(bytes []byte) error {
	if len(bytes) < 0 {
		return ErrInvalidLength(4)
	}

	record := schema.GetRootAsOKVersion(bytes, 0)

	o.Duration = time.Duration(record.Duration())

	version := record.Records(nil).Version()
	if len(version) < 0 {
		return ErrInvalidLength(5)
	}
	o.Records = Version{string(version)}

	return nil
}

type OKInt struct {
	Duration time.Duration
	Records  int
}

func (o OKInt) Write(fb *flatbuffers.Builder) ([]byte, error) {
	schema.OKIntStart(fb)
	schema.OKIntAddDuration(fb, int64(o.Duration))
	schema.OKIntAddRecords(fb, int64(o.Records))

	position := schema.OKIntEnd(fb)

	fb.Finish(position)
	return fb.FinishedBytes(), nil
}

func (o *OKInt) Read(bytes []byte) error {
	if len(bytes) < 0 {
		return ErrInvalidLength(6)
	}

	record := schema.GetRootAsOKInt(bytes, 0)

	o.Duration = time.Duration(record.Duration())
	o.Records = int(record.Records())

	return nil
}

type OKNoContent struct {
	Duration time.Duration
}

func (o OKNoContent) Write(fb *flatbuffers.Builder) ([]byte, error) {
	schema.OKNoContentStart(fb)
	schema.OKNoContentAddDuration(fb, int64(o.Duration))

	position := schema.OKNoContentEnd(fb)

	fb.Finish(position)
	return fb.FinishedBytes(), nil
}

func (o *OKNoContent) Read(bytes []byte) error {
	if len(bytes) < 0 {
		return ErrInvalidLength(6)
	}

	record := schema.GetRootAsOKNoContent(bytes, 0)

	o.Duration = time.Duration(record.Duration())

	return nil
}

type OKKeyFieldScoreTxnValues struct {
	Duration time.Duration
	Records  []KeyFieldScoreTxnValue
}

func (o OKKeyFieldScoreTxnValues) Write(fb *flatbuffers.Builder) ([]byte, error) {
	var (
		num       = len(o.Records)
		positions = make([]flatbuffers.UOffsetT, num, num)
	)

	for k, v := range o.Records {
		position, err := v.WriteSub(fb)
		if err != nil {
			return nil, err
		}
		positions[num-1-k] = position
	}

	schema.OKKeyFieldScoreTxnValuesStartRecordsVector(fb, num)

	for _, v := range positions {
		fb.PrependUOffsetT(v)
	}

	vector := fb.EndVector(num)

	schema.OKKeyFieldScoreTxnValuesStart(fb)
	schema.OKKeyFieldScoreTxnValuesAddDuration(fb, int64(o.Duration))
	schema.OKKeyFieldScoreTxnValuesAddRecords(fb, vector)

	position := schema.OKKeyFieldScoreTxnValuesEnd(fb)

	fb.Finish(position)
	return fb.FinishedBytes(), nil
}

func (o *OKKeyFieldScoreTxnValues) Read(bytes []byte) error {
	record := schema.GetRootAsOKKeyFieldScoreTxnValues(bytes, 0)

	o.Duration = time.Duration(record.Duration())

	var (
		num    = record.RecordsLength()
		vector = make([]KeyFieldScoreTxnValue, 0, num)
	)

	for i := 0; i < num; i++ {
		var k schema.KeyFieldScoreTxnValue
		if !record.Records(&k, i) {
			return ErrInvalidRecord(0)
		}

		vector = append(vector, KeyFieldScoreTxnValue{
			Key:   bs.Key(string(k.Key())),
			Field: bs.Key(string(k.Field())),
			Score: k.Score(),
			Txn:   bs.Key(string(k.Txn())),
			Value: string(k.Value()),
		})
	}

	o.Records = vector

	return nil
}

type OKKeyFieldScoreTxnValue struct {
	Duration time.Duration
	Records  KeyFieldScoreTxnValue
}

func (o OKKeyFieldScoreTxnValue) Write(fb *flatbuffers.Builder) ([]byte, error) {
	valuePosition, err := o.Records.WriteSub(fb)
	if err != nil {
		return nil, err
	}

	schema.OKKeyFieldScoreTxnValuesStart(fb)
	schema.OKKeyFieldScoreTxnValuesAddDuration(fb, int64(o.Duration))
	schema.OKKeyFieldScoreTxnValuesAddRecords(fb, valuePosition)

	position := schema.OKKeyFieldScoreTxnValuesEnd(fb)

	fb.Finish(position)
	return fb.FinishedBytes(), nil
}

func (o *OKKeyFieldScoreTxnValue) Read(bytes []byte) error {
	var (
		record = schema.GetRootAsOKKeyFieldScoreTxnValue(bytes, 0)
		result = record.Records(nil)
	)

	o.Duration = time.Duration(record.Duration())
	o.Records = KeyFieldScoreTxnValue{
		Key:   bs.Key(string(result.Key())),
		Field: bs.Key(string(result.Field())),
		Score: result.Score(),
		Txn:   bs.Key(string(result.Txn())),
		Value: string(result.Value()),
	}

	return nil
}

type OKQuery struct {
	Duration time.Duration
	Records  []QueryRecord
}

func (o OKQuery) Write(fb *flatbuffers.Builder) ([]byte, error) {
	var (
		num       = len(o.Records)
		positions = make([]flatbuffers.UOffsetT, num, num)
	)

	for k, v := range o.Records {
		position, err := v.WriteSub(fb)
		if err != nil {
			return nil, err
		}
		positions[num-1-k] = position
	}

	schema.OKQueryStartRecordsVector(fb, num)

	for _, v := range positions {
		fb.PrependUOffsetT(v)
	}

	vector := fb.EndVector(num)

	schema.OKQueryStart(fb)
	schema.OKQueryAddDuration(fb, int64(o.Duration))
	schema.OKQueryAddRecords(fb, vector)

	position := schema.OKQueryEnd(fb)

	fb.Finish(position)
	return fb.FinishedBytes(), nil
}

func (o *OKQuery) Read(bytes []byte) error {
	record := schema.GetRootAsOKQuery(bytes, 0)

	o.Duration = time.Duration(record.Duration())

	var (
		num    = record.RecordsLength()
		vector = make([]QueryRecord, 0, num)
	)

	for i := 0; i < num; i++ {
		var k schema.QueryRecord
		if !record.Records(&k, i) {
			return ErrInvalidRecord(1)
		}

		vector = append(vector, QueryRecord{
			Key:    bs.Key(string(k.Key())),
			Field:  bs.Key(string(k.Field())),
			Record: string(k.Record()),
		})
	}

	o.Records = vector

	return nil
}
