package lorenz

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SimonRichardson/colossus/blist/common"
	"github.com/SimonRichardson/colossus/blist/selectors"
	cli "github.com/SimonRichardson/colossus/blist/services/lorenz/client"
)

type transaction struct {
	Transaction string `json:"transaction"`
}

type transactionBody struct {
	Records transaction `json:"records"`
}

func charge(client cli.Client,
	event selectors.Event,
	user selectors.User,
	payment selectors.Payment,
) (selectors.KeyTxn, error) {
	cost := event.Tickets.Cost
	bytes, err := json.Marshal(lorenzPayment{
		Transaction: common.StringPtr(payment.Key.String()),
		Cost: &lorenzCost{
			Amount:   cost.Amount,
			Currency: cost.Currency,
			Fees:     getChargeFees(cost.Fees),
		},
		Count: common.IntPtr(payment.Count),
		Method: &lorenzMethod{
			Type:  lorenzMethodType(payment.Method.Type.String()),
			Token: common.StringPtr(payment.Method.Token.String()),
		},
		UserInfo: &lorenzUserInfo{
			FullName:   common.StringPtr(payment.UserInfo.FullName),
			PostalCode: common.StringPtr(payment.UserInfo.PostalCode),
		},
	})
	if err != nil {
		return selectors.KeyTxn{Key: event.Id}, err
	}

	response, err := client.Post(fmt.Sprintf("/events/%s/charge", event.Id.String()),
		bytes,
		cli.Versioned,
		func(headers http.Header) {
			headers.Set("Authorization", fmt.Sprintf("Bearer %s", user.Access.Token))
		},
	)
	if err != nil {
		return selectors.KeyTxn{Key: event.Id}, err
	}

	if err := readError(response, http.StatusOK); err != nil {
		return selectors.KeyTxn{Key: event.Id}, err
	}

	var record transactionBody
	if err := json.Unmarshal(response.Bytes, &record); err != nil {
		return selectors.KeyTxn{Key: event.Id}, err
	}

	return selectors.KeyTxn{
		Key: event.Id,
		Txn: selectors.Key(record.Records.Transaction),
	}, nil
}

func getChargeFees(fees map[selectors.FeeType]selectors.Fee) map[string]lorenzFee {
	res := make(map[string]lorenzFee)
	for k, v := range fees {
		res[k.String()] = lorenzFee{
			Type:    v.Type.String(),
			Percent: fmt.Sprintf("%f", v.Percent),
			Fixed:   v.Fixed,
		}
	}

	return res
}

// Payment describes the information required for making a payment for the
// items.
type lorenzPayment struct {
	Transaction *string         `json:"transaction"`
	Cost        *lorenzCost     `json:"cost"`
	Count       *int            `json:"count"`
	Method      *lorenzMethod   `json:"method"`
	UserInfo    *lorenzUserInfo `json:"user_info"`
}

// Cost type describes an amount for paying for a transaction. Note that
// currency can not be empty, even if the amount is free.
type lorenzCost struct {
	Currency string               `json:"currency"`
	Amount   uint64               `json:"amount"`
	Fees     map[string]lorenzFee `json:"fees"`
}

type lorenzFee struct {
	Type    string `json:"type"`
	Percent string `json:"percent"`
	Fixed   uint64 `json:"fixed"`
}

// Method type encapsulates a method type and token required to authorize the
// transaction.
type lorenzMethod struct {
	Type  lorenzMethodType `json:"type"`
	Token *string          `json:"token,omitempty"`
}

// MethodType defines a new type for different payment types.
type lorenzMethodType string

// UserInfo defines tailored user information for used with purchasing
type lorenzUserInfo struct {
	FullName   *string `json:"full_name"`
	PostalCode *string `json:"postal_code,omitempty"`
}
