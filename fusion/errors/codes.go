package errors

import "github.com/SimonRichardson/colossus/typex"

const (
	Source typex.ErrorSource = "Fusion"
)

var (
	InvalidContentType = typex.BadRequest.With("Invalid Content Type")
	InvalidArgument    = typex.BadRequest.With("Invalid Argument")

	Fatal                   = typex.InternalServerError.With("Fatal")
	Complete                = typex.InternalServerError.With("Complete")
	Partial                 = typex.InternalServerError.With("Partial")
	Repair                  = typex.InternalServerError.With("Repair")
	UnexpectedArgument      = typex.InternalServerError.With("Unexpected Argument")
	UnexpectedParseArgument = typex.InternalServerError.With("Unexpected Parse Argument")
	UnexpectedResults       = typex.InternalServerError.With("Unexpected Results")
	NativeError             = typex.InternalServerError.With("Native Error")
	NoCaseFound             = typex.InternalServerError.With("No Case Found")
	ExpiredNode             = typex.InternalServerError.With("Expired Node")
	RateLimited             = typex.InternalServerError.With("Rate Limited")
	MaxSize                 = typex.InternalServerError.With("Max Size")

	MissingContent = typex.NotFound.With("Missing Content")
)
