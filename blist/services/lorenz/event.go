package lorenz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/SimonRichardson/colossus/blist/common"
	"github.com/SimonRichardson/colossus/blist/errors"
	"github.com/SimonRichardson/colossus/blist/selectors"
	cli "github.com/SimonRichardson/colossus/blist/services/lorenz/client"
	"github.com/SimonRichardson/colossus/typex"
)

var (
	defaultBehaviourTime = time.Date(2017, 12, 1, 1, 1, 1, 1, time.UTC)
)

type eventsBody struct {
	Records []lorenzEvent `json:"records"`
}

type eventBody struct {
	Records lorenzEvent `json:"records"`
}

func readEvents(client cli.Client, offset, limit int) ([]selectors.Event, error) {
	var (
		uri           = fmt.Sprintf("/events%s", getOffset(offset, limit))
		response, err = client.Get(uri, cli.Versioned, func(http.Header) {})
	)
	if err != nil {
		return nil, err
	}

	if err := readError(response, http.StatusOK); err != nil {
		return nil, err
	}

	var record eventsBody
	if err := json.Unmarshal(response.Bytes, &record); err != nil {
		return nil, err
	}

	res := make([]selectors.Event, 0, len(record.Records))
	for _, v := range record.Records {
		res = append(res, unselector(v))
	}
	return res, nil
}

func getOffset(offset, limit int) string {
	var (
		off, lim string
	)
	if offset >= 0 {
		off = fmt.Sprintf("offset=%d", offset)
	}
	if limit > 0 {
		lim = fmt.Sprintf("limit=%d", limit)
	}

	if off != "" && lim != "" {
		return fmt.Sprintf("?%s&%s", off, lim)
	}
	if off != "" {
		return fmt.Sprintf("?%s", off)
	}
	if lim != "" {
		return fmt.Sprintf("?%s", lim)
	}
	return ""
}

func readEvent(client cli.Client, key selectors.Key) (selectors.Event, error) {
	response, err := client.Get(fmt.Sprintf("/events/%s", key.String()),
		cli.Versioned,
		func(http.Header) {},
	)
	if err != nil {
		return selectors.Event{}, err
	}

	if err := readError(response, http.StatusOK); err != nil {
		return selectors.Event{}, err
	}

	var record eventBody
	if err := json.Unmarshal(response.Bytes, &record); err != nil {
		return selectors.Event{}, err
	}

	event := record.Records
	return unselector(event), nil
}

func unselector(event lorenzEvent) selectors.Event {
	return selectors.Event{
		Behaviours:  getSelectorBehaviours(event),
		Id:          selectors.Key(event.Id.Hex()),
		Name:        common.StringUnptr(event.Name),
		PermName:    common.StringUnptr(event.PermName),
		Description: common.StringUnptr(event.Description),
		Place: selectors.Place{
			Venue:    common.StringUnptr(event.Venue),
			Address:  common.StringUnptr(event.Address),
			Location: getLocation(event),
		},
		Tickets: selectors.Tickets{
			Cost: selectors.Cost{
				Currency: string(event.Cost.Currency),
				Amount:   event.Cost.Amount,
				Fees:     getSelectorFees(event.Cost),
			},
			TotalNum:       common.Uint64Unptr(event.TotalNumTickets),
			ExpiryDuration: common.DurationUnptr(event.TicketExpiryDuration),
		},
		WaitingLists: selectors.WaitingLists{
			Cost: selectors.Cost{
				Currency: string(event.Cost.Currency), // Make sure we never have to differing currencies.
				Amount:   event.WaitingLists.Cost.Amount,
				Fees:     getSelectorFees(&event.WaitingLists.Cost),
			},
		},
		State: selectors.State{
			CodeLocked: contains(common.StringsUnptr(event.Flags), "code-locked"),
		},
		Dates: selectors.Dates{
			Start: common.TimeUnptr(event.Date),
			End:   common.TimeUnptr(event.DateEnd),
		},
		PaymentAccounts: selectors.PaymentAccounts{},
		CodeTypes: selectors.CodeTypes{
			BarcodeType: selectors.BarcodeType(common.StringUnptr(event.BarcodeType)),
		},
	}
}

func getSelectorBehaviours(event lorenzEvent) selectors.Behaviours {
	if event.Behaviours != nil {
		return selectors.Behaviours{
			CreatedAt: event.Behaviours.CreatedAt,
			UpdatedAt: event.Behaviours.UpdatedAt,
		}
	}
	return selectors.Behaviours{
		CreatedAt: defaultBehaviourTime,
		UpdatedAt: defaultBehaviourTime,
	}
}

