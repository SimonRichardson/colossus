package lorenz

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SimonRichardson/colossus/blist/selectors"
	cli "github.com/SimonRichardson/colossus/blist/services/lorenz/client"
)

type codeBody struct {
	Records lorenzCodeSet `json:"records"`
}

func readCode(client cli.Client, event selectors.Event, user selectors.User) (selectors.CodeSet, error) {
	response, err := client.Get(
		fmt.Sprintf("/events/%s/code?type=%s", event.Id, event.CodeTypes.BarcodeType.String()),
		cli.Versioned,
		func(headers http.Header) {
			headers.Set("Authorization", fmt.Sprintf("Bearer %s", user.Access.Token))
		},
	)

	if err != nil {
		return selectors.CodeSet{}, err
	}

	if err := readError(response, http.StatusOK); err != nil {
		return selectors.CodeSet{}, err
	}

	var record codeBody
	if err := json.Unmarshal(response.Bytes, &record); err != nil {
		return selectors.CodeSet{}, err
	}

	codeSet := record.Records
	return selectors.CodeSet{
		Barcode: selectors.Barcode{
			Type:   selectors.BarcodeType(codeSet.Barcode.Type),
			Origin: codeSet.Barcode.Origin,
			Source: codeSet.Barcode.Source,
		},
		QRCode: selectors.QRCode{
			Type:   selectors.QRCodeType(codeSet.QRCode.Type),
			Source: codeSet.QRCode.Source,
		},
	}, nil
}

type lorenzCodeSet struct {
	Barcode lorenzBarcode `json:"barcode"`
	QRCode  lorenzQRCode  `json:"qrcode"`
}

type lorenzBarcode struct {
	Type   string `json:"type"`
	Origin string `json:"origin"`
	Source string `json:"source"`
}

type lorenzQRCode struct {
	Type   string `json:"type"`
	Source string `json:"source"`
}
