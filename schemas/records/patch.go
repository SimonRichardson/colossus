package records

import (
	"time"

	"github.com/SimonRichardson/colossus/schemas/schema"
	flatbuffers "github.com/google/flatbuffers/go"
)

type PatchRecords struct {
	Operations []Operation
	Score      float64
	MaxSize    int64
	Expiry     time.Duration
}

func (r PatchRecords) Write(fb *flatbuffers.Builder) ([]byte, error) {
	var (
		num       = len(r.Operations)
		positions = make([]flatbuffers.UOffsetT, num, num)
	)

	for k, v := range r.Operations {
		position, err := v.WriteSub(fb)
		if err != nil {
			return nil, err
		}
		positions[num-1-k] = position
	}

	schema.PatchRequestStartOperationsVector(fb, num)

	for _, v := range positions {
		fb.PrependUOffsetT(v)
	}

	vector := fb.EndVector(num)

	schema.PatchRequestStart(fb)
	schema.PatchRequestAddScore(fb, r.Score)
	schema.PatchRequestAddOperations(fb, vector)
	schema.PatchRequestAddMaxSize(fb, uint64(r.MaxSize))
	schema.PatchRequestAddExpiry(fb, uint64(r.Expiry.Nanoseconds()))
	position := schema.PatchRequestEnd(fb)

	fb.Finish(position)
	return fb.FinishedBytes(), nil
}

type Operation struct {
	Op    string
	Path  string
	Value string
}

func (r Operation) WriteSub(fb *flatbuffers.Builder) (flatbuffers.UOffsetT, error) {
	var (
		opPosition    = fb.CreateString(r.Op)
		pathPosition  = fb.CreateString(r.Path)
		valuePosition = fb.CreateString(r.Value)
	)

	schema.OperationStart(fb)
	schema.OperationAddOp(fb, opPosition)
	schema.OperationAddPath(fb, pathPosition)
	schema.OperationAddValue(fb, valuePosition)

	return schema.OperationEnd(fb), nil
}
