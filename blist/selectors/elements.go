package selectors

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Namespace represents an alias for a semaphore lock
type Namespace string

// Lock describes a way of returning a Namespace as a lock.
func (n Namespace) Lock() string {
	return fmt.Sprintf("ns-%s/lock", n.String())
}

// Prefix inserts a prefix for the current namespace, which returns a new
// namespace to use.
func (n Namespace) Prefix(ns Namespace) Namespace {
	return Namespace(fmt.Sprintf("%s-%s", ns.String(), n.String()))
}

func (n Namespace) String() string {
	return string(n)
}

// SemaphoreUnlock is a type alias for unlocking a lock.
type SemaphoreUnlock func() error

// The HealthStatus of the service.
type HealthStatus string

// Queue defines a type alias for describing the channel sent through the
// message bus
type Queue string

func (q Queue) String() string {
	return string(q)
}

// Class defines a type alias for describing the class of the message with in
// the queue.
type Class string

func (c Class) String() string {
	return string(c)
}

type Failure struct {
	Error error
}

// Key defines an alias for a key selector
type Key string

// Len returns the length of the string
func (k Key) Len() int {
	return len(k.String())
}

// Namespace provides a namespace from a key
func (k Key) Namespace() Namespace {
	return Namespace(strings.ToLower(k.String()))
}

func (k Key) String() string {
	return string(k)
}

// EmptyKey defines a key that is completely empty/zero
const EmptyKey = Key("")

type IdempotencyKey string

func (v IdempotencyKey) String() string {
	return string(v)
}

type KeyTxn struct {
	Key, Txn Key
}

// Range defines a a min and max number
type Range struct {
	Min, Max int
}

// WithIn defines a way to see if a value is within the range (intersects)
func (l Range) WithIn(x int) bool {
	return x >= l.Min && x <= l.Max
}

type Behaviours struct {
	CreatedAt, UpdatedAt time.Time
}

// Event represents an event from the backing store, which has been normalised.
type Event struct {
	Behaviours      Behaviours
	Id              Key
	Name            string
	PermName        string
	Dates           Dates
	Description     string
	Colour          Colour
	Agenda          []Agenda
	Tags            Tags
	PaymentAccounts PaymentAccounts
	CodeTypes       CodeTypes
	Place           Place
	State           State
	Related         Related
	Media           []Media
	Tickets         Tickets
	WaitingLists    WaitingLists
}

type Tickets struct {
	Cost           Cost
	Limits         Range
	ExpiryDuration time.Duration
	Transfer       TicketMode
	Returns        TicketMode
	TotalNum       uint64
}

type WaitingLists struct {
	Cost        Cost
	ExpiryRange DurationRange
}

type Agenda struct {
	Name, Details string
}

type Tags struct {
	Tags, Inherited []string
}

type MediaType string

func (m MediaType) String() string {
	return strings.ToLower(string(m))
}

type Media struct {
	Type   MediaType
	Values map[string]string
}

func (m Media) Get(key string) (string, bool) {
	if m.Type.String() == key && len(m.Values) > 0 {
		if val, ok := m.Values["url"]; ok {
			return val, true
		}
	}
	return "", false
}

type Medias []Media

func (m Medias) Get(key string) (string, bool) {
	for _, v := range m {
		if val, ok := v.Get(key); ok {
			return val, true
		}
	}
	return "", false
}

func (m Medias) Set(key, field, value string) Medias {
	if len(value) < 1 || len(field) < 1 || m.Contains(key) {
		return m
	}

	return append(m, Media{
		Type: MediaType(key),
		Values: map[string]string{
			field: value,
		},
	})
}

func (m Medias) Contains(key string) bool {
	lower := strings.ToLower(key)
	for _, v := range m {
		if v.Type.String() == lower {
			return true
		}
	}
	return false
}

func (m Medias) Val() []Media {
	return []Media(m)
}

type TicketMode struct {
	Mode     int
	Deadline time.Time
}

type Related struct {
	CityIds, ArtistIds, EventIds []bson.ObjectId
}

type State struct {
	Type                string
	SoldOut, CodeLocked bool
	Hidden, Deleted     bool
	Flags               EventFlags
}

type Place struct {
	Venue    string
	Address  string
	Location Location
}

type Location struct {
	Place    string
	Accuracy float64
	Lat, Lng float64
}

type Colour struct {
	Red, Green, Blue, Alpha int
}

func (c Colour) Hex() string {
	return fmt.Sprintf("#%02x%02x%02x", uint8(c.Red), uint8(c.Green), uint8(c.Blue))
}

func (c Colour) Get(name string) int {
	switch name {
	case "red":
		return c.Red
	case "green":
		return c.Green
	case "blue":
		return c.Blue
	case "alpha":
		return c.Alpha
	default:
		return 0
	}
}

func (c Colour) String() string {
	return fmt.Sprintf("%d,%d,%d,%d", c.Red, c.Green, c.Blue, c.Alpha)
}