func getSelectorFees(cost *lorenzEventCost) map[selectors.FeeType]selectors.Fee {
	res := make(map[selectors.FeeType]selectors.Fee)

	for k, v := range cost.Fees {
		percent, err := strconv.ParseFloat(v.Percent, 64)
		if err != nil {
			typex.Errorf(errors.Source, errors.InvalidArgument,
				"Unable to parse fees as percent").With(err)
			continue
		}
		res[selectors.FeeType(k)] = selectors.Fee{
			Type:    selectors.FeeAmountType(v.Type),
			Percent: percent,
			Fixed:   v.Fixed,
		}
	}

	return res
}

func getLocation(event lorenzEvent) selectors.Location {
	if event.Location != nil {
		return selectors.Location{
			Lat:      event.Location.Lat,
			Lng:      event.Location.Lng,
			Place:    event.Location.Place,
			Accuracy: event.Location.Accuracy,
		}
	}

	return selectors.Location{}
}

type lorenzEvent struct {
	Behaviours                *lorenzBehaviours                `json:"behaviours"`
	Id                        *bson.ObjectId                   `json:"id,omitempty"`
	Date                      *time.Time                       `json:"date"`
	DateEnd                   *time.Time                       `json:"date_end"`
	Name                      *string                          `json:"name"`
	Cost                      *lorenzEventCost                 `json:"cost"`
	Description               *string                          `json:"description"`
	BackgroundImage           *string                          `json:"background_image"`
	TicketImage               *string                          `json:"ticket_image"`
	Media                     *[]lorenzMedia                   `json:"media"`
	Colour                    *lorenzColour                    `json:"colour"`
	Type                      *string                          `json:"type"`
	Venue                     *string                          `json:"venue"`
	Address                   *string                          `json:"address"`
	Location                  *lorenzEventLocation             `json:"location"`
	SoldOut                   *bool                            `json:"sold_out"`
	Hidden                    *bool                            `json:"hidden"`
	Deleted                   *bool                            `json:"deleted"`
	PermName                  *string                          `json:"perm_name"`
	Flags                     *[]string                        `json:"flags"`
	Tags                      *[]string                        `json:"tags"`
	InheritedTags             *[]string                        `json:"inherited_tags"`
	BarcodeType               *string                          `json:"barcode_type"`
	MaxTicketsLimit           *uint64                          `json:"max_tickets_limit"`
	TotalNumTickets           *uint64                          `json:"total_num_tickets"`
	TicketExpiryDuration      *time.Duration                   `json:"ticket_expiry_duration"`
	WaitingListExpiryDuration *time.Duration                   `json:"waitinglist_expiry_duration"`
	CityIds                   *[]bson.ObjectId                 `json:"related_city_ids"`
	ArtistsIds                *[]bson.ObjectId                 `json:"related_artists_ids"`
	EventIds                  *[]bson.ObjectId                 `json:"related_event_ids"`
	TicketTransfer            *lorenzTicketMode                `json:"transfer"`
	TicketReturns             *lorenzTicketMode                `json:"returns"`
	PaymentAccounts           *map[string]lorenzPaymentAccount `json:"payments"`
	WaitingLists              *lorenzWaitingLists              `json:"waiting_lists"`
}

type lorenzBehaviours struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type lorenzCurrency string

type lorenzEventCost struct {
	Currency lorenzCurrency                `json:"currency"`
	Amount   uint64                        `json:"amount"`
	Fees     map[string]lorenzEventCostFee `json:"fees"`
}

type lorenzEventCostFee struct {
	Type    string `json:"type"`
	Fixed   uint64 `json:"fixed"`
	Percent string `json:"percent"`
}

type lorenzMedia struct {
	Type   string            `json:"type"`
	Values map[string]string `json:"values"`
}

type lorenzEventLocation struct {
	Place    string  `json:"place"`
	Accuracy float64 `json:"accuracy"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
}

type lorenzColour struct {
	Red   int `json:"red"`
	Green int `json:"green"`
	Blue  int `json:"blue"`
	Alpha int `json:"alpha"`
}

type lorenzTicketMode struct {
	Mode     lorenzMode `json:"mode"`
	Deadline string     `json:"deadline"`
}

type lorenzMode int

type lorenzPaymentAccount struct {
	// Legacy settings
	Account string `json:"account,omitempty"`

	// Stripe connect settings
	DestinationKey string `json:"destination_key,omitempty"`
}

type lorenzWaitingLists struct {
	Cost lorenzEventCost `json:"cost"`
}
