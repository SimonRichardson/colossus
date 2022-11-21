package records

import (
	"github.com/SimonRichardson/colossus/schemas/schema"
	flatbuffers "github.com/google/flatbuffers/go"
)

type Error struct {
	Error       string
	Code        string
	Description string
}

func (e Error) Write(fb *flatbuffers.Builder) ([]byte, error) {
	var (
		errorPosition       = fb.CreateString(e.Error)
		codePosition        = fb.CreateString(e.Code)
		descriptionPosition = fb.CreateString(e.Description)
	)

	schema.ErrorStart(fb)
	schema.ErrorAddError(fb, errorPosition)
	schema.ErrorAddCode(fb, codePosition)
	schema.ErrorAddDescription(fb, descriptionPosition)

	position := schema.ErrorEnd(fb)

	fb.Finish(position)
	return fb.FinishedBytes(), nil
}

func (e *Error) Read(bytes []byte) error {
	if len(bytes) < 0 {
		return ErrInvalidLength(3)
	}

	record := schema.GetRootAsError(bytes, 0)

	e.Error = string(record.Error())
	e.Code = string(record.Code())
	e.Description = string(record.Description())

	return nil
}