type Cost struct {
	Currency string
	Amount   uint64
	Fees     map[FeeType]Fee
}

func (c Cost) String() string {
	return fmt.Sprintf("%s:%d", c.Currency, c.Amount)
}

type Fee struct {
	Type    FeeAmountType
	Fixed   uint64
	Percent float64
}

type FeeType string

func (t FeeType) String() string {
	return string(t)
}

type FeeAmountType string

func (t FeeAmountType) String() string {
	return string(t)
}

// DurationRange defines a min and max duration
type DurationRange struct {
	Min, Max time.Duration
}

// IsZero checks to see if the duration range is zero.
func (l DurationRange) IsZero() bool {
	return l.Min == 0 && l.Max == 0
}

type Dates struct {
	Start, End time.Time
}

// PaymentAccounts defines aliases for a typed accounts
type PaymentAccounts map[string]PaymentAccount

// Get returns the account if found and returns bool to know if it has been.
func (p PaymentAccounts) Get(name string) (PaymentAccount, bool) {
	account, ok := p[name]
	return account, ok
}

// PaymentAccountType describes the account type for the PaymentAccount
type PaymentAccountType int

// PaymentAccount represents a primary and destination pairing.
type PaymentAccount struct {
	Type        PaymentAccountType
	Destination string
}

// Payment represents a set of values required for purchasing items.
type Payment struct {
	Key            Key
	Txn            Key
	Cost           Cost
	Count          int
	Method         PaymentMethod
	UserInfo       PaymentUserInfo
	IdempotencyKey IdempotencyKey
}

func (p Payment) SetTxn(txn Key) Payment {
	return Payment{
		Key:            p.Key,
		Txn:            txn,
		Cost:           p.Cost,
		Count:          p.Count,
		Method:         p.Method,
		UserInfo:       p.UserInfo,
		IdempotencyKey: p.IdempotencyKey,
	}
}

func (p Payment) SetIdempotencyKey(key IdempotencyKey) Payment {
	return Payment{
		Key:            p.Key,
		Txn:            p.Txn,
		Cost:           p.Cost,
		Count:          p.Count,
		Method:         p.Method,
		UserInfo:       p.UserInfo,
		IdempotencyKey: key,
	}
}

// PaymentMethod defines a method of payment
type PaymentMethod struct {
	Type  PaymentMethodType
	Token Key
}

// PaymentMethodType describes how the type should be made
type PaymentMethodType string

func (t PaymentMethodType) String() string {
	return string(t)
}

// PaymentUserInfo represents a set of optional values that could be used for
// payments.
type PaymentUserInfo struct {
	FullName   string
	PostalCode string
}

// CodeTypes holds all the various different types that the event could be
// using.
type CodeTypes struct {
	BarcodeType BarcodeType
	QRCodeType  QRCodeType
}

// BarcodeType is a type alias for knowing which type of barcode is in use.
type BarcodeType string

func (b BarcodeType) String() string {
	return string(b)
}

// QRCodeType is a type alias for knowing which type of barcode is in use.
type QRCodeType string

func (b QRCodeType) String() string {
	return string(b)
}

// CodeSet represents a tuple of codes required for tickets
type CodeSet struct {
	Barcode Barcode
	QRCode  QRCode
}

// Barcode represents a high level version of what a barcode encompases
type Barcode struct {
	Type   BarcodeType
	Origin string
	Source string
}

// QRCode represents a high level version of what a qrcode encompases
type QRCode struct {
	Type   QRCodeType
	Source string
}

// User represents a user from the backing store.
type User struct {
	Id                          Key
	FirstName, LastName         string
	Phone, Email, DOB, Postcode string
	State                       UserState
	Codes, Tickets              map[Key][]Key
	Limits                      map[Key]Range
	Accounts                    []Account
	Access                      UserAccess
	PushIds                     UserPushIds
}

type Account struct {
	Type  AccountType
	Id    Key
	Cards []Card
}

type AccountType int

func (a AccountType) Int() int {
	return int(a)
}

// CardType defines the alternatives of card types we can use for payments
type CardType int

func (c CardType) Int() int {
	return int(c)
}

// Card represents a card used for payments
type Card struct {
	Type        CardType
	Id          Key
	Name        string
	Last4       string
	Brand       string
	Funding     string
	ExpMonth    uint
	ExpYear     uint
	FingerPrint string
}

type UserAccess struct {
	Token  string
	Expiry time.Time
}

type UserState struct {
	Status          UserStatus
	Type            UserType
	Deleted, Locked bool
	AllowMarketing  bool
	MultipleDevices bool
}

type UserStatus int

func (u UserStatus) Val() int {
	return int(u)
}

func (u UserStatus) String() string {
	switch u {
	case 1:
		return "non-authorized"
	case 2:
		return "authorized"
	}
	return "invalid"
}

type UserPushIds struct {
	Android, IOS []string
}

// Version defines an alias for the sum of a version
type Version string

func (v Version) String() string {
	return string(v)
}
