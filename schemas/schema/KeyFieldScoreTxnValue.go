// automatically generated by the FlatBuffers compiler, do not modify

package schema

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type KeyFieldScoreTxnValue struct {
	_tab flatbuffers.Table
}

func GetRootAsKeyFieldScoreTxnValue(buf []byte, offset flatbuffers.UOffsetT) *KeyFieldScoreTxnValue {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &KeyFieldScoreTxnValue{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *KeyFieldScoreTxnValue) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *KeyFieldScoreTxnValue) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *KeyFieldScoreTxnValue) Key() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *KeyFieldScoreTxnValue) Field() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *KeyFieldScoreTxnValue) Score() float64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetFloat64(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *KeyFieldScoreTxnValue) MutateScore(n float64) bool {
	return rcv._tab.MutateFloat64Slot(8, n)
}

func (rcv *KeyFieldScoreTxnValue) Txn() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *KeyFieldScoreTxnValue) Value() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func KeyFieldScoreTxnValueStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func KeyFieldScoreTxnValueAddKey(builder *flatbuffers.Builder, key flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(key), 0)
}
func KeyFieldScoreTxnValueAddField(builder *flatbuffers.Builder, field flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(field), 0)
}
func KeyFieldScoreTxnValueAddScore(builder *flatbuffers.Builder, score float64) {
	builder.PrependFloat64Slot(2, score, 0.0)
}
func KeyFieldScoreTxnValueAddTxn(builder *flatbuffers.Builder, txn flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(txn), 0)
}
func KeyFieldScoreTxnValueAddValue(builder *flatbuffers.Builder, value flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(value), 0)
}
func KeyFieldScoreTxnValueEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
